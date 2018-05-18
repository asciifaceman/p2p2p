// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

// P2Server defines the http API Server
type P2Server struct {
	router *mux.Router
	Srv    *http.Server
}
