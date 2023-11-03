package idsequence

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("case-normal", func(t *testing.T) {
		Init()
		fmt.Println("IdSequenceMap is: ", IdSequenceMap)
	})
}

func TestGetIdSequence(t *testing.T) {
	t.Run("case-success", func(t *testing.T) {
		Init()
		got, err := GetIdSequence("demo")
		assert.Nil(t, err)
		assert.Equal(t, got.biz, "demo")
	})

	t.Run("case-biz not support", func(t *testing.T) {
		Init()
		biz := "demo1"
		got, err := GetIdSequence(biz)
		assert.Equal(t, err, fmt.Errorf("biz=(%s) nor support", biz))
		assert.Nil(t, got)
	})
}

func TestStop(t *testing.T) {
	t.Run("case-normal", func(t *testing.T) {
		Init()

		got, err := GetIdSequence("demo")
		assert.Nil(t, err)

		Stop()

		_, ok := <-got.ids
		assert.Equal(t, ok, false)
	})
}