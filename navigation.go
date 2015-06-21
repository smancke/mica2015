package main

import (
	//"log"
)

func findNearestFieldToDiscover(startField *field, startDirection Direction) (path *NavigationPath) {

	pq := NewNavigationPathPQ()
	pq.Put(&NavigationPath{start: startField, end: startField, moves: []Direction{}})
	//log.Printf("start loop")

	visitedPositions := make(map[string]bool)

	for nextField := pq.Get(); pq.Len() >= 0; nextField = pq.Get() {

		visitedPositions[nextField.end.pos.String()] = true
		
		//log.Printf("next field is: %v -> %v ", nextField.moves, nextField.end.pos)
		if ! nextField.end.buttonStatusKnown {
			//log.Printf("buttonStatusKnown=false, .. so returngin")
			return nextField
		}
		for _, direction := range []Direction{NORTH, EAST, WEST, SOUTH} {
			neighbour := nextField.end.beside[direction]
			if neighbour == nil {
				//log.Printf("neighbour = nil, so .. returngin")

				return nextField
			}
			//log.Printf("%v, neighbour.isWall=%v", direction, neighbour.isWall)
			
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
 	return nil	
}
