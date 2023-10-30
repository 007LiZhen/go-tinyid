package idsequence

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"

	"gitee.com/git-lz/go-tinyid/common/config"
	"gitee.com/git-lz/go-tinyid/common/mysql"
)

func init() {
	config.Viper.AddConfigPath("../../conf/")
	config.Init("")
	mysql.Init()
}

func TestNewIdSequence(t *testing.T) {
	type args struct {
		idListLength int64
		biz          string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case1",
			args: args{
				idListLength: 100,
				biz:          "demo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIdSequence(tt.args.idListLength, tt.args.biz)
			time.Sleep(2 * time.Second)
			got.stopMonitor <- true
			assert.Equal(t, int64(len(got.ids)), tt.args.idListLength)

			id, err := got.GetOne()
			assert.Equal(t, err, nil)
			fmt.Println("id is: ", id)

			got.Close()
			got.saveLastId(context.Background())
		})
	}
}
