package main

import(
	"fmt"
	"container/heap"
)

type NavigationPath struct {
	start *field
	end *field
	moves []Direction
}

func (path *NavigationPath) String() string {
	return fmt.Sprintf("%v->%v->%v", path.start.pos, path.moves, path.end.pos)
}

type NavigationPathPQ []*NavigationPath


func NewNavigationPathPQ() (*NavigationPathPQ) {
	pq := make(NavigationPathPQ, 0)
	return &pq
}

func (pq *NavigationPathPQ) Put(path *NavigationPath)  {
	heap.Push(pq, path)
}

func (pq *NavigationPathPQ) Get() (*NavigationPath) {
	path := heap.Pop(pq);
	return path.(*NavigationPath)
}

func (pq NavigationPathPQ) Len() int {
	return len(pq)
}

func (pq NavigationPathPQ) Less(i, j int) bool {
	return len(pq[i].moves) < len(pq[j].moves)
}

func (pq NavigationPathPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *NavigationPathPQ) Push(path interface{}) {
	item := path.(*NavigationPath)
	*pq = append(*pq, item)
}

func (pq *NavigationPathPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}


func (pq *NavigationPathPQ) ContainsPathTo(targetField *field) bool {
	for _, path := range *pq {
		if path.end.pos == targetField.pos {
			return true
		}
	}
	return false
}
