package main

import (
	"log"
)

func findPathTo(startField *field, startDirection Direction, endField *field) (path *NavigationPath) {
	log.Printf("findPathTo %v", endField.pos)

	targetCiterium := func(possibleTarget *field) bool {
		return possibleTarget.pos.String() == endField.pos.String()
	}
	return findPath(startField, startDirection, targetCiterium, false)
}

func findNearestFieldToDiscover(startField *field, startDirection Direction) (path *NavigationPath) {
	targetCiterium := func(possibleTarget *field) bool {
		if startField.pos.String() == possibleTarget.pos.String() {
			return false
		}
		if ! possibleTarget.buttonStatusKnown {
			return true
		}
		unknownBesides := 0
		for _, direction := range []Direction{NORTH, EAST, WEST, SOUTH} {			
			if possibleTarget.beside[direction] == nil {
				return true
			}
			if possibleTarget.beside[direction].wallStatusKnown == false || possibleTarget.beside[direction].buttonStatusKnown == false {
				unknownBesides++
			}
		}
		if unknownBesides >= 3 {
			return true
		}
		return false
	}
	return findPath(startField, startDirection, targetCiterium, true)
}

func findPath(startField *field, startDirection Direction, targetCiterium (func(*field) bool), withBenefit bool) (path *NavigationPath) {

	pq := NewNavigationPathPQ()
	pq.Put(&NavigationPath{start: startField, startDirection: startDirection, end: startField, moves: []Direction{}, takeBenefitIntoAccount: withBenefit})

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
