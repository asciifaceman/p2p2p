// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Server defines the http API Server
type Server struct {
	router *mux.Router
	Srv    *http.Server
	Host   string
	Port   int
	Name   string
	Me     *Node
}
