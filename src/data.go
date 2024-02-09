package main

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	hoursInDay = 8
	// BUG: selections are weird if blocksPerHour is 2 and hoursInDay is 8
	blocksPerHour = 1

	dayStartTime = 9

	dayEndTime = dayStartTime + hoursInDay

	taskBackgroundColor = "21"
)

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

func (m model) Init() tea.Cmd {
	// return textinput.Blink
	// `nil` means "no I/O right now, please."
	return nil
}

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
