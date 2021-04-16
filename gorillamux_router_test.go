// +build !gingonic,!echo,gorillamux

package api2go

import (
	"log"

	"github.com/charm-jp/api2go/routing"
	"github.com/gorilla/mux"
)

func newTestRouter() routing.Routeable {
	router := mux.NewRouter()
	router.MethodNotAllowedHandler = notAllowedHandler{}
	return routing.Gorilla(router)
}

func init() {
	log.Println("Testing with gorilla router")
}
