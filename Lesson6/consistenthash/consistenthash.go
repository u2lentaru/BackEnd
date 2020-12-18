package consistenthash

import (
	"errors"
	"hash/crc32"
	"sort"
	"sync"
)

//ErrNodeNotFound - node not found
var ErrNodeNotFound = errors.New("node not found")

//Ring struct
type Ring struct {
	Nodes Nodes
	sync.Mutex
}

//NewRing func
func NewRing() *Ring {
	return &Ring{Nodes: Nodes{}}
}

//AddNode func
func (r *Ring) AddNode(id string) {
	r.Lock()
	defer r.Unlock()

	node := NewNode(id)
	r.Nodes = append(r.Nodes, node)

	sort.Sort(r.Nodes)
}

//RemoveNode func
func (r *Ring) RemoveNode(id string) error {
	r.Lock()
	defer r.Unlock()

	i := r.search(id)
	if i >= r.Nodes.Len() || r.Nodes[i].ID != id {
		return ErrNodeNotFound
	}

	r.Nodes = append(r.Nodes[:i], r.Nodes[i+1:]...)

	return nil
}

//Get func
func (r *Ring) Get(id string) string {
	i := r.search(id)
	if i >= r.Nodes.Len() {
		i = 0
	}

	return r.Nodes[i].ID
}

func (r *Ring) search(id string) int {
	searchfn := func(i int) bool {
		return r.Nodes[i].HashID >= hashID(id)
	}

	return sort.Search(r.Nodes.Len(), searchfn)
}

//Node struct
type Node struct {
	ID     string
	HashID uint32
}

//NewNode func
func NewNode(id string) *Node {
	return &Node{
		ID:     id,
		HashID: hashID(id),
	}
}

//Nodes type
type Nodes []*Node

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Nodes) Less(i, j int) bool { return n[i].HashID < n[j].HashID }

func hashID(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
