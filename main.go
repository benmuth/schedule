package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO (rough priority order)
// soon
// -[ ] label time blocks with start and end times (based on blocksPerHour)
// -[ ] "visual mode": expand/contract a task to cover more/fewer blocks
// -[ ] highlight current time block
// -[ ] grey out past time blocks
// later
// -[ ] specific todo list per task
// -[ ] autocomplete menu for activities
// -[ ] undo/redo (+navigable history?)
// -[ ] copy and paste blocks
// -[ ] import/export JSON (or something else)
// -[ ] increase/decrease time resolution (globally and per block)
// -[ ] configuration
// -[ ] sqlite persistent storage (store json blobs?)
// -[ ] navigate through previous days and later days
// -[ ] pagination/scrolling for small terminal windows
// -[ ] alternate tabular view

const (
	hoursInDay = 8
	// BUG: selections are weird if blocksPerHour is 2 and hoursInDay is 8
	blocksPerHour = 1

	taskBackgroundColor = "21"
	blockWidth          = 40
)

// styles
var (
	normalTask   = lipgloss.NewStyle().Width(blockWidth).Margin(1, 0, 0).Border(lipgloss.NormalBorder(), true)
	selectedTask = normalTask.Copy().Background(lipgloss.Color(taskBackgroundColor))
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	time     int              // the total number of minutes in the workday
	blocks   int              // how many blocks of time in the workday
	tasks    []string         // what to do in each block of time
	cursor   int              // which time block our cursor is pointing at
	selected map[int]struct{} // which time blocks are selected

	textInput textinput.Model
	isEditing bool
	err       error

	width  int
	height int
}

type errMsg error

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter a task"
	ti.Blur()
	ti.CharLimit = 156
	ti.Width = blockWidth - 10
	ti.TextStyle.Background(lipgloss.Color(taskBackgroundColor))
	ti.PlaceholderStyle.Background(lipgloss.Color(taskBackgroundColor))
	ti.Cursor.SetMode(cursor.CursorStatic)

	blocks := hoursInDay * blocksPerHour
	activities := make([]string, blocks)

	return model{
		time: 60 * hoursInDay,

		blocks: blocks,
		tasks:  activities,

		// The keys refer to the indices of the `activities` slice.
		selected: make(map[int]struct{}),

		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	// return textinput.Blink
	// `nil` means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if !m.isEditing {
				return m, tea.Quit
			}

		case "up", "k":
			m.cursor = m.moveCursorUp(1)

		case "down", "j":
			m.cursor = m.moveCursorDown(1)

		// selects a block to move around
		case "enter":
			if !m.isEditing {
				m.toggleSelectedBlock()
			}

		// insert mode (edit activity text)
		case "i":
			if !m.textInput.Focused() {
				m.textInput.SetValue(m.tasks[m.cursor])
				m.textInput.Focus()
				m.isEditing = true
			}
			if m.hasSelectedBlock() {
				m.toggleSelectedBlock()
			}

		// visual mode (select a block and expand the length of it)
		// case "v":

		// return to "normal mode"
		case "esc":
			if m.textInput.Focused() || m.isEditing {
				m.textInput.Blur()
				m.tasks[m.cursor] = m.textInput.Value()
				m.isEditing = false
				m.textInput.Reset()
			}
			if m.hasSelectedBlock() {
				m.toggleSelectedBlock()
			}
		}

	case errMsg:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.assertInvariants()

	// Return the updated model to the Bubble Tea runtime
	return m, cmd
}

func (m model) View() string {
	s := "Schedule:\n\n"

	for i, task := range m.tasks {

		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		var block string
		// fill the block
		if m.isEditing && m.cursor == i {
			block = fmt.Sprintf("%s\n", m.textInput.View())
		} else {
			block = fmt.Sprintf("%s %s\n", cursor, task)
		}

		if _, ok := m.selected[i]; ok {
			s += selectedTask.Render(block)
		} else {
			s += normalTask.Render(block)
		}

	}

	s += "\nPress q to quit.\n"
	return s
}

func (m model) moveCursorUp(amount int) int {
	initial := m.cursor
	if m.cursor > 0 && !m.isEditing {
		final := initial - amount
		m.moveSelectedBlock(initial, final)
		return final
	}
	return initial
}

func (m model) moveCursorDown(amount int) int {
	initial := m.cursor
	if m.cursor < len(m.tasks)-1 && !m.isEditing {
		final := initial + amount
		m.moveSelectedBlock(initial, final)
		return final
	}
	return initial
}

func (m model) moveSelectedBlock(initial, final int) {
	if _, ok := m.selected[initial]; ok {
		swapBlocks(m.tasks, initial, final)
		delete(m.selected, initial)
		m.selected[final] = struct{}{}
	}
}

func swapBlocks(tasks []string, a, b int) {
	tasks[a], tasks[b] = tasks[b], tasks[a]
}

func (m model) hasSelectedBlock() bool {
	_, ok := m.selected[m.cursor]
	return ok
}

func (m model) toggleSelectedBlock() {
	if m.hasSelectedBlock() {
		delete(m.selected, m.cursor)
	} else {
		m.selected[m.cursor] = struct{}{}
	}
}

func (m model) assertInvariants() {
	if len(m.selected) > 1 {
		panic(fmt.Sprintf("too many elements selected! want 1 have %v", len(m.selected)))
	}

	if m.hasSelectedBlock() && m.isEditing {
		panic(fmt.Sprintf("selected block while editing! selected block at index %v. cursor at %v", m.selected, m.cursor))
	}
}
