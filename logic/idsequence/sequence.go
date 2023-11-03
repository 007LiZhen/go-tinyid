package idsequence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/007LiZhen/go-tinyid/dao"
	"github.com/007LiZhen/go-tinyid/model"
)

// IdSequence - The generator struct
type IdSequence struct {
	idListLength int64      // the length of id list, you can define it by your business
	biz          string     // business type
	ids          chan int64 // the channel of id list
	stopMonitor  chan bool  // the stop monitor signal
}

// NewIdSequence - Create a new IdSequence by idListLength and biz
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

// Close - You can use it to close the id list channel and stop monitor goroutine
func (is *IdSequence) Close() {
	close(is.ids)
	close(is.stopMonitor)
}

// Monitor - This function can monitor the capacity of the id list. If the capacity of the id list is zero, it will
//
//	add new ids to the list.
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

// GetOne - Get a new id from the id list channel.
func (is *IdSequence) GetOne() (int64, error) {
	// if there is no new id, sleep 1s
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

// add - Add new ids to the id list channel
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

// getNewIdListLoop - Add new id to the id list channel for spin with optimistic lock.
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

// getNewIdListSync - Get new max id from the mysql with optimistic lock.
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
	version := record.Version + 1

	rowsEffect, err := idSequenceDao.UpdateByCond(ctx, map[string]interface{}{
		"version": record.Version,
		"biz":     is.biz,
	}, map[string]interface{}{
		"value":   newMaxId,
		"version": version,
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

// saveLastId - save the last id when the process is being stopped
func (is *IdSequence) saveLastId(ctx context.Context) {
	lastId, ok := <-is.ids
	if ok && lastId != 0 {
		dao.NewIdSequenceDao().UpdateByCond(ctx, map[string]interface{}{
			"biz": is.biz,
		}, map[string]interface{}{
			"value": lastId,
		})
	}
}
