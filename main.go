package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// model encapsulates our data for displaying and updating.
type model struct {
	text string
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
		}
	}

	return m, nil
}

// View is required for building what we want to show on the screen.
// For now, lets just print a text.
func (m *model) View() string {
	// Simply print our models text to the screen.
	return m.text
}

func main() {
	// Create our initial model with a sample text.
	model := &model{
		text: "Hello world!",
	}

	// Program setup to initialize bubbletee and use full screen.
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Run bubbletee and exist with a message if an error occurs.
	if _, err := program.Run(); err != nil {
		fmt.Println("Unexpected error: %v", err)
		os.Exit(1)
	}
}
