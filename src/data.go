package main

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	hoursInDay = 24

	blocksPerHour = 2

	dayStartTime = 0
)

type model struct {
	// TODO: update current time while app is open
	currentTime time.Time

	// timeSpan int // the total number of minutes in the workday

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

	mode int
	err  error

	width  int
	height int

	blockHeight int
	blockWidth  int

	blockLabels []string

	styles *styles

	logger *slog.Logger

	// ready bool

	// vpStart holds the index of the first time block visible in the terminal window
	vpStart int

	// vpRange holds the number of time blocks that will be visible in a single window.
	// It's derived from the window size.
	vpRange int
}

func (m model) Init() tea.Cmd {
	// `nil` means "no I/O right now, please."
	return readConfig
}

type configMsg map[string]string

func readConfig() tea.Msg {
	contents, err := os.ReadFile("../config.ini")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(contents), "\n")
	config := make(configMsg)
	for i := range lines {
		if len(lines[i]) < 3 {
			continue
		}
		key, val, found := strings.Cut(lines[i], "=")
		if !found {
			panic("Invalid line in config")
		}

		config[strings.TrimSpace(key)] = strings.TrimSpace(val)
	}
	return config
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

		vpStart: 0,
	}
}
