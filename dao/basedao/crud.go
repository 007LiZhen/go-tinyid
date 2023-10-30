package basedao

import (
	"context"
	"reflect"
)

func (b *BaseDao) InsertOneRecord(ctx context.Context) error {
	db := b.Db.WithContext(ctx).Table(b.TableName).Create(b.Model)
	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (b *BaseDao) GetsByCond(ctx context.Context, cond map[string]interface{}, res interface{}) error {
	db := b.Db.WithContext(ctx).Table(b.TableName)

	for k, v := range cond {
		switch k {
		case "limit":
			db = db.Limit(v.(int)).Order("ID DESC")
		case "offset":
			db = db.Offset(v.(int))
		default:
			switch reflect.TypeOf(v).Kind() {
			case reflect.Array, reflect.Slice:
				db = db.Where(k+" in (?)", v)
			default:
				db = db.Where(k+" = ?", v)
			}
		}
	}

	db = db.Find(res)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
