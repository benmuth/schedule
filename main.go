package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const hoursInDay = 8

const blocksPerHour = 1

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	time       int              // the total number of minutes in the workday
	blocks     int              // how many blocks of time in the workday
	activities []string         // what to do in each block of time
	cursor     int              // which time block our cursor is pointing at
	selected   map[int]struct{} // which time blocks are selected

	textInput textinput.Model
	editing   bool
	err       error
}

type errMsg error

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter an activity"
	ti.Blur()
	ti.CharLimit = 156
	ti.Width = 20

	blocks := hoursInDay * blocksPerHour
	activities := make([]string, blocks)

	return model{
		time: 60 * hoursInDay,

		blocks: blocks,
		// Our to-do list is a grocery list
		activities: activities,

		// A map which indicates which choices are selected. The keys refer to the indices
		// of the `activities` slice.
		selected: make(map[int]struct{}),

		textInput: ti,
		err:       nil,
	}
}

// models have 3 methods
// init
// update
// view

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 && !m.editing {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.activities)-1 && !m.editing {
				m.cursor++
			}

		// The "enter" key and the spacebar toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "i":
			if !m.textInput.Focused() {
				m.textInput.SetValue(m.activities[m.cursor])
				m.textInput.Focus()
				m.editing = true
			}
		case "esc":
			if m.textInput.Focused() {
				m.textInput.Blur()
				m.activities[m.cursor] = m.textInput.Value()
				m.editing = false
				m.textInput.Reset()
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	// Return the updated model to the Bubble Tea runtime
	return m, cmd
}

func (m model) View() string {
	// The header
	s := "Schedule:\n\n----------------\n"

	// s += fmt.Sprintf("editing: %v\n", m.editing)
	// Iterate over our choices
	for i, activity := range m.activities {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
			// if m.editing {
			// s += m.textInput.View()
			// }
		}

		// Is this choice selected?
		// checked := " " // not selected
		// if _, ok := m.selected[i]; ok {
		// 	checked = "x" // selected!
		// }

		// Render the row
		if m.editing && m.cursor == i {
			s += fmt.Sprintf("%s\n", m.textInput.View())
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, activity)
		}
		s += "----------------\n"
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
