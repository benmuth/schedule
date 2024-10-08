package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	hoursInDay = 5
	// BUG: selections and cursors are weird if there are too many blocks for the screen size
	blocksPerHour = 1

	dayStartTime = 11

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
	// selected is a set containing the indices of the selected activities
	selected map[int]struct{}

	// spans holds the ends of time block spans
	spans []int

	textInput textinput.Model
	// insertMode bool
	mode int
	err  error

	width  int
	height int

	blockLabels []string

	styles *styles

	logger *slog.Logger

	viewport viewport.Model

	ready bool
}

func (m model) Init() tea.Cmd {
	// return textinput.Blink
	// `nil` means "no I/O right now, please."
	return nil
}

func initialModel() model {
	styles := defaultStyles()

	ti := textinput.New()
	ti.Placeholder = "Enter a task"
	ti.Blur()
	ti.CharLimit = 156
	ti.TextStyle.Inherit(styles.tiTextStyle)
	ti.PlaceholderStyle.Inherit(ti.PlaceholderStyle)
	ti.Cursor.SetMode(cursor.CursorStatic)

	numBlocks := hoursInDay * blocksPerHour
	labels := makeBlockLabels(numBlocks, dayStartTime, blocksPerHour)
	activities := make([]string, numBlocks)

	spans := make([]int, numBlocks)
	for i := range spans {
		spans[i] = i
	}

	f, err := os.Create("../rescheduler.log")
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(f, nil))

	width := 20
	height := 2

	return model{
		currentTime: time.Now(),

		timeSpan: 60 * hoursInDay,

		numBlocks: numBlocks,
		tasks:     activities,

		selected: make(map[int]struct{}),

		spans: spans,

		textInput: ti,
		err:       nil,

		blockLabels: labels,

		styles: &styles,

		logger: logger,

		width:  width,
		height: height,
	}
}
