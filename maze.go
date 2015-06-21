package main

import (
	"fmt"
	"os"
	"log"
	"math"
	tm "github.com/buger/goterm"
)

type Direction int8

const (
	NORTH Direction = 0
	EAST Direction = 1
	SOUTH Direction = 2
	WEST Direction = 3
)

func (dir Direction) right() Direction {
	return Direction(math.Mod(float64(dir+1), 4))
}

func (dir Direction) opposite() Direction {
	return Direction(math.Mod(float64(dir+2), 4))
}

func (dir Direction) left() Direction {
	return Direction(math.Mod(float64(dir+3), 4))
}

type Position struct {
	EAST int
	NORTH int
}

func (p *Position) String() (string) {
	return fmt.Sprintf("%d,%d", p.EAST, p.NORTH)	
}

func (p *Position) next(direction Direction) (newPosition Position) {
	newPosition.EAST = p.EAST
	newPosition.NORTH = p.NORTH
	if direction == NORTH {
		newPosition.NORTH++
	} else if direction == EAST {
		newPosition.EAST++
	} else if direction == SOUTH {
		newPosition.NORTH--
	} else if direction == WEST {
		newPosition.EAST--
	}
	return
}

type field struct {
	beside map[Direction]*field
	wallStatusKnown bool
	isWall bool
	buttonStatusKnown bool
	buttonId int
	pos *Position
}

func (field *field) String() string {
	icon := "?"
	if field.buttonStatusKnown || field.wallStatusKnown {
		if field.isWall {
			icon = "#"
		} else if field.buttonId >= 0 {
			icon = fmt.Sprintf("%d", field.buttonId)
		} else {
			icon = "_"
		}
	}

	s := fmt.Sprintf("%v: %v", field.pos, icon)
	return s
}

type maze struct {
	client MazeClient
	fields map[string]*field
	buttonFields []*field
	buttonToCollect int
	robotPosition *field
	robotDirection Direction
}

func NewMaze(client MazeClient) (this *maze) {
	this = &maze{}
	this.client = client
	this.fields = make(map[string]*field)
	this.buttonFields = make([]*field, 10)
	this.robotPosition = this.NewField(Position{0,0})
	this.robotDirection = NORTH
	return
}

func (this *maze) NewField(position Position) *field {
	newField := &field{
		wallStatusKnown: false,
		buttonStatusKnown: false,
		buttonId: -1,
		beside: make(map[Direction]*field),
		pos: &position,
	}
	this.fields[position.String()] = newField
	return newField
}

func (this *maze) getFieldByPosition(pos Position) (f *field){	
	f = this.fields[pos.String()]
	if f == nil {
		f = this.NewField( pos )
	}
	return
}

func (this *maze) FindButtons() {
	for this.buttonToCollect < 10 {
		nextButtonField := this.buttonFields[this.buttonToCollect]
		
		if nextButtonField != nil {
			log.Printf("next button: %v -- going to take it at: %v", this.buttonToCollect, nextButtonField)
			this.goTo(nextButtonField)
		} else {
			log.Printf("button %v not found; discover from %v", this.buttonToCollect, this.robotPosition)
			this.discover()
		}
	}
	log.Printf("got all buttons!!!")
}

func (this *maze) discover() {
	
	this.ensureButtonStatusIsKnown()

	for this.shouldTakeALookAtAnyNeighbour() {
		// TODO: optimize, not only to turn right
		if this.shouldTakeALook(this.robotDirection) {
			this.look()
		}
		if this.shouldTakeALookAtAnyNeighbour() {
			this.turnRight()
		}
		if this.buttonFields[this.buttonToCollect] != nil {
			return
		}
	}

	navigationPath := findNearestFieldToDiscover(this.robotPosition, this.robotDirection)
	log.Printf("NearestFieldToDiscover %v", navigationPath)
	if navigationPath != nil {
		this.doMoves(navigationPath)
	}
}

func (this *maze) turnTo(direction Direction) {
	//log.Printf("start turnTo(%v): %v", direction, this.robotDirection)
	//TODO: optimize
	for direction != this.robotDirection {
		//log.Printf("turnTo(%v): %v", direction, this.robotDirection)
		this.turnRight()
	}
}

func (this *maze) doMoves(moves *NavigationPath) {
	for _,direction := range moves.moves {
		this.turnTo(direction)
		this.client.Walk()
		
		this.robotPosition = this.robotPosition.beside[this.robotDirection]
		this.plot()
	}
}
 	// don't forget to look on buttons on every step
// func (this *maze) findPath(fOrigin *field, fTarget *field) *[]move


func (this *maze) goTo(f *field) {
	// todo: implement
	log.Printf("TODO: implement goto %v!!!", f)
	os.Exit(1)
}

func (this *maze) turnRight() {
	this.client.Right();
	this.robotDirection = this.robotDirection.right()
	this.plot()
}

func (this *maze) turnLeft() {
	this.client.Left();
	this.robotDirection = this.robotDirection.left()
	this.plot()
}

func (this *maze) ensureButtonStatusIsKnown() {
	this.robotPosition.isWall = false

	if this.robotPosition.buttonStatusKnown {
		return
	}

	button := this.client.Push()
	this.robotPosition.buttonStatusKnown = true	
	if button > -1 {
		this.robotPosition.buttonId = button
	}	
}
	
func (this *maze) shouldTakeALook(d Direction) bool {
	return this.robotPosition.beside[d] == nil
}

func (this *maze) shouldTakeALookAtAnyNeighbour() bool {
	for _, direction := range []Direction{NORTH, EAST, WEST, SOUTH} {
		if this.shouldTakeALook(direction) {
			return true
		}
	}
	return false;
}

func (this *maze) left() Direction {	
	return this.robotDirection.left()
}

func (this *maze) right() Direction {	
	return this.robotDirection.right()
}

func (this *maze) look() {
	looks := this.client.Look()
	this.updateMaze(looks)
	
	this.plot()
}

func (this *maze) updateMaze(looks [5]LookDescription) {
	prev := this.robotPosition	
	for _,look := range looks {

		// already seen and connected
		current := prev.beside[this.robotDirection]
		if current == nil {

			// already seen but new connection or new
			current = this.getFieldByPosition(  prev.pos.next(this.robotDirection) )
			prev.beside[this.robotDirection] = current
		}

		// doing updates	
		current.beside[this.robotDirection.opposite()] = prev
		current.buttonStatusKnown = true
		current.wallStatusKnown = true

		if look.isWall {
			current.isWall = true
			return
		}

		if look.hasButton {
			current.buttonStatusKnown = true
			current.buttonId = look.buttonId
			this.buttonFields[current.buttonId] = current
		}
		
		if current.beside[this.left()] == nil {
			current.beside[this.left()] = this.getFieldByPosition( current.pos.next(this.left()) )
			current.beside[this.left()].isWall = ! look.left
			current.beside[this.left()].wallStatusKnown = true
		}

		if current.beside[this.right()] == nil {
			current.beside[this.right()] = this.getFieldByPosition( current.pos.next(this.right()) )
			current.beside[this.right()].isWall = ! look.right
			current.beside[this.right()].wallStatusKnown = true
		}

		prev = current
	}
}


func (this *maze) plot() {
	tm.Clear()
	tm.MoveCursor(1,1)
	tm.Print("\n\nMaze:\n")

	size := 20
	pos := Position{size,size}
	for ; pos.NORTH >= -1*size; pos.NORTH-- {
		for pos.EAST=-1*size; pos.EAST <= size; pos.EAST++ {
			
			field := this.getFieldByPosition( pos )
			if this.robotPosition.pos.NORTH == pos.NORTH && this.robotPosition.pos.EAST == pos.EAST {
				if this.robotDirection == NORTH {
					tm.Print("^")
				} else if this.robotDirection == EAST {
					tm.Print(">")
				} else if this.robotDirection == SOUTH {
					tm.Print("!")
				} else if this.robotDirection == WEST {
					tm.Print("<")
				}
				
			} else {
				if ! field.buttonStatusKnown && ! field.wallStatusKnown {
					tm.Print("?")
				} else {
					if field.isWall {
						tm.Printf("#")
					} else if field.buttonId >= 0 {
						tm.Printf("%d", field.buttonId)
					} else {
						tm.Printf(" ")
					}
				}
			}
		}
		tm.Print("\n")
	}
	tm.Flush()
}

