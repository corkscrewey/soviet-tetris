package main

import "github.com/rusq/dagproc"

var _ dagproc.Node = node{}

type node struct {
	id  string
	par string
	do  func() error
}

func (n node) ID() string {
	return n.id
}

func (n node) ParentIDs() []string {
	return []string{n.par}
}

func (n node) Do() error {
	return n.do()
}

func dagFunc(id string, parID string, do func() error) node {
	return node{
		id:  id,
		par: parID,
		do:  do,
	}
}
