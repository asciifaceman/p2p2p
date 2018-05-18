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

func (s *Server) informNode(node *Node) error {
	var conn *grpc.ClientConn

	conn, cerr := grpc.Dial(lib.FormatHostPort(node.Host, node.Port), grpc.WithInsecure())
	if cerr != nil {
		log.Printf("Could not contact: %v\n", cerr)
		return cerr
	}

	defer conn.Close()

	ni := NewInformServiceClient(conn)

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
		if poolNode.Name == node.Name && poolNode.Port == node.Port {
			continue
		}
		thisNode := &NodeMessage{
			Name: poolNode.Name,
			Host: poolNode.Host,
			Port: int32(poolNode.Port),
		}
		informPayload.Pool = append(informPayload.Pool, thisNode)
	}

	// Send the target node our phonebook
	response, rerr := ni.InformNode(context.Background(), informPayload)
	if rerr != nil {
		log.Printf("Error calling remote server %s\n", rerr)
		return rerr
	}

	if !response.Response {
		return fmt.Errorf("%s could not consume my phonebook", node.Name)
	}

	return nil
}
