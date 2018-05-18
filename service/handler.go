// Copyright © 2018 Charles Corbett <nafredy@gmail.com>
//

package service

import (
	"context"
	"log"
)

// SayName generates a response to a Name request
func (s *Server) SayName(ctx context.Context, in *NodeMessage) (*NodeMessage, error) {
	log.Printf("Name request received from %s, checking against my phonebook...\n", in.Name)

	// Naively maintaining a cache of Nodes
	err := s.AddNodeToPool(in)
	if err != nil {
		log.Printf("Something serious went wrong! %v", err)
	}
	log.Printf("[%s@%s:%d] Added. Responding....\n", in.Name, in.Host, in.Port)

	return &NodeMessage{Name: s.Name, Host: s.Host, Port: int32(s.Port)}, nil
}

// SendWhisper receives and processes a sent whisper and may respond with a receipt
func (s *Server) SendWhisper(ctx context.Context, in *WhisperMessage) (*WhisperAck, error) {
	log.Printf("Receive message from [%s]: %s\n", in.Source, in.Body)
	return &WhisperAck{Response: true}, nil
}

// InformNode informs the target node of who it is, and retrieves the others identity
func (s *Server) InformNode(ctx context.Context, in *NodeInformMessage) (*NodeInformMessage, error) {
	log.Printf("Received phonebook from [%s@%s:%d], processing...\n", in.Informer.Name, in.Informer.Host, in.Informer.Port)
	for _, poolNode := range in.Pool {
		err := s.AddNodeToPool(poolNode)
		if err != nil {
			// I haven't fully thought out how to inform the informer of this
			log.Printf("Something serious went wrong! %v", err)
			continue
		}
	}

	informPool := s.BuildInformPool(in.Informer.Name, int(in.Informer.Port))

	log.Printf("Acknowledging receipt & attempted incorporation.\n")
	return informPool, nil
}

// RequestNode responds to a request to check our pool for a specific node
func (s *Server) RequestNode(ctx context.Context, in *NodeRequestMessage) (*NodeRequestReply, error) {
	log.Printf("Received a request from [%s@%s:%d] to look up [%s@??:??], processing...", in.Informer.Name, in.Informer.Host, in.Informer.Port, in.Request)
	node, found := s.CheckPoolForNodeByName(in.Request)
	if found {
		return &NodeRequestReply{Found: true, Contents: node}, nil
	}

	return &NodeRequestReply{Found: false, Contents: &NodeMessage{}}, nil
}
