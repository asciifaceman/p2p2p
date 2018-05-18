// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import (
	"context"
	"fmt"
	"log"

	"github.com/asciifaceman/p2p2p/lib"
	"google.golang.org/grpc"
)

const emptyString string = ""

func (s *Server) getNodeNameFromNode(host string, port int) (*Node, error) {
	var conn *grpc.ClientConn

	conn, cerr := grpc.Dial(lib.FormatHostPort(host, port), grpc.WithInsecure())
	if cerr != nil {
		log.Printf("Could not contact: %v\n", cerr)
		return &Node{}, cerr
	}

	defer conn.Close()

	n := NewNameClient(conn)

	response, rerr := n.SayName(context.Background(), &NodeMessage{Name: s.Name, Host: s.Host, Port: int32(s.Port)})
	if rerr != nil {
		log.Printf("Error calling remote server: %s\n", rerr)
		return &Node{}, rerr
	}

	node := &Node{Name: response.Name, Host: response.Host, Port: int(response.Port)}

	return node, nil
}

// BuildInformPool builds the pool in a format for gRPC
func (s *Server) BuildInformPool(name string, port int) *NodeInformMessage {
	// TODO: Fix this variables to be a struct, but must think on it
	// This was the quickest path for me
	// On a future refactor once everything is working I will shift this
	// burden on to the recipients

	// Our basic inform payload
	informPayload := &NodeInformMessage{
		Informer: &NodeMessage{
			Name: s.Name,
			Host: s.Host,
			Port: int32(s.Port),
		},
		Pool: []*NodeMessage{},
	}

	// Iter through our nodes, rip out our destination, and pack the rest
	for _, poolNode := range s.Me.Pool.nodes {
		if poolNode.Name == name && poolNode.Port == port {
			continue
		}
		thisNode := &NodeMessage{
			Name: poolNode.Name,
			Host: poolNode.Host,
			Port: int32(poolNode.Port),
		}
		informPayload.Pool = append(informPayload.Pool, thisNode)
	}

	return informPayload
}

func (s *Server) informNode(node *Node) error {
	var conn *grpc.ClientConn

	conn, cerr := grpc.Dial(lib.FormatHostPort(node.Host, node.Port), grpc.WithInsecure())
	if cerr != nil {
		log.Printf("Could not contact: %v\n", cerr)
		return cerr
	}

	defer conn.Close()

	ni := NewInformServiceClient(conn)

	informPayload := s.BuildInformPool(node.Name, node.Port)

	// Send the target node our phonebook
	response, rerr := ni.InformNode(context.Background(), informPayload)
	if rerr != nil {
		log.Printf("Error calling remote server %s\n", rerr)
		return rerr
	}

	// incorporate their response
	if len(response.Pool) > 0 {
		fmt.Printf("Incorporating %s's phonebook...\n", node.Name)

		for _, poolNode := range response.Pool {
			err := s.AddNodeToPool(poolNode)
			if err != nil {
				log.Printf("Something serious went wrong! %v\n", err)
				continue
			}
		}
	}

	return nil
}

// CheckPoolForNode checks my pool for a nodes name and port
func (s *Server) CheckPoolForNode(node *NodeMessage) bool {
	poolNames := make(map[string]int)

	// If we already know the server we don't have to do aything
	for _, poolNode := range s.Me.Pool.nodes {
		poolNames[poolNode.Name] = poolNode.Port
	}

	if val, ok := poolNames[node.Name]; ok && val == int(node.Port) {
		return true
	}

	return false
}

// AddNodeToPool checks for a nodes existence in its pool and adds if it doesn't
func (s *Server) AddNodeToPool(node *NodeMessage) error {
	// need to think on this, it feels wrong
	if s.CheckPoolForNode(node) {
		log.Printf("I already know [%s@%s:%d]. Welcome back!\n", node.Name, node.Host, node.Port)
	} else {
		log.Printf("[%s@%s:%d] is new to me. Adding to my phonebook.\n", node.Name, node.Host, node.Port)
		newNode := &Node{
			Name: node.Name,
			Host: node.Host,
			Port: int(node.Port),
		}
		s.Me.Pool.nodes = append(s.Me.Pool.nodes, newNode)
	}
	return nil
}
