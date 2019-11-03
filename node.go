package nametree

import (
	"strings"
	"sync"
)

//New creates a root node without a parent
func New() INode {
	return &node{
		parent: nil,
		name:   "",
		sub:    make(map[string]*node),
	}
}

//INode in the name tree
type INode interface {
	Name() string                         //write the full name e.g. "/" or "/a/b/c"
	Named(name string, create bool) INode //returns the named sub item, creating a new one if necessary when create==true
	Parent() INode
	Root() INode
}

//node implements INode
type node struct {
	mutex  sync.Mutex
	parent *node
	name   string
	sub    map[string]*node
}

func (n *node) Parent() INode {
	if n == nil {
		return nil
	}
	return n.parent
}

func (n *node) Root() INode {
	if n.parent == nil {
		return n
	}
	return n.parent.Root()
}

func (n *node) Name() string {
	if n.parent == nil {
		return "/"
	}
	if n.parent.Name() == "/" {
		return "/" + n.name
	}
	return n.parent.Name() + "/" + n.name
}

func (n *node) Named(name string, create bool) INode {
	return n.named(name, create)
}

func (n *node) named(name string, create bool) *node {
	if n == nil {
		return nil
	}
	if name == "" {
		return n
	}
	parts := strings.SplitN(name, "/", 2)
	if len(parts) < 1 {
		return n
	}

	//>=1 parts
	sub, ok := n.sub[parts[0]]
	if !ok {
		//does not exist
		if !create {
			return nil
		}
		sub = &node{
			parent: n,
			name:   parts[0],
			sub:    make(map[string]*node),
		}
		n.sub[parts[0]] = sub
	}

	//found
	if len(parts) == 1 {
		return sub
	}

	return sub.named(parts[1], create)
}
