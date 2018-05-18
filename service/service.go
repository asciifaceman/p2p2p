// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import (
	"errors"
	"log"
	"strings"

	"github.com/asciifaceman/p2p2p/lib"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
)

// Start is the external entrypoint to boot up and listen
func (s *Server) Start(nodePool string) error {

	// Define this node
	s.Me = &Node{Name: s.Name, Host: s.Host, Port: s.Port}

	// Create our mux'd listener
	m := s.NewListener()

	// Multiplex
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	// Run them wrapped in an errgroup
	g := new(errgroup.Group)
	g.Go(func() error { return s.grpcServe(grpcListener) })
	g.Go(func() error { return s.httpServe(httpListener) })
	g.Go(func() error { return s.checkBootnodes(nodePool) })
	g.Go(func() error { return m.Serve() })

	log.Printf("==============================:\n")
	log.Printf("Welcome to p2p2p - The naive graph network that does very little.\n")
	log.Printf("==============================:\n")
	log.Printf("\n\n")
	log.Printf("Node [%s] is booting. Listening on [%s:%d]\n", s.Name, s.Host, s.Port)

	log.Println("Server running:", g.Wait())

	// Check bootnodes and configure/test if available
	//if len(nodePool) > 0 {
	//  fmt.Printf("I have bootnodes, contacting...\n")
	//}

	return nil
}

func (s *Server) checkBootnodes(nodes string) error {
	if len(nodes) > 0 {
		log.Printf("[network] I have bootnodes... contacting...\n")

		// Format my input, contact the nodes, and add them to my pool
		s.Me.Pool = *s.formatBootnodes(nodes)

		//ierr := s.InformPoolOfNodes()
		//if ierr != nil {
		//	log.Printf("%v", ierr)
		//}

		return nil
	}

	return errors.New("No bootnodes")
}

// formatBootnodes takes CLI string, contacts nodes, and generates pool
func (s *Server) formatBootnodes(nodes string) *Pool {
	pool := &Pool{}
	splitPool := strings.Split(nodes, ",")
	for i := range splitPool {
		host := strings.Split(splitPool[i], ":")
		bootNode, err := s.getNodeNameFromNode(host[0], lib.ToInt(host[1]))
		if err != nil {
			log.Fatalf("Failed to retrieve the name for node [%s:%d]. Ignoring and moving on (backoff not implemented).\n", host[0], lib.ToInt(host[1]))
			continue
		}
		thisNode := Node{Name: bootNode.Name, Host: host[0], Port: lib.ToInt(host[1])}
		log.Printf("Found %s@%s:%d!, adding to available node pool.\n", thisNode.Name, thisNode.Host, thisNode.Port)
		pool.nodes = append(pool.nodes, &thisNode)
	}
	return pool
}
