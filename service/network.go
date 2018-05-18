// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/asciifaceman/p2p2p/lib"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

func (s *Server) buildRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", s.healthHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	return r
}

// serve http
func (s *Server) httpServe(l net.Listener) error {
	s.router = s.buildRoutes()
	s.Srv = s.NewServer(s.router)
	log.Printf("API Listening...\n")
	return s.Srv.Serve(l)
}

// NewServer generates a new server object for http connections
func (s *Server) NewServer(r *mux.Router) *http.Server {
	serv := &http.Server{
		Handler:      handlers.RecoveryHandler()(r),
		Addr:         lib.FormatHostPort(s.Host, s.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return serv
}

// NewListener generates a new listener object and returns a mux of it
func (s *Server) NewListener() cmux.CMux {
	lis, err := net.Listen("tcp", lib.FormatHostPort(s.Host, s.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	m := cmux.New(lis)

	return m
}

// serve grpc
func (s *Server) grpcServe(l net.Listener) error {
	// gRPC server object
	grpcServer := grpc.NewServer()

	// attach the name service
	RegisterNameServer(grpcServer, s)
	RegisterWhisperServer(grpcServer, s)

	log.Printf("gRPC Listening...\n")

	// start gRPC server
	return grpcServer.Serve(l)
}
