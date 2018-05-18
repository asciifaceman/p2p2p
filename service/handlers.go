// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/asciifaceman/p2p2p/lib"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("/health hit.")
	lib.RespondJSON(w, 200, ResponseMessage{Status: "Ok", Body: fmt.Sprintf("[%s] is Alive and listening on [%d]", s.Name, s.Port)})
}
