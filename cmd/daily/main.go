package main

import (
	"net/http"

	"daily/internal/controller"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", controller.NewHome)
	router.GET("/blog", controller.NewBlog)
	http.ListenAndServe(":4000", router)
}
