package main

import (
	"net/http"

	"daily/internal/config"
	"daily/internal/controller"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// 加载环境变量
	config.LoadEnv()
	
	router := httprouter.New()
	router.GET("/", controller.NewHome)
	router.GET("/blog", controller.NewBlog)
	// New endpoint for importing tweets from URL
	router.GET("/api/site/import", controller.ImportFromURL)
	http.ListenAndServe(":4000", router)
}
