// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import (
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
	g.Go(func() error { return m.Serve() })

	log.Printf("Node [%s] is booting. Listening on [%s:%d]", s.Name, s.Host, s.Port)

	log.Println("Server running:", g.Wait())

	// Check bootnodes and configure/test if available
	//if len(nodePool) > 0 {
	//  fmt.Printf("I have bootnodes, contacting...\n")
	//}

	return nil
}

// formatBootnodes takes CLI string and formats into the Pool struct
func (s *Server) formatBootnodes(nodes string) *Pool {
	pool := &Pool{}
	splitPool := strings.Split(nodes, ",")
	for i := range splitPool {
		host := strings.Split(splitPool[i], ":")
		thisNode := Node{Host: host[0], Port: lib.ToInt(host[1])}
		pool.nodes = append(pool.nodes, &thisNode)
	}
	return pool
}
