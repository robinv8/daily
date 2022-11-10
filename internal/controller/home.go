package controller

import (
	"fmt"
	"net/http"

	"daily/internal/service"

	"github.com/julienschmidt/httprouter"

	"encoding/json"
)

func NewHome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	siteInfo := service.Home()

	json, err := json.Marshal(siteInfo)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, string(json))
}
