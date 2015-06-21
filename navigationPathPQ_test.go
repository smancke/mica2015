package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue(t *testing.T) {

	pq := NewNavigationPathPQ()

	pq.Put( createTestPath(3) )
	pq.Put( createTestPath(2) )
	pq.Put( createTestPath(1) )
	pq.Put( createTestPath(5) )
	pq.Put( createTestPath(4) )

	assert.Equal(t, len(*pq), 5)

	assert.Equal(t, 1, len(pq.Get().moves))
	assert.Equal(t, 2, len(pq.Get().moves))
	assert.Equal(t, 3, len(pq.Get().moves))
	assert.Equal(t, 4, len(pq.Get().moves))
	assert.Equal(t, 5, len(pq.Get().moves))
}

func createTestPath(length int) (path *NavigationPath) {
	path = &NavigationPath{
		nil, nil,
		[]Direction{
		},
	}
	for ; length > 0; length-- {
		path.moves =  append(path.moves, NORTH)
	}
	return
	
}
