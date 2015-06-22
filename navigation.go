package main

import (
	"log"
)

func findPathTo(startField *field, startDirection Direction, endField *field) (path *NavigationPath) {
	log.Printf("findPathTo %v", endField.pos)

	targetCiterium := func(possibleTarget *field) bool {
		return possibleTarget.pos.String() == endField.pos.String()
	}
	return findPath(startField, startDirection, targetCiterium)
}

func findNearestFieldToDiscover(startField *field, startDirection Direction) (path *NavigationPath) {
	targetCiterium := func(possibleTarget *field) bool {
		if ! possibleTarget.buttonStatusKnown {
			return true
		}
		for _, direction := range []Direction{NORTH, EAST, WEST, SOUTH} {
			if possibleTarget.beside[direction] == nil {
				return true
			}
		}
		return false;
	}
	return findPath(startField, startDirection, targetCiterium)
}

func findPath(startField *field, startDirection Direction, targetCiterium (func(*field) bool)) (path *NavigationPath) {

	pq := NewNavigationPathPQ()
	pq.Put(&NavigationPath{start: startField, end: startField, moves: []Direction{}})

	visitedPositions := make(map[string]bool)

	for pq.Len() > 0  {
		nextField := pq.Get()
		visitedPositions[nextField.end.pos.String()] = true

		log.Printf("next field is: %v -> %v ", nextField.moves, nextField.end.pos)
		
		if targetCiterium(nextField.end) {
			return nextField
		}
		for _, direction := range []Direction{NORTH, EAST, WEST, SOUTH} {
			neighbour := nextField.end.beside[direction]

			if neighbour != nil {
				_, visited := visitedPositions[neighbour.pos.String()]				
				if (! visited && ! neighbour.isWall && ! pq.ContainsPathTo(neighbour)) {
					newPath := *nextField
					newPath.end = neighbour
					newPath.moves = make([]Direction, len(nextField.moves))
					copy(newPath.moves, nextField.moves)
					newPath.moves = append(newPath.moves, direction)
					pq.Put(&newPath)
				}
			}
		}
	}
 	return nil	
}
