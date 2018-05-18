// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/asciifaceman/p2p2p/lib"
	"github.com/gorilla/mux"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Health check.")

	nodeCount := len(s.Me.Pool.nodes)

	lib.RespondJSON(w, 200, ResponseMessage{Status: "Ok", Responder: s.Me.Name, Body: fmt.Sprintf("[%s] is Alive and listening on [%d]. It is aware of [%d] nodes in its network.", s.Name, s.Port, nodeCount)})
}

func (s *Server) whisperHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]

	v := r.URL.Query()
	message := v.Get("message")

	if len(strings.TrimSpace(name)) == 0 || len(strings.TrimSpace(message)) == 0 {
		lib.RespondJSON(w, 400, ResponseMessage{Status: "Ok", Responder: s.Me.Name, Body: fmt.Sprintf("Request missing parameters.")})
		return
	}

	log.Printf("Received request to whisper %s...\n", name)

	node, found := s.CheckPoolForNodeByName(name)
	if !found {
		log.Printf("I am not familiar with %s, so I am going to ask my peers...", name)
		for _, thisNode := range s.Me.Pool.nodes {
			found, ferr := s.requestNode(thisNode, name)
			if ferr != nil {
				log.Printf("%v", ferr)
			}
			s.AddNodeToPool(found)
			break
		}
		node, found = s.CheckPoolForNodeByName(name)
		if !found {
			log.Printf("Failed to retrieve node from my network. Giving up.")
			lib.RespondJSON(w, 500, ResponseMessage{Status: "Ok", Responder: s.Me.Name, Body: fmt.Sprintf("Failed to discover %s.", name)})
			return
		}
	}

	s.sendWhisper(node, message)

	lib.RespondJSON(w, 200, ResponseMessage{Status: "Ok", Responder: s.Me.Name, Body: fmt.Sprintf("Whisper hit.")})

	log.Printf("Sending whisper to %s\n", name)

}
