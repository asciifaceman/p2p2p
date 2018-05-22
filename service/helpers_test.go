package service

import (
	"fmt"
	"testing"
)

func TestAddNodeToPool(t *testing.T) {
	s := &Server{Host: "0.0.0.0", Name: "Tester", Port: 1000}
	s.Me = &Node{Name: "Tester", Host: "0.0.0.0", Port: 1000}

	addNode := &NodeMessage{Name: "Test01", Host: "0.0.0.0", Port: 1000}

	if len(s.Me.Pool.nodes) != 0 {
		t.Fatalf("Expected 0 nodes in pool but got %d. - %v", len(s.Me.Pool.nodes), s.Me.Pool.nodes)
	}

	if err := s.AddNodeToPool(addNode); err != nil {
		t.Fatalf("%v", err)
	}

	if len(s.Me.Pool.nodes) != 1 {
		t.Fatalf("Expected 1 nodes in pool but got %d. - %v", len(s.Me.Pool.nodes), s.Me.Pool.nodes)
	}

	addNode2 := &NodeMessage{Name: "Test02", Host: "0.0.0.0", Port: 1000}

	if err := s.AddNodeToPool(addNode2); err != nil {
		t.Fatalf("%v", err)
	}

	if len(s.Me.Pool.nodes) != 2 {
		t.Fatalf("Expected 2 nodes in pool but got %d. - %v", len(s.Me.Pool.nodes), s.Me.Pool.nodes)
	}

	if s.Me.Pool.nodes[0].Name != addNode.Name {
		t.Fatalf("Expected node at index 0 to have name %s, but got %v - %v", addNode.Name, s.Me.Pool.nodes[0].Name, s.Me.Pool.nodes[0])
	}

}

func TestBuildInformPool(t *testing.T) {
	s := &Server{Host: "0.0.0.0", Name: "Tester", Port: 1000}
	s.Me = &Node{Name: "Tester", Host: "0.0.0.0", Port: 1000}
	//&NodeMessage{Name: "Test01", Host: "0.0.0.0", Port: 1000}
	testPool := [4]NodeMessage{}
	for i := range testPool {
		name := fmt.Sprintf("Test0%d", i)
		host := fmt.Sprintf("%[1]d.%[1]d.%[1]d.%[1]d", i)
		s.AddNodeToPool(&NodeMessage{Name: name, Host: host, Port: int32(i)})
	}

	informPool := s.BuildInformPool("Test00", 0)

	if informPool.Informer.Name != s.Name {
		t.Fatalf("Expected pool informer to be %s but got %v - %v", s.Name, informPool.Informer.Name, informPool.Informer)
	}

	if len(informPool.Pool) != 3 {
		t.Fatalf("Expected pool to contain 3 nodes, received %d - %v", len(informPool.Pool), informPool.Pool)
	}

	if informPool.Pool[0].Port != 1 {
		t.Fatalf("Expected first node in Pool to have port 1, had %v - %v", informPool.Pool[0].Port, informPool.Pool[0])
	}

}
