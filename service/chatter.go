// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/asciifaceman/p2p2p/lib"
	"google.golang.org/grpc"
)

const (
	emptyString string        = ""
	nanosecond  time.Duration = 1
	microsecond               = 1000 * nanosecond
	millisecond               = 1000 * microsecond
	second                    = 1000 * millisecond
	minute                    = 60 * second
	hour                      = 60 * minute
)

func (s *Server) getNodeNameFromNode(host string, port int) (*Node, error) {
	var conn *grpc.ClientConn

	conn, cerr := grpc.Dial(lib.FormatHostPort(host, port), grpc.WithInsecure(), grpc.WithBackoffMaxDelay(minute))
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

func (s *Server) sendWhisper(node *NodeMessage, message string) error {
	var conn *grpc.ClientConn

	conn, cerr := grpc.Dial(lib.FormatHostPort(node.Host, int(node.Port)), grpc.WithInsecure(), grpc.WithBackoffMaxDelay(minute))
	if cerr != nil {
		log.Printf("Could not contact: %v\n", cerr)
		return cerr
	}

	defer conn.Close()

	sw := NewWhisperClient(conn)

	response, rerr := sw.SendWhisper(context.Background(), &WhisperMessage{Source: s.Me.Name, Body: message})
	if rerr != nil {
		log.Printf("Failed to send a whisper to %s\n", node.Name)
		return rerr
	}

	fmt.Printf("ACK: %v", response)

	return nil
}

func (s *Server) requestNode(node *Node, name string) (*NodeMessage, error) {
	var conn *grpc.ClientConn

	conn, cerr := grpc.Dial(lib.FormatHostPort(node.Host, node.Port), grpc.WithInsecure(), grpc.WithBackoffMaxDelay(minute))
	if cerr != nil {
		log.Printf("Could not contact: %v\n", cerr)
		return &NodeMessage{}, cerr
	}

	defer conn.Close()

	nq := NewInformServiceClient(conn)

	requestPayload := &NodeRequestMessage{
		Informer: s.NodeToMessage(s.Me),
		Request:  name,
	}

	// ask our target node for this node
	response, rerr := nq.RequestNode(context.Background(), requestPayload)
	if rerr != nil {
		log.Printf("Error calling remote server %s\n", rerr)
		return &NodeMessage{}, rerr
	}

	if response.Found {
		return response.Contents, nil
	}

	return &NodeMessage{}, errors.New("Not found")
}

func (s *Server) informNode(node *Node) error {
	var conn *grpc.ClientConn

	conn, cerr := grpc.Dial(lib.FormatHostPort(node.Host, node.Port), grpc.WithInsecure(), grpc.WithBackoffMaxDelay(minute))
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
