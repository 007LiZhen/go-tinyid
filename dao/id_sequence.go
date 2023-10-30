package dao

import (
	"context"

	"gitee.com/git-lz/go-tinyid/dao/basedao"
	"gitee.com/git-lz/go-tinyid/model"
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

func (d *idSequenceDao) UpdateByCond(ctx context.Context, version int, update map[string]interface{}) (int64, error) {
	db := d.Db.WithContext(ctx).Table(d.TableName).Where("version = (?)", version).
		Updates(update)
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
