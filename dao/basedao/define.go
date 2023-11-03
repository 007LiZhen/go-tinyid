package basedao

import (
	"gorm.io/gorm"

	"github.com/007LiZhen/go-tinyid/common/mysql"
)

type BaseDao struct {
	Db        *gorm.DB
	Model     interface{}
	TableName string
}

func NewBaseDao(tableName string, m interface{}) *BaseDao {
	return &BaseDao{
		Db:        mysql.DB,
		TableName: tableName,
		Model:     m,
	}
}

func (b *BaseDao) SetModel(m interface{}) *BaseDao {
	tx := b
	tx.Model = m
	return tx
}

func (b *BaseDao) Begin() *BaseDao {
	db := *b.Db
	tx := &BaseDao{
		Db:        db.Begin(),
		TableName: b.TableName,
		Model:     b.Model,
	}
	return tx
}

func (b *BaseDao) Commit() *BaseDao {
	tx := *b
	tx.Db.Commit()
	return &tx
}

func (b *BaseDao) Rollback() *BaseDao {
	tx := *b
	tx.Db.Rollback()
	return &tx
}
