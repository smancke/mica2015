package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFindNearestFieldToDiscover(t *testing.T) {

	maze := NewMaze(&Client{})
	looks := [5]LookDescription{
		LookDescription{},
		LookDescription{left: true},
		LookDescription{right: true},
		LookDescription{isWall: true},
	}	
	maze.updateMaze(looks)

	maze.robotPosition.buttonStatusKnown = true
	maze.getFieldByPosition(Position{1,0}).isWall = true
	maze.getFieldByPosition(Position{-1,0}).isWall = true
	maze.getFieldByPosition(Position{0,-1}).isWall = true
	maze.getFieldByPosition(Position{1,0}).buttonStatusKnown = true
	maze.getFieldByPosition(Position{-1,0}).buttonStatusKnown = true
	maze.getFieldByPosition(Position{0,-1}).buttonStatusKnown = true
	maze.robotPosition.beside[WEST] = maze.getFieldByPosition(Position{1,0});	
	maze.robotPosition.beside[EAST] = maze.getFieldByPosition(Position{-1,0});	
	maze.robotPosition.beside[SOUTH] = maze.getFieldByPosition(Position{0,-1});	
	
	navigationPath := findNearestFieldToDiscover(maze.robotPosition, maze.robotDirection)
	assert.NotNil(t, navigationPath)

	assert.Equal(t, maze.getFieldByPosition(Position{0,0}), navigationPath.start)
	assert.Equal(t, maze.getFieldByPosition(Position{-1,2}), navigationPath.end)

	moves := navigationPath.moves

	assert.Equal(t, 3, len(moves)) 
	assert.Equal(t, NORTH, moves[0])
	assert.Equal(t, NORTH, moves[1])
	assert.Equal(t, WEST, moves[2])

	assert.Equal(t, 4, navigationPath.Cost())
	maze.enablePlot = true
	maze.plot()
}
