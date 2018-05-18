// Copyright © 2018 Charles Corbett <nafredy@gmail.com>
//

package service

import (
	"context"
	"log"
)

// SayName generates a response to a Name request
func (s *Server) SayName(ctx context.Context, in *NameMessage) (*NameMessage, error) {
	log.Printf("Receive message: %s", in.Name)
	return &NameMessage{Name: s.Name}, nil
}

// SendWhisper receives and processes a sent whisper and may respond with a receipt
func (s *Server) SendWhisper(ctx context.Context, in *WhisperMessage) (*WhisperAck, error) {
	log.Printf("Receive message from [%s]: %s", in.Source, in.Body)
	return &WhisperAck{Response: true}, nil
}
