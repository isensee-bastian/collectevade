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
	player         = '0'
	item           = '$'
	enemy          = 'X'
)

// model encapsulates our data for displaying and updating.
type model struct {
	table [tableHeight][tableWidth]rune

	playerRow int
	playerCol int

	score    int
	gameOver bool
}

func (m *model) spawnItem() {
	row, col := m.randomFreeCoordinates()
	m.table[row][col] = item
}

func (m *model) spawnEnemy() {
	row, col := m.randomFreeCoordinates()
	m.table[row][col] = enemy
}

func (m *model) moveEnemies() {
	for row := 1; row < tableHeight-1; row++ {
		for col := 1; col < tableWidth-1; col++ {
			if m.table[row][col] == enemy {
				// Move enemy randomly to one of the four directly neighbooring cells if it is empty
				// or contains the player. Borders, items and other enemies will block enemy moves though.

				neighboors := [4][2]int{
					[2]int{row - 1, col}, // Top neighboor.
					[2]int{row, col + 1}, // Right neighboor.
					[2]int{row + 1, col}, // Bottom neighboor.
					[2]int{row, col - 1}, // Right neighboor.
				}

				targetIndex := rand.Intn(len(neighboors))
				targetRow := neighboors[targetIndex][0]
				targetCol := neighboors[targetIndex][1]

				if m.table[targetRow][targetCol] == empty {
					// Target cell is empty. Move enemy and clear the old one.
					m.table[targetRow][targetCol] = enemy
					m.table[row][col] = empty
				} else if m.table[targetRow][targetCol] == player {
					// Target cell contains the player. Attack and stop further processing, this game is over!
					m.gameOver = true

					return
				}
			}
		}
	}
}

func (m *model) randomFreeCoordinates() (row, col int) {
	// Generate some random coordinates.
	row, col = randomCoordinates()

	// Check that the random cell is empty.
	// If not repeat randomizing until we find an empty cell.
	for m.table[row][col] != empty {
		row, col = randomCoordinates()
	}

	return
}

// randomCoordinates returns a new set of random coordinates within the
// playing field excluding borders. However, it is not guaranteed that the
// cell under the returned coordinates is actually empty.
func randomCoordinates() (row, col int) {
	row = rand.Intn(tableHeight-2) + 1
	col = rand.Intn(tableWidth-2) + 1

	return
}

func (m *model) movePlayer(row, col int) {
	if m.gameOver {
		return
	}

	// Clear old player location.
	m.table[m.playerRow][m.playerCol] = empty

	if m.table[row][col] == enemy {
		// We ran into an enemy. Signal game over and skip further
		// processing.
		m.gameOver = true
		return
	}

	if m.table[row][col] == item {
		// We collected an item. A new item and enemy needs to be
		// spawned. Increase the score.
		m.score++
		m.spawnItem()
		m.spawnEnemy()
	}

	// Set new player location.
	m.table[row][col] = player
	m.playerRow = row
	m.playerCol = col
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

// init is responsible for initializing or resetting a model that is ready
// to use for a new game.
func (m *model) init() {
	// Clear and reset all fields as init() is also used for restarting
	// an existing game. Therefore, our model needs to be fresh.

	// Initially, set every cell to our empty symbol.
	for row := 0; row < tableHeight; row++ {
		for col := 0; col < tableWidth; col++ {
			m.table[row][col] = empty
		}
	}

	m.playerRow = 0
	m.playerCol = 0
	m.score = 0
	m.gameOver = false

	// Set the four corners.
	m.table[0][0] = corner
	m.table[0][tableWidth-1] = corner
	m.table[tableHeight-1][0] = corner
	m.table[tableHeight-1][tableWidth-1] = corner

	// Draw horizontal borders at the top and bottom.
	for col := 1; col < tableWidth-1; col++ {
		m.table[0][col] = lineHorizontal
		m.table[tableHeight-1][col] = lineHorizontal
	}

	// Draw vertical borders on the left and right side.
	for row := 1; row < tableHeight-1; row++ {
		m.table[row][0] = lineVertical
		m.table[row][tableWidth-1] = lineVertical
	}

	// Spawn our player near the top left corner.
	m.playerRow = 1
	m.playerCol = 1
	m.table[m.playerRow][m.playerCol] = player

	m.spawnItem()
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

		// If our game is lost, any key shall restart the game.
		if m.gameOver {

			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				m.init()
				return m, nil
			}
		}

		switch msg.String() {

		// Exit program on ctrl+c or q typing.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.playerUp()
			m.moveEnemies()
		case "down":
			m.playerDown()
			m.moveEnemies()
		case "left":
			m.playerLeft()
			m.moveEnemies()
		case "right":
			m.playerRight()
			m.moveEnemies()
		}
	}

	return m, nil
}

// View is required for building what we want to show on the screen.
// That means we need to translate our model data into a string for displaying.
func (m *model) View() string {
	builder := strings.Builder{}

	if m.gameOver {
		// Just inform about game over and don't continue.
		builder.WriteString("\n\n\n\n\n")
		builder.WriteString("          You died, Game Over!")
		builder.WriteString("\n\n")
		builder.WriteString(fmt.Sprintf("          Your score: %d", m.score))
		builder.WriteString("\n\n")
		builder.WriteString("          Press enter to restart or q to quit")

		return builder.String()
	}

	// Iterate our table (2d array) and print the cells.
	for _, row := range m.table {

		for _, cell := range row {
			builder.WriteRune(cell)
		}

		// Go to next line after each row.
		builder.WriteString("\n")
	}

	return builder.String()
}

func main() {
	// Create our initial model.
	model := &model{}
	model.init()

	// Program setup to initialize bubbletea and use full screen.
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Run bubbletea and exist with a message if an error occurs.
	if _, err := program.Run(); err != nil {
		fmt.Println("Unexpected error: %v", err)
		os.Exit(1)
	}
}
