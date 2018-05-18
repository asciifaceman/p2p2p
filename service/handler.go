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

	poolNames := make(map[string]int)

	// If we already know the server we don't have to do aything
	for _, poolNode := range s.Me.Pool.nodes {
		poolNames[poolNode.Name] = poolNode.Port
	}

	// Naively checking if I know the node or not
	if val, ok := poolNames[in.Name]; ok && val == int(in.Port) {

		log.Printf("I already know the source node. Long time no see, %s!\n", in.Name)

	} else {
		log.Printf("[%s@%s:%d] is new to me. Adding to my phonebook.\n", in.Name, in.Host, in.Port)
		newNode := &Node{
			Name: in.Name,
			Host: in.Host,
			Port: int(in.Port),
		}
		s.Me.Pool.nodes = append(s.Me.Pool.nodes, newNode)
		log.Printf("[%s@%s:%d] Added. Responding....\n", in.Name, in.Host, in.Port)
	}

	return &NodeMessage{Name: s.Name, Host: s.Host, Port: int32(s.Port)}, nil
}

// SendWhisper receives and processes a sent whisper and may respond with a receipt
func (s *Server) SendWhisper(ctx context.Context, in *WhisperMessage) (*WhisperAck, error) {
	log.Printf("Receive message from [%s]: %s\n", in.Source, in.Body)
	return &WhisperAck{Response: true}, nil
}

// InformNode informs the target node of who it is, and retrieves the others identity
func (s *Server) InformNode(ctx context.Context, in *NodeInformMessage) (*WhisperAck, error) {
	log.Printf("Received phonebook from [%s@%s:%d], processing...\n", in.Informer.Name, in.Informer.Host, in.Informer.Port)
	log.Printf("%v", in.Pool)
	// add to pool
	// respond
	log.Printf("Acknowledging successful receipt.")
	return &WhisperAck{Response: true}, nil
}
