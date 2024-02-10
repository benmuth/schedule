package main

import (
	"time"

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
)

type model struct {
	currentTime time.Time

	timeSpan int // the total number of minutes in the workday

	// NOTE: maybe this can be a constant
	// numBlocks is the number of blocks of time in the workday
	numBlocks int

	// tasks hold what is scheduled for each block of time. A task
	// may be scheduled for more than one block.
	tasks []string

	// cursor indicates which time block our cursor is pointed at
	cursor int

	// selected holds the selected time blocks
	selected map[int]struct{}

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
	ti.TextStyle.Background(lipgloss.Color(selectedBlockBackgroundColor))
	ti.PlaceholderStyle.Background(lipgloss.Color(selectedBlockBackgroundColor))
	ti.Cursor.SetMode(cursor.CursorStatic)

	numBlocks := hoursInDay * blocksPerHour
	labels := makeBlockLabels(numBlocks, dayStartTime, blocksPerHour)
	activities := make([]string, numBlocks)

	return model{
		currentTime: time.Now(),

		timeSpan: 60 * hoursInDay,

		numBlocks: numBlocks,
		tasks:     activities,

		// selected is a set containing the indices of the selected activities
		selected: make(map[int]struct{}),

		textInput: ti,
		err:       nil,

		blockLabels: labels,
	}
}
