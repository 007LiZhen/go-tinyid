package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitee.com/git-lz/go-tinyid/common/config"
	"gitee.com/git-lz/go-tinyid/common/mysql"
	"gitee.com/git-lz/go-tinyid/logic/idsequence"
	"gitee.com/git-lz/go-tinyid/router"
)

func main() {
	config.Viper.AddConfigPath("./conf")
	config.Init("")
	mysql.Init()
	idsequence.Init()
	go router.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR2, syscall.SIGKILL)

	for {
		select {
		case s := <-c:
			switch s {
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR2, syscall.SIGKILL, syscall.SIGQUIT:
				fmt.Printf("signal recieve: %s\n", s)
				idsequence.Stop()
				return
			}
		}
	}
}
