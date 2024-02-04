package main

import (
	"fmt"
	"math"
	"os"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO (rough priority order)
// soon
// -[ ] "select mode": expand/contract a task to cover more/fewer blocks
// -[ ] highlight current time block
// -[ ] grey out past time blocks
// -[ ] refactor into separate files
// later
// -[ ] improve styles
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
// completed
// -[x] label time blocks with start and end times (based on blocksPerHour)

const (
	hoursInDay = 8
	// BUG: selections are weird if blocksPerHour is 2 and hoursInDay is 8
	blocksPerHour = 1

	dayStartTime = 9

	dayEndTime = dayStartTime + hoursInDay

	taskBackgroundColor = "21"
)

const (
	normalMode = iota
	insertMode
	selectMode
)

// styles
var (
	normalTask   = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true)
	selectedTask = normalTask.Copy().Background(lipgloss.Color(taskBackgroundColor))
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	time int // the total number of minutes in the workday
	// NOTE: maybe this can be a constant
	numBlocks int              // how many blocks of time in the workday
	tasks     []string         // what to do in each block of time
	cursor    int              // which time block our cursor is pointing at
	selected  map[int]struct{} // which time blocks are selected

	textInput textinput.Model
	// insertMode bool
	mode int
	err  error

	width  int
	height int

	blockLabels []string
}

type errMsg error

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter a task"
	ti.Blur()
	ti.CharLimit = 156
	// ti.Width = 40 - 10
	ti.TextStyle.Background(lipgloss.Color(taskBackgroundColor))
	ti.PlaceholderStyle.Background(lipgloss.Color(taskBackgroundColor))
	ti.Cursor.SetMode(cursor.CursorStatic)

	numBlocks := hoursInDay * blocksPerHour
	labels := makeBlockLabels(numBlocks)
	activities := make([]string, numBlocks)

	return model{
		time: 60 * hoursInDay,

		numBlocks: numBlocks,
		tasks:     activities,

		// The keys refer to the indices of the `activities` slice.
		selected: make(map[int]struct{}),

		textInput: ti,
		err:       nil,

		blockLabels: labels,
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
			if !(m.mode == insertMode) {
				return m, tea.Quit
			}

		case "up", "k":
			if !(m.mode == insertMode) {
				m.cursor = m.moveCursorUp(1)
			}

		case "down", "j":
			if !(m.mode == insertMode) {
				m.cursor = m.moveCursorDown(1)
			}

		// selects a block to move around
		case "enter":
			if !(m.mode == insertMode) {
				m.toggleSelectedBlock()
			}

		// insert mode (edit activity text)
		case "i":
			if !m.textInput.Focused() || !(m.mode == insertMode) {
				m.textInput.SetValue(m.tasks[m.cursor])
				m.textInput.Focus()
				// m.insertMode = true
				m.mode = insertMode
			}
			if m.hasSelectedBlock() {
				m.toggleSelectedBlock()
			}

		// select  mode (select a block and expand the length of it)
		// case "v":

		// return to "normal mode"
		case "esc":
			if m.textInput.Focused() || m.mode == insertMode {
				m.textInput.Blur()
				m.tasks[m.cursor] = m.textInput.Value()
				// m.insertMode = false
				m.mode = normalMode
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
		m.resize()
		// normalTask.Width(m.width - (m.width / 10))
	}

	m.assertInvariants()

	// Return the updated model to the Bubble Tea runtime
	return m, cmd
}

func (m model) View() string {
	s := "Schedule"

	for i, task := range m.tasks {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		var block string
		// fill the block
		block += fmt.Sprintf("%s\n", m.blockLabels[i])
		if m.mode == insertMode && m.cursor == i {
			block += fmt.Sprintf("%s\n", m.textInput.View())
		} else {
			block += fmt.Sprintf("%s %s\n", cursor, task)
		}

		if _, ok := m.selected[i]; ok {
			s += selectedTask.Render(block)
		} else {
			s += normalTask.Render(block)
		}
	}

	s += "\nPress q to quit.\n"
	s += m.debugInfo()
	return s
}

func (m model) moveCursorUp(amount int) int {
	initial := m.cursor
	if m.cursor > 0 {
		final := initial - amount
		m.moveSelectedBlock(initial, final)
		return final
	}
	return initial
}

func (m model) moveCursorDown(amount int) int {
	initial := m.cursor
	if m.cursor < len(m.tasks)-1 {
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

	if m.hasSelectedBlock() && m.mode == insertMode {
		panic(fmt.Sprintf("selected block while editing! selected block at index %v. cursor at %v", m.selected, m.cursor))
	}
}

func (m model) resize() {
	height := int(math.Floor(float64(m.height)/float64(m.numBlocks)) * float64(0.6))
	if height < 2 {
		height = 2
	}
	width := int(math.Floor(float64(m.width) - float64(float64(m.width)/float64(10))))
	if width < 20 {
		width = 20
	}
	normalTask.Width(width).Padding(0, 1, 0).Height(height).Margin(1, 1, 0)
	selectedTask = normalTask.Copy().Background(lipgloss.Color(taskBackgroundColor))
}

func (m model) showMode() string {
	switch m.mode {
	case normalMode:
		return "NOR"
	case insertMode:
		return "INS"
	case selectMode:
		return "SEL"
	}
	return ""
}

func (m model) debugInfo() string {
	return fmt.Sprintf("\n%s | height: %v | width: %v \n", m.showMode(), m.height, m.width)
}

// type blockLabels struct {
// 	startTime string
// 	endTime   string
// }

func makeBlockLabels(numBlocks int) []string {
	labels := make([]string, numBlocks)

	time := float64(dayStartTime)
	interval := float64(1) / float64(blocksPerHour)
	for i := 0; i < len(labels); i++ {
		labels[i] = conv24To12(time)
		time += float64(interval)
	}
	return labels
}

// conv24To12 converts a 24 hour timestamp into a 12 hour clock time string.
// time24 represents the hour of the day and must be between 0.0 and 24.0.
func conv24To12(time24 float64) string {
	integer, fraction := math.Modf(time24)
	hrs := int(integer) % 12
	mins := math.Floor(fraction * 60)

	var period string
	if time24 < 12 {
		period = "am"
	} else {
		period = "pm"
	}
	return fmt.Sprintf("%v:%02v %s", hrs, mins, period)
}
