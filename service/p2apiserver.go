// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

import (
	"net/http"

	"github.com/gorilla/mux"
)

// P2APIServer defines the http API Server
type P2APIServer struct {
	router *mux.Router
	Srv    *http.Server
}
