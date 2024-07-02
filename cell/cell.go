package cell

type Coords struct {
	X, Y int
}

type Cell struct {
	isAlive bool
	age     int
	coords  Coords
}

func New(isAlive bool, age int, coords Coords) Cell {
	return Cell{
		isAlive: isAlive,
		age:     age,
		coords:  coords,
	}
}

func Status(cell Cell) bool {
	return cell.isAlive
}

func Age(cell Cell) int {
	return cell.age
}

func increaseAge(cell Cell) Cell {
	return New(true, cell.age+1, cell.coords)
}

func kill(cell Cell) Cell {
	return New(false, 0, cell.coords)
}

func spawn(cell Cell) Cell {
	return New(true, 1, cell.coords)
}

func NextGeneration(cell Cell, cells [][]Cell) Cell {
	aliveNeighbors := 0
	coords := cell.coords

	for i := coords.X - 1; i <= coords.X+1; i++ {
		for j := coords.Y - 1; j <= coords.Y+1; j++ {
			if i == coords.X && j == coords.Y {
				continue
			}

			if i < 0 || j < 0 || i >= len(cells) || j >= len(cells[i]) {
				continue
			}

			if Status(cells[i][j]) {
				aliveNeighbors++
			}
		}
	}

	if aliveNeighbors < 2 || aliveNeighbors > 3 {
		return kill(cell)
	}

	if cell.isAlive {
		return increaseAge(cell)
	}

	if aliveNeighbors == 3 {
		return spawn(cell)
	} else {
		return cell
	}
}
