// Copyright © 2018 Charles Corebtt <nafredy@gmail.com>
//

package service

// Node defines the structure of a node
type Node struct {
	Host string
	Port int
	Name string
	Pool Pool
}
