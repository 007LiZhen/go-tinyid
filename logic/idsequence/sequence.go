package idsequence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitee.com/git-lz/go-tinyid/dao"
	"gitee.com/git-lz/go-tinyid/model"
)

type IdSequence struct {
	idListLength int64
	biz          string
	ids          chan int64
	stopMonitor  chan bool
}

func Stop() {
	for _, idSequence := range IdSequenceMap {
		idSequence.stopMonitor <- true
		idSequence.Close()
		idSequence.saveLastId(context.Background())
	}
}

func GetIdSequence(biz string) (*IdSequence, error) {
	idSequence, ok := IdSequenceMap[biz]
	if !ok {
		return nil, fmt.Errorf("biz=(%s) nor support", biz)
	}

	return idSequence, nil
}

func NewIdSequence(idListLength int64, biz string) *IdSequence {
	is := &IdSequence{
		ids:          make(chan int64, idListLength),
		stopMonitor:  make(chan bool, 1),
		idListLength: idListLength,
		biz:          biz,
	}

	go is.Monitor(context.Background())
	return is
}

func (is *IdSequence) Close() {
	close(is.ids)
	close(is.stopMonitor)
}

func (is *IdSequence) Monitor(ctx context.Context) {
	for {
		select {
		case <-is.stopMonitor:
			fmt.Println("stop monitor success...")
			return
		default:
			if len(is.ids) == 0 {
				if err := is.add(ctx); err != nil {
					fmt.Println("monitor id sequence err: ", err)
				} else {
					fmt.Println("add ids success...")
				}
			}
		}
	}
}

func (is *IdSequence) GetOne() (int64, error) {
	// 如果取不到，则会等待1s
	ticker := time.After(time.Second)

	for {
		select {
		case id, ok := <-is.ids:
			if ok && id != 0 {
				return id, nil
			}
		case <-ticker:
			return 0, fmt.Errorf("get next id timeout")
		}
	}
}

func (is *IdSequence) add(ctx context.Context) error {
	minId, maxId, err := is.getNewIdListLoop(ctx)
	if err != nil {
		return err
	}

	for i := minId; i < maxId; i++ {
		is.ids <- i
	}

	return nil
}

func (is *IdSequence) getNewIdListLoop(ctx context.Context) (int64, int64, error) {
	ticker := time.After(10 * time.Second)

	var (
		errCh                  = make(chan error, 1)
		curMaxIdCh, newMaxIdCh = make(chan int64, 1), make(chan int64, 1)
	)

	for {
		select {
		case <-ticker:
			return 0, 0, errors.New("timeout")
		case err := <-errCh:
			if err == nil {
				return <-curMaxIdCh, <-newMaxIdCh, nil
			}
		default:
			is.getNewIdListSync(ctx, curMaxIdCh, newMaxIdCh, errCh)
			time.Sleep(time.Millisecond)
		}
	}
}

func (is *IdSequence) getNewIdListSync(ctx context.Context, curMaxIdCh, newMaxIdCh chan int64, errCh chan error) {
	idSequenceDao := dao.NewIdSequenceDao()

	var records []model.IdSequence
	if err := idSequenceDao.GetsByCond(ctx, map[string]interface{}{
		"biz": is.biz,
	}, &records); err != nil {
		errCh <- err
		return
	}

	if len(records) == 0 {
		if err := idSequenceDao.SetModel(model.NewIdSequence().SetBiz(is.biz).SetValue(0).SetValue(0)).
			InsertOneRecord(ctx); err != nil {
			errCh <- err
			return
		}

		curMaxIdCh <- 0
		newMaxIdCh <- is.idListLength
		errCh <- nil
		return
	}

	if len(records) != 1 {
		errCh <- fmt.Errorf("data=(%v) has not only one record", records)
		return
	}

	record := records[0]
	newMaxId := record.Value + is.idListLength

	rowsEffect, err := idSequenceDao.UpdateByCond(ctx, record.Version, map[string]interface{}{
		"value":   newMaxId,
		"version": 0,
	})
	if err != nil {
		errCh <- err
		return
	}

	if rowsEffect != 1 {
		errCh <- errors.New("data has been changed, need to retry")
		return
	}

	curMaxIdCh <- record.Value
	newMaxIdCh <- newMaxId
	errCh <- nil
}

func (is *IdSequence) saveLastId(ctx context.Context) {
	lastId, ok := <-is.ids
	if ok && lastId != 0 {
		dao.NewIdSequenceDao().UpdateByCond(ctx, 0, map[string]interface{}{
			"value":   lastId,
			"version": 0,
		})
	}
}
