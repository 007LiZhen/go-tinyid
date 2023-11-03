package dao

import (
	"context"

	"github.com/007LiZhen/go-tinyid/dao/basedao"
	"github.com/007LiZhen/go-tinyid/model"
)

type idSequenceDao struct {
	*basedao.BaseDao
}

func NewIdSequenceDao() *idSequenceDao {
	return &idSequenceDao{
		basedao.NewBaseDao(model.TableNameIdSequence, model.NewIdSequence()),
	}
}

func (d *idSequenceDao) SetModel(model *model.IdSequence) *idSequenceDao {
	tx := d
	tx.BaseDao.SetModel(model)
	return tx
}

func (d *idSequenceDao) UpdateByCond(ctx context.Context, cond map[string]interface{}, update map[string]interface{}) (int64, error) {
	db := d.Db.WithContext(ctx).Table(d.TableName)

	for k, v := range cond {
		db = db.Where(k+" = ?", v)
	}

	db = db.Updates(update)
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
