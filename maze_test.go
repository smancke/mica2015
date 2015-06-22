package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRightAndLeft(t *testing.T) {
	assert.Equal(t, NORTH.right(), EAST)
	assert.Equal(t, NORTH.right().right().right(), NORTH.left())
	assert.Equal(t, NORTH.opposite(), SOUTH)
}

func Test_getFieldByPosition(t *testing.T) {
	maze := NewMaze(&Client{})
	pos := Position{7,6}
	field := maze.getFieldByPosition(pos)
	assert.Equal(t, field.pos, &pos)
}

func TestUpdateMaze_BasicCase(t *testing.T) {
	
	maze := NewMaze(&Client{})
	looks := [5]LookDescription{
		LookDescription{},
		LookDescription{left: true, right: true},
		LookDescription{hasButton: true, buttonId: 4},
		LookDescription{isWall: true},
	}

	maze.updateMaze(looks)

	// step 1
	assert.Equal(t, maze.getFieldByPosition(Position{0,1}), maze.robotPosition.beside[NORTH])
	assert.True(t, maze.getFieldByPosition(Position{0,1}).isWall == false)
	assert.True(t, maze.getFieldByPosition(Position{1,1}).isWall)
	assert.True(t, maze.getFieldByPosition(Position{-1,1}).isWall)

	// step 2
	assert.Equal(t, maze.getFieldByPosition(Position{0,2}), maze.getFieldByPosition(Position{0,1}).beside[NORTH])
	assert.Equal(t, maze.getFieldByPosition(Position{0,2}).beside[SOUTH], maze.getFieldByPosition(Position{0,1}))
	assert.True(t, maze.getFieldByPosition(Position{0,2}).isWall == false)

	assert.Equal(t, maze.getFieldByPosition(Position{0,2}).beside[EAST], maze.getFieldByPosition(Position{1,2}))
	assert.Equal(t, maze.getFieldByPosition(Position{1,2}).beside[WEST], maze.getFieldByPosition(Position{0,2}))
	assert.True(t, maze.getFieldByPosition(Position{1,2}).isWall == false)
	assert.True(t, maze.getFieldByPosition(Position{1,2}).wallStatusKnown)

	assert.Equal(t, maze.getFieldByPosition(Position{0,2}).beside[WEST], maze.getFieldByPosition(Position{-1,2}))
	assert.Equal(t, maze.getFieldByPosition(Position{-1,2}).beside[EAST], maze.getFieldByPosition(Position{0,2}))
	assert.True(t, maze.getFieldByPosition(Position{-1,2}).isWall == false)
	assert.True(t, maze.getFieldByPosition(Position{-1,2}).wallStatusKnown)

	// step 3
	assert.True(t, maze.getFieldByPosition(Position{0,3}).buttonStatusKnown)
	assert.Equal(t, maze.getFieldByPosition(Position{0,3}).buttonId, 4)
	assert.Equal(t, maze.buttonFields[4], maze.getFieldByPosition(Position{0,3}))

	// step 4
	assert.True(t, maze.getFieldByPosition(Position{0,4}).isWall)
}


func TestUpdateMaze_ConnectionToExistingFields(t *testing.T) {
	
	maze := NewMaze(&Client{})
	maze.getFieldByPosition(Position{-1,1}).buttonId = 4
	maze.getFieldByPosition(Position{-1,1}).buttonStatusKnown = true
	
	looks := [5]LookDescription{
		LookDescription{left: true},
		LookDescription{isWall: true},
	}

	maze.updateMaze(looks)

	assert.Equal(t, maze.getFieldByPosition(Position{0,1}).beside[WEST].buttonId, 4)
}
