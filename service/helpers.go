// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

import "log"

// MessageToNode converts NodeMessage{} to Node{}
func (s *Server) MessageToNode(inNode *NodeMessage) *Node {
	newNode := &Node{
		Name: inNode.Name,
		Host: inNode.Host,
		Port: int(inNode.Port),
	}
	return newNode
}

// NodeToMessage converts Node{} to NodeMessage{}
func (s *Server) NodeToMessage(inNode *Node) *NodeMessage {
	newNode := &NodeMessage{
		Name: inNode.Name,
		Host: inNode.Host,
		Port: int32(inNode.Port),
	}
	return newNode
}

// CheckPoolForNode checks my pool for a nodes name and port
func (s *Server) CheckPoolForNode(node *NodeMessage) bool {
	poolNames := make(map[string]int)

	// If we already know the server we don't have to do aything
	for _, poolNode := range s.Me.Pool.nodes {
		poolNames[poolNode.Name] = poolNode.Port
	}

	if _, ok := poolNames[node.Name]; ok { //&& val == int(node.Port) saving this for now
		return true
	}

	return false
}

// InformPoolOfNodes informs the entire pool of changes to its pool
func (s *Server) InformPoolOfNodes() error {
	// Give each of my nodes my phonebook, minus themselves
	if len(s.Me.Pool.nodes) > 0 {
		log.Printf("[network] Informing my bootnodes of their graph...\n")
		for i := range s.Me.Pool.nodes {
			log.Printf("[network] Sending %s my phonebook...\n", s.Me.Pool.nodes[i].Name)
			ierr := s.informNode(s.Me.Pool.nodes[i])
			if ierr != nil {
				log.Printf("[network] Informing failed: %v\n", ierr)
				return ierr
			}
			log.Printf("[network] %s and I have swapped phonebooks.\n", s.Me.Pool.nodes[i].Name)
		}
	}

	return nil
}

// CheckPoolForNodeByName checks my pool for nodes by name
func (s *Server) CheckPoolForNodeByName(name string) (*NodeMessage, bool) {
	for _, node := range s.Me.Pool.nodes {
		if node.Name == name {
			return s.NodeToMessage(node), true
		}
	}
	return &NodeMessage{}, false
}

// AddNodeToPool checks for a nodes existence in its pool and adds if it doesn't
func (s *Server) AddNodeToPool(node *NodeMessage) error {
	// need to think on this, it feels wrong
	if !s.CheckPoolForNode(node) {
		log.Printf("[%s@%s:%d] is new to me. Adding to my phonebook.\n", node.Name, node.Host, node.Port)

		s.Me.Pool.nodes = append(s.Me.Pool.nodes, s.MessageToNode(node))
	}

	return nil
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

		informPayload.Pool = append(informPayload.Pool, s.NodeToMessage(poolNode))
	}

	return informPayload
}
