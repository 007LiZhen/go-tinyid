package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var Viper = viper.New()
var once = &sync.Once{}

func Init(path string) {
	once.Do(func() {
		Viper.SetConfigName("mysql")
		Viper.SetConfigType("toml")
		Viper.AddConfigPath(path)
		if err := Viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("err is (%w)", err))
		}

		for _, config := range Viper.GetStringSlice("include.configs") {
			Viper.SetConfigName(config)
			if err := Viper.MergeInConfig(); err != nil {
				panic(fmt.Errorf("err is (%w)", err))
			}
		}
		fmt.Println("config init success!")
	})
}
