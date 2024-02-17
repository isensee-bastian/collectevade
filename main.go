package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	// Define the size of our playing field, add 2 for the borders.
	tableWidth  = 60 + 2
	tableHeight = 30 + 2

	// Define our symbols for representing the playing field.
	corner         = '+'
	lineVertical   = '|'
	lineHorizontal = '-'
	empty          = ' '
	player         = 'o'
	item           = '$'
)

// model encapsulates our data for displaying and updating.
type model struct {
	table [tableHeight][tableWidth]rune

	playerRow int
	playerCol int

	itemRow int
	itemCol int
}

func (m *model) spawnItem() {
	// Generate some random coordinates.
	row, col := randomCoordinates()

	// Check that the random cell is empty.
	// If not repeat randomizing until we find an empty cell.
	for m.table[row][col] != 0 {
		row, col = randomCoordinates()
	}

	m.table[row][col] = item
	m.itemRow = row
	m.itemCol = col
}

// randomCoordinates returns a new set of random coordinates within the
// playing field excluding borders. However, it is not guaranteed that the
// cell under the returned coordinates is actually empty.
func randomCoordinates() (row int, col int) {
	row = rand.Intn(tableHeight-2) + 1
	col = rand.Intn(tableWidth-2) + 1

	return
}

func (m *model) movePlayer(row, col int) {
	// Clear old player location.
	m.table[m.playerRow][m.playerCol] = 0

	// Set new player location.
	m.table[row][col] = player
	m.playerRow = row
	m.playerCol = col

	// Check if we are at the position of an item. Then we need to
	// spawn a new item.
	if m.itemRow == row && m.itemCol == col {
		m.spawnItem()
	}
}

func (m *model) playerUp() {
	if m.playerRow <= 1 {
		// Do nothing as we are already at the border and cannot move.
		return
	}

	m.movePlayer(m.playerRow-1, m.playerCol)
}

func (m *model) playerDown() {
	if m.playerRow >= tableHeight-2 {
		// Do nothing as we are already at the border and cannot move.
		return
	}

	m.movePlayer(m.playerRow+1, m.playerCol)
}

func (m *model) playerLeft() {
	if m.playerCol <= 1 {
		// Do nothing as we are already at the border and cannot move.
		return
	}

	m.movePlayer(m.playerRow, m.playerCol-1)
}

func (m *model) playerRight() {
	if m.playerCol >= tableWidth-2 {
		// Do nothing as we are already at the border and cannot move.
		return
	}

	m.movePlayer(m.playerRow, m.playerCol+1)
}

// newModel is responsible for creating an initial model that is ready to use.
func newModel() *model {
	// By default, all entries in our table have value zero as type rune
	// is based on int and represents symbols. We only set non-empty fields
	// explicitly.
	model := &model{}

	// Set the four corners.
	model.table[0][0] = corner
	model.table[0][tableWidth-1] = corner
	model.table[tableHeight-1][0] = corner
	model.table[tableHeight-1][tableWidth-1] = corner

	// Draw horizontal borders at the top and bottom.
	for col := 1; col < tableWidth-1; col++ {
		model.table[0][col] = lineHorizontal
		model.table[tableHeight-1][col] = lineHorizontal
	}

	// Draw vertical borders on the left and right side.
	for row := 1; row < tableHeight-1; row++ {
		model.table[row][0] = lineVertical
		model.table[row][tableWidth-1] = lineVertical
	}

	// Spawn our player near the top left corner.
	model.playerRow = 1
	model.playerCol = 1
	model.table[model.playerRow][model.playerCol] = player

	model.spawnItem()

	return model
}

// Init can be used to setup initial command to perform.
// We don't need anything here. Therefore we return nil.
func (m *model) Init() tea.Cmd {
	return nil
}

// Update is called whenever something happens like a key is pressed or
// another event occurs. Then, we have the option of reacting to it by
// modifying our model.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		// Exit program on ctrl+c or q typing.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.playerUp()
		case "down":
			m.playerDown()
		case "left":
			m.playerLeft()
		case "right":
			m.playerRight()
		}
	}

	return m, nil
}

// View is required for building what we want to show on the screen.
// That means we need to translate our model data into a string for displaying.
func (m *model) View() string {
	builder := strings.Builder{}

	// Iterate our table (2d array) and print non-empty fields as set in
	// our model. Empty (zero) fields shall be printed with a blank space.
	for _, row := range m.table {

		for _, cell := range row {
			if cell == 0 {
				builder.WriteRune(empty)
			} else {
				builder.WriteRune(cell)
			}
		}

		// Go to next line after each row.
		builder.WriteString("\n")
	}

	return builder.String()
}

func main() {
	// Create our initial model.
	model := newModel()

	// Program setup to initialize bubbletea and use full screen.
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Run bubbletea and exist with a message if an error occurs.
	if _, err := program.Run(); err != nil {
		fmt.Println("Unexpected error: %v", err)
		os.Exit(1)
	}
}
