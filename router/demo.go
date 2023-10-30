package router

import (
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"gitee.com/git-lz/go-tinyid/controller"
)

const (
	pathPrefix = "/demo"
)

func Init() {
	r := gin.New()

	groupV1 := r.Group(pathPrefix + "/v1")
	groupV1.GET("/id", controller.NewID().Get)

	if err := endless.ListenAndServe(":8099", r); err != nil {
		fmt.Println("server err is: ", err)
	}
}
