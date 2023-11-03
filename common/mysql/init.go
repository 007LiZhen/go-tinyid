package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/007LiZhen/go-tinyid/common/config"
)

var DB = new(gorm.DB)

func Init() func() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.Viper.GetString("mysql.dsn")))
	if err != nil {
		panic(fmt.Errorf("mysql err is (%w)", err))
	}

	mysqlDb, err := DB.DB()
	if err != nil {
		panic(fmt.Errorf("mysql err is (%w)", err))
	}
	mysqlDb.SetConnMaxLifetime(time.Second * config.Viper.GetDuration("mysql.max_life_time"))
	mysqlDb.SetMaxOpenConns(config.Viper.GetInt("mysql.max_open_conns"))
	mysqlDb.SetMaxIdleConns(config.Viper.GetInt("mysql.max_idle_conns"))

	cancelFunc := func() {
		if err := mysqlDb.Close(); err != nil {
			fmt.Println("mysql close err is: ", err)
		} else {
			fmt.Println("mysql close success")
		}
	}

	fmt.Println("mysql init success!")
	return cancelFunc
}
