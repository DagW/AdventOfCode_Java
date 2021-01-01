package main

import (
	"fmt"
	"math"
)

const (
	UNKNOWN = 0
	VISITED = 1
	WALL    = 2
	OXYGEN  = 3
)
const (
	HIT_A_WALL = 0
	SUCCESS    = 1
	FOUND      = 2
)
const (
	NORTH = 1
	SOUTH = 2
	WEST  = 3
	EAST  = 4
)

type position struct {
	y    int
	x    int
	from int
}

func mapWorld(computer *IntCodeComputer, world *map[position]int, current position) (bool, position) {

	tries := [][3]int{{NORTH, -1, 0}, {SOUTH, 1, 0}, {WEST, 0, -1}, {EAST, 0, 1}}
	for _, testpos := range tries {
		modpos := position{
			y:    current.y + testpos[1],
			x:    current.x + testpos[2],
			from: testpos[0],
		}
		if (*world)[modpos] == UNKNOWN {
			computer.run([]int{modpos.from})
			result := computer.outputs[0]
			computer.clearOutputs()
			switch result {
			case HIT_A_WALL:
				//Hit a wall, position not updated, Mark as wall to not try again
				(*world)[modpos] = WALL
			case SUCCESS:
				//Success, run again from the new position
				(*world)[modpos] = VISITED
				found, pos := mapWorld(computer, world, modpos)
				if found {
					return found, pos
				}
			case FOUND:
				//Found oxygen, return the position
				(*world)[modpos] = OXYGEN
				mapWorld(computer, world, modpos)
				return true, modpos

			}
		}
	}
	// Backtrack here
	switch current.from {
	case NORTH:
		computer.run([]int{SOUTH})
		computer.clearOutputs()
	case SOUTH:
		computer.run([]int{NORTH})
		computer.clearOutputs()
	case WEST:
		computer.run([]int{EAST})
		computer.clearOutputs()
	case EAST:
		computer.run([]int{WEST})
		computer.clearOutputs()
	}

	return false, current
}

func bfsSearch(world map[position]int, from position, to position) int {
	searchStack := []position{from}
	distances := map[position]int{from: 0}
	for len(searchStack) > 0 {
		current := searchStack[len(searchStack)-1]
		searchStack = searchStack[0 : len(searchStack)-1]
		neighbours := []position{
			{y: current.y + 1, x: current.x},
			{y: current.y - 1, x: current.x},
			{y: current.y, x: current.x + 1},
			{y: current.y, x: current.x - 1},
		}
		for _, neighbour := range neighbours {
			locationtype, inmap := world[neighbour]
			if inmap && locationtype != WALL {
				_, visited := distances[neighbour]
				if !visited {
					searchStack = append(searchStack, neighbour)
				}
				if distances[neighbour] == 0 || distances[neighbour] > distances[current]+1 {
					distances[neighbour] = distances[current] + 1
				}
			}
		}
	}
	return distances[to]
}

func printWorld(world map[position]int, pos position) {
	var minY, minX = math.MaxInt32, math.MaxInt32
	var maxY, maxX = math.MinInt32, math.MinInt32
	for pos := range world {
		if pos.x < minX {
			minX = pos.x
		}
		if pos.x > maxX {
			maxX = pos.x
		}
		if pos.y < minY {
			minY = pos.y
		}
		if pos.y > maxY {
			maxY = pos.y
		}
	}
	for x := minX - 4; x < maxX+4; x++ {
		fmt.Print("-")
	}
	fmt.Println()
	for y := minY - 3; y < maxY+3; y++ {
		fmt.Print("|")
		for x := minX - 3; x < maxX+3; x++ {
			if pos.x == x && pos.y == y && pos.from != 0 {
				fmt.Print("X")
			} else {
				found := false
				val := -1
				for k, v := range world {
					if k.y == y && k.x == x {
						val = v
						found = true
					}
				}
				if found {
					switch val {
					case VISITED:
						fmt.Print(".")
					case WALL:
						fmt.Print("#")
					case OXYGEN:
						fmt.Print("^")
					}
				} else {
					fmt.Print(" ")
				}

			}
		}
		fmt.Print("|")
		fmt.Println()
	}
	for x := minX - 4; x < maxX+4; x++ {
		fmt.Print("-")
	}
	fmt.Println()
}

func main() {
	program := []int{3, 1033, 1008, 1033, 1, 1032, 1005, 1032, 31, 1008, 1033, 2, 1032, 1005, 1032, 58, 1008, 1033, 3, 1032, 1005, 1032, 81, 1008, 1033, 4, 1032, 1005, 1032, 104, 99, 1002, 1034, 1, 1039, 1002, 1036, 1, 1041, 1001, 1035, -1, 1040, 1008, 1038, 0, 1043, 102, -1, 1043, 1032, 1, 1037, 1032, 1042, 1106, 0, 124, 1001, 1034, 0, 1039, 1002, 1036, 1, 1041, 1001, 1035, 1, 1040, 1008, 1038, 0, 1043, 1, 1037, 1038, 1042, 1106, 0, 124, 1001, 1034, -1, 1039, 1008, 1036, 0, 1041, 1002, 1035, 1, 1040, 1001, 1038, 0, 1043, 101, 0, 1037, 1042, 1105, 1, 124, 1001, 1034, 1, 1039, 1008, 1036, 0, 1041, 102, 1, 1035, 1040, 1001, 1038, 0, 1043, 101, 0, 1037, 1042, 1006, 1039, 217, 1006, 1040, 217, 1008, 1039, 40, 1032, 1005, 1032, 217, 1008, 1040, 40, 1032, 1005, 1032, 217, 1008, 1039, 39, 1032, 1006, 1032, 165, 1008, 1040, 3, 1032, 1006, 1032, 165, 1102, 1, 2, 1044, 1106, 0, 224, 2, 1041, 1043, 1032, 1006, 1032, 179, 1102, 1, 1, 1044, 1106, 0, 224, 1, 1041, 1043, 1032, 1006, 1032, 217, 1, 1042, 1043, 1032, 1001, 1032, -1, 1032, 1002, 1032, 39, 1032, 1, 1032, 1039, 1032, 101, -1, 1032, 1032, 101, 252, 1032, 211, 1007, 0, 59, 1044, 1105, 1, 224, 1102, 1, 0, 1044, 1105, 1, 224, 1006, 1044, 247, 101, 0, 1039, 1034, 1001, 1040, 0, 1035, 101, 0, 1041, 1036, 1002, 1043, 1, 1038, 1002, 1042, 1, 1037, 4, 1044, 1105, 1, 0, 93, 27, 71, 56, 88, 17, 30, 78, 5, 57, 79, 56, 3, 82, 62, 58, 16, 2, 21, 89, 95, 33, 12, 32, 90, 12, 7, 76, 83, 31, 8, 13, 27, 89, 60, 33, 7, 40, 22, 50, 8, 63, 35, 45, 57, 94, 81, 4, 65, 33, 47, 73, 28, 98, 11, 70, 95, 17, 82, 39, 19, 73, 62, 56, 80, 85, 23, 91, 39, 86, 91, 82, 50, 37, 86, 4, 90, 83, 8, 65, 56, 63, 15, 99, 51, 3, 60, 60, 77, 58, 90, 82, 5, 52, 14, 87, 37, 74, 85, 43, 17, 61, 91, 35, 31, 81, 19, 12, 34, 54, 9, 66, 34, 69, 67, 21, 4, 14, 87, 22, 76, 26, 82, 79, 4, 69, 48, 73, 8, 73, 57, 61, 83, 23, 83, 60, 3, 41, 75, 67, 53, 44, 91, 27, 52, 84, 66, 13, 65, 95, 81, 83, 30, 26, 60, 12, 33, 92, 81, 46, 78, 25, 13, 72, 87, 26, 63, 57, 35, 2, 60, 96, 63, 26, 2, 76, 95, 21, 38, 60, 5, 79, 86, 89, 47, 42, 12, 91, 30, 52, 69, 55, 67, 73, 47, 44, 5, 86, 8, 52, 69, 81, 23, 70, 3, 38, 41, 89, 88, 58, 41, 9, 96, 27, 67, 21, 14, 68, 67, 35, 84, 23, 20, 91, 63, 47, 75, 34, 70, 57, 13, 54, 82, 33, 61, 27, 97, 88, 46, 44, 56, 74, 14, 5, 96, 71, 16, 40, 86, 61, 84, 41, 81, 81, 16, 88, 51, 41, 96, 76, 28, 97, 44, 41, 65, 87, 50, 73, 58, 71, 46, 73, 51, 43, 18, 46, 99, 74, 65, 9, 89, 3, 77, 22, 34, 93, 94, 39, 54, 96, 12, 35, 62, 87, 56, 69, 64, 9, 34, 91, 64, 71, 28, 10, 94, 1, 96, 20, 67, 92, 39, 37, 26, 79, 68, 16, 76, 57, 83, 92, 46, 75, 99, 26, 64, 39, 72, 65, 37, 93, 65, 5, 53, 62, 36, 13, 97, 14, 38, 85, 33, 76, 56, 99, 29, 64, 84, 28, 19, 91, 92, 55, 33, 88, 32, 70, 38, 53, 76, 1, 76, 35, 26, 75, 18, 18, 7, 88, 19, 53, 65, 22, 91, 20, 85, 15, 13, 72, 82, 13, 31, 75, 62, 68, 4, 56, 91, 89, 56, 10, 46, 63, 7, 74, 50, 15, 85, 87, 64, 77, 12, 95, 10, 66, 77, 51, 6, 61, 75, 91, 75, 85, 61, 78, 4, 97, 99, 4, 90, 34, 89, 44, 44, 68, 89, 30, 20, 70, 24, 22, 81, 22, 77, 61, 33, 89, 2, 11, 75, 50, 85, 13, 43, 56, 78, 73, 49, 27, 38, 78, 56, 90, 17, 94, 72, 51, 5, 55, 67, 32, 19, 81, 81, 45, 83, 18, 96, 33, 75, 53, 4, 29, 87, 80, 33, 57, 78, 80, 43, 68, 57, 71, 83, 10, 18, 98, 70, 36, 61, 31, 73, 33, 69, 24, 78, 76, 43, 88, 96, 16, 14, 91, 43, 66, 15, 98, 44, 48, 68, 57, 72, 48, 49, 89, 62, 31, 55, 83, 68, 86, 97, 16, 25, 87, 13, 74, 40, 82, 43, 48, 85, 40, 45, 72, 33, 60, 84, 4, 47, 96, 19, 92, 75, 73, 46, 6, 69, 4, 81, 98, 89, 48, 55, 89, 24, 64, 31, 47, 50, 93, 72, 47, 72, 36, 79, 7, 24, 66, 60, 65, 18, 81, 93, 40, 37, 36, 62, 94, 48, 8, 77, 21, 82, 22, 65, 20, 46, 85, 47, 52, 70, 55, 74, 19, 65, 15, 72, 81, 57, 67, 46, 94, 21, 16, 94, 84, 36, 43, 62, 82, 48, 47, 79, 5, 96, 39, 58, 85, 80, 31, 7, 98, 23, 69, 22, 99, 37, 69, 35, 66, 36, 70, 3, 69, 47, 6, 64, 38, 69, 42, 57, 91, 89, 21, 89, 13, 42, 78, 24, 44, 79, 74, 65, 63, 85, 10, 50, 71, 94, 26, 78, 55, 5, 26, 71, 46, 20, 83, 96, 51, 87, 2, 99, 83, 5, 38, 86, 8, 13, 94, 61, 93, 39, 67, 23, 60, 74, 87, 57, 30, 72, 23, 19, 95, 57, 93, 83, 58, 34, 83, 35, 4, 47, 81, 88, 24, 87, 34, 93, 79, 70, 18, 24, 73, 98, 76, 77, 24, 93, 18, 66, 56, 87, 25, 29, 7, 7, 97, 40, 61, 56, 96, 96, 1, 42, 21, 92, 73, 11, 10, 97, 69, 58, 93, 2, 82, 27, 96, 7, 84, 44, 67, 57, 63, 13, 79, 56, 72, 34, 89, 26, 94, 24, 86, 99, 71, 73, 98, 26, 89, 10, 98, 5, 64, 70, 85, 32, 61, 35, 67, 0, 0, 21, 21, 1, 10, 1, 0, 0, 0, 0, 0, 0}
	world := map[position]int{}
	computer := IntCodeComputer{}
	computer.program(program)
	start := position{}

	_, target := mapWorld(&computer, &world, start)

	// Have to clear the direction from all positions
	worldWithoutDirections := map[position]int{}
	for k, v := range world {
		worldWithoutDirections[position{y: k.y, x: k.x}] = v
	}
	target.from = 0

	printWorld(worldWithoutDirections, start)
	fmt.Println(bfsSearch(worldWithoutDirections, start, target))
	fmt.Println(oxygenate(worldWithoutDirections))
}

func oxygenate(world map[position]int) int {
	minutes := 1
	visitedPositions := map[position]bool{}
	for {
		for pos, v := range world {
			visited := visitedPositions[pos]
			if !visited && v == OXYGEN {
				neighbours := []position{
					{y: pos.y + 1, x: pos.x},
					{y: pos.y - 1, x: pos.x},
					{y: pos.y, x: pos.x + 1},
					{y: pos.y, x: pos.x - 1},
				}
				for _, neighbour := range neighbours {
					value, existsInMap := world[neighbour]
					if existsInMap && value == VISITED {
						world[neighbour] = OXYGEN
					}
				}
				visitedPositions[pos] = true
			}
		}
		numNotOxygenated := 0
		for _, v := range world {
			if v == VISITED {
				numNotOxygenated++
			}
		}
		fmt.Println("num:", numNotOxygenated)
		minutes++
		if numNotOxygenated == 0 {
			break
		}
	}
	printWorld(world, position{})
	return minutes
}
