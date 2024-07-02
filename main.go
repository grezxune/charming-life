package main

import (
	"fmt"
	"os"

	"github.com/grezxune/charm-life/cell"
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

const timeout = time.Hour
const interval = time.Millisecond * 100
const gridSize = 50
const liveChance = 5

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
	next  key.Binding
}

type board struct {
	cells     [][]cell.Cell
	iteration int
	timer     timer.Model
	keymap    keymap
}

func getRandomBool() bool {
	return rand.Intn(liveChance) == 0
}

func initialModel() board {
	return board{
		cells:     createNewCells(),
		iteration: 0,
		timer:     timer.NewWithInterval(timeout, interval),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "quit"),
			),
			next: key.NewBinding(
				key.WithKeys("n"),
				key.WithHelp("n", "next"),
			),
		},
	}
}

func createNewCells() [][]cell.Cell {
	cells := make([][]cell.Cell, gridSize)

	for i := range cells {
		cells[i] = make([]cell.Cell, gridSize)
		for j := range cells[i] {
			isNewCellAlive := getRandomBool()

			if isNewCellAlive {
				cells[i][j] = cell.New(isNewCellAlive, 1, cell.Coords{X: i, Y: j})
			} else {
				cells[i][j] = cell.New(isNewCellAlive, 0, cell.Coords{X: i, Y: j})
			}
		}
	}

	return cells
}

func (board board) Init() tea.Cmd {
	return board.timer.Init()
}

func toggleCells(board *board) {
	// Create a new slice with the same dimensions as the existing one
	newCells := make([][]cell.Cell, len(board.cells))

	// Toggle the isAlive state and assign to the new slice
	for i, row := range board.cells {
		newCells[i] = make([]cell.Cell, len(board.cells[i]))

		for j, item := range row {
			newCells[i][j] = cell.NextGeneration(item, board.cells)
		}
	}

	// Assign the updated cells back to the board
	board.cells = newCells
}

func iterateBoard(board *board) {
	board.iteration++
	toggleCells(board)
}

func (board board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		board.timer, cmd = board.timer.Update(msg)
		iterateBoard(&board)
		return board, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		board.timer, cmd = board.timer.Update(msg)
		board.keymap.stop.SetEnabled(board.timer.Running())
		board.keymap.start.SetEnabled(!board.timer.Running())
		return board, cmd

	case timer.TimeoutMsg:
		board.timer.Timeout = timeout
		return board, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, board.keymap.quit):
			return board, tea.Quit

		case key.Matches(msg, board.keymap.reset):
			board.timer.Timeout = timeout
			board.iteration = 0
			board.cells = createNewCells()

		case key.Matches(msg, board.keymap.start, board.keymap.stop):
			return board, board.timer.Toggle()

		case key.Matches(msg, board.keymap.next):
			var cmd tea.Cmd
			iterateBoard(&board)
			return board, cmd
		}
	}

	return board, nil
}

func (board board) View() string {
	s := "Iteration: " + fmt.Sprint(board.iteration) + "\n\n"

	deadCellStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(2).Height(1).Background(lipgloss.Color("#383838")).Foreground(lipgloss.Color("#ffffff"))
	liveCellStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(2).Height(1).Background(lipgloss.Color("#9DAEFF"))

	liveCell := liveCellStyle.Render()
	deadCell := deadCellStyle.Render()

	for _, row := range board.cells {
		for _, item := range row {
			if cell.Status(item) {
				s += liveCell
			} else {
				s += deadCell
			}
		}
		s += "\n"
	}

	s += "\nPress q to quit, r to reset board, n to manually iterate the board\n"

	textStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(30).Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#000000"))
	test := textStyle.Render("Testing!")

	boardStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9DAEFF")).
		Padding(0).
		Margin(3, 5)

	return boardStyle.Render(s + test)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
