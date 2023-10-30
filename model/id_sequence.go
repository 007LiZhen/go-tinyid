package model

import "time"

// IdSequence id-sequence表
type IdSequence struct {
	ID         uint64    `gorm:"primaryKey" json:"id"` // primary id
	Biz        string    `gorm:"biz" json:"biz"`
	Value      int64     `gorm:"column:value;NOT NULL" json:"value"`
	Version    int       `gorm:"version" json:"version"`
	IsDel      int8      `gorm:"column:is_del;default:0;NOT NULL" json:"is_del"` // 是否软删除 0-正常未被软删除 1-已被软删除
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
}

const (
	TableNameIdSequence = "sequence"
)

func NewIdSequence() *IdSequence {
	return &IdSequence{}
}

func (m *IdSequence) TableName() string {
	return TableNameIdSequence
}

func (m *IdSequence) SetBiz(biz string) *IdSequence {
	tx := m
	tx.Biz = biz
	return tx
}

func (m *IdSequence) SetValue(value int64) *IdSequence {
	tx := m
	tx.Value = value
	return tx
}

func (m *IdSequence) SetVersion(version int) *IdSequence {
	tx := m
	tx.Version = version
	return tx
}
