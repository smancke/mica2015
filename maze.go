package main

import (
	"fmt"
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
	enablePlot bool
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
	for _, direction := range []Direction{NORTH, EAST, WEST, SOUTH} {
		if f.beside[direction] == nil {
			posBeside := f.pos.next(direction)
			f.beside[direction] = this.fields[posBeside.String()]
		}
	}

	return
}

func (this *maze) FindButtons() {
	for this.buttonToCollect < 10 {
		nextButtonField := this.buttonFields[this.buttonToCollect]
		
		if nextButtonField != nil {
			log.Printf("next button: %v -- going to take it at: %v", this.buttonToCollect, nextButtonField)
			this.goTo(nextButtonField)
			this.plot()
			button := this.client.Push()
			if button == this.buttonToCollect {
				this.buttonToCollect++
			} else {
				log.Fatalf("button %v expected, but found %v at %v", this.buttonToCollect, button, this.robotPosition)
			}
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
		if this.shouldTakeALook(this.robotDirection) {
			this.look()
			this.plot()
		}
		if this.shouldTakeALook(this.left()) {
			this.turnLeft()
			this.plot()
		} else {
			if this.shouldTakeALookAtAnyNeighbour() {
				this.turnRight()
				this.plot()
			}
		}
		if this.buttonFields[this.buttonToCollect] != nil {
			return
		}
	}

	navigationPath := findNearestFieldToDiscover(this.robotPosition, this.robotDirection)
	log.Printf("NearestFieldToDiscover %v", navigationPath)
	if navigationPath != nil {
		this.doMoves(navigationPath)
		this.plot()
	} else {
		log.Printf("No more Fields to discover")
	}
}

func (this *maze) turnTo(direction Direction) {
	if direction == this.robotDirection.left() {
		this.turnLeft()
		return
	}
	
	for direction != this.robotDirection {
		this.turnRight()
	}
}

func (this *maze) doMoves(moves *NavigationPath) {
	for _,direction := range moves.moves {
		this.turnTo(direction)
		this.client.Walk()
		
		this.robotPosition = this.robotPosition.beside[this.robotDirection]
	}
}
 	// don't forget to look on buttons on every step


func (this *maze) goTo(f *field) {
	navigationPath := findPathTo(this.robotPosition, this.robotDirection, f)
	if navigationPath != nil {
		this.doMoves(navigationPath)
	} else {
		log.Printf("No path found from %v to %v", this.robotPosition, f)
	}
}

func (this *maze) turnRight() {
	this.client.Right();
	this.robotDirection = this.robotDirection.right()
}

func (this *maze) turnLeft() {
	this.client.Left();
	this.robotDirection = this.robotDirection.left()
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
		this.buttonFields[button] = this.robotPosition
	}	
}
	
func (this *maze) shouldTakeALook(d Direction) bool {
	steps := 0
	beside := this.robotPosition.beside[d]
	for steps < 2  {
		if (beside == nil || ! beside.wallStatusKnown) {
			return true
		}
		if beside.isWall {
			return false
		}
		beside = beside.beside[d]
		steps++
	}
	return false
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
			current.beside[this.left()].beside[this.right()] = current
			current.beside[this.left()].isWall = ! look.left
			current.beside[this.left()].wallStatusKnown = true
		}

		if current.beside[this.right()] == nil {
			current.beside[this.right()] = this.getFieldByPosition( current.pos.next(this.right()) )
			current.beside[this.right()].beside[this.left()] = current
			current.beside[this.right()].isWall = ! look.right
			current.beside[this.right()].wallStatusKnown = true
		}

		prev = current
	}
}

func (this *maze) plot() {
	if ! this.enablePlot {
		return 
	}

	waitForAnyKey()
	tm.MoveCursor(1,1)
	tm.Clear()
	if (this.buttonToCollect < 10) {
		nextButtonField := this.buttonFields[this.buttonToCollect]
		if nextButtonField == nil {
			fmt.Printf("\nMaze: DISCOVER %v\n", this.buttonToCollect)
		} else {
			fmt.Printf("\nMaze: FETCH    %v\n", this.buttonToCollect)
		}
	}
	
	size := 20
	pos := Position{size,size}
	for ; pos.NORTH >= -1*size; pos.NORTH-- {
		line := ""
		for pos.EAST=-1*size; pos.EAST <= size; pos.EAST++ {
			
			field := this.getFieldByPosition( pos )
			if this.robotPosition.pos.NORTH == pos.NORTH && this.robotPosition.pos.EAST == pos.EAST {
				if this.robotDirection == NORTH {
					line += "^"
				} else if this.robotDirection == EAST {
					line += ">"
				} else if this.robotDirection == SOUTH {
					line += "!"
				} else if this.robotDirection == WEST {
					line += "<"
				}
				
			} else {
				if ! field.buttonStatusKnown && ! field.wallStatusKnown {
					line += "?"
				} else {
					if field.isWall {
						line += "#"
					} else if field.buttonId >= 0 {
						line += fmt.Sprintf("%d", field.buttonId)
					} else {
						line += " "
					}
				}
			}
		}
		fmt.Println(line)
	}
	//tm.Flush()
}

