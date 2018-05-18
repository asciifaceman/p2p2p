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

// AskPeersForNode asks the peers for a node
func (s *Server) AskPeersForNode(name string, exclude []string) ([]string, error) {
	log.Printf("I am not familiar with %s, so I am going to ask my peers...", name)
	for _, thisNode := range s.Me.Pool.nodes {
		for _, e := range exclude {
			if e == thisNode.Name {
				continue
			}
		}
		log.Printf("Asking %s for %s", thisNode.Name, name)
		found, ferr := s.requestNode(thisNode, name, exclude)
		if ferr != nil {
			log.Printf("%v", ferr)
			exclude = append(exclude, thisNode.Name)
			continue
		}
		s.AddNodeToPool(found)
		break
	}
	return exclude, nil
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
		var exclude []string
		exclude = append(exclude, s.Me.Name)
		s.AskPeersForNode(name, exclude)
	}

	node, found = s.CheckPoolForNodeByName(name)
	if !found {
		log.Printf("Failed to retrieve node from my network. Giving up.")
		lib.RespondJSON(w, 500, ResponseMessage{Status: "Ok", Responder: s.Me.Name, Body: fmt.Sprintf("Failed to discover %s. I will ask my network, try again later.", name)})
		return
	}

	werr := s.sendWhisper(node, message)
	if werr != nil {
		lib.RespondJSON(w, 404, ResponseMessage{Status: "Ok", Responder: s.Me.Name, Body: fmt.Sprintf("Failed to send message to %s. Reason: %v", name, werr)})
		return
	}

	lib.RespondJSON(w, 200, ResponseMessage{Status: "Ok", Responder: s.Me.Name, Body: fmt.Sprintf("Whisper hit.")})
}
