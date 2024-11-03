package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO (rough priority order)
// soon
// - [ ]"select mode":
//     - move cursor around to highlight a chunk of blocks
//     - exit select mode to keep only the time stamp of the first block
//     - shift all affected blocks later
// - [ ]configuration
//     - [x]start/end of work day
//     - [ ]colors
// - [ ]start viewport at configured time
// later
// - [ ]on startup, start viewport and cursor at current hour by default
// - [ ]improve styles
// - [ ]specific todo list per task
// - [ ]autocomplete menu for activities
// - [ ]undo/redo (+navigable history?)
// - [ ]copy and paste blocks
// - [ ]import/export JSON (or something else)
// - [ ]option to increase/decrease time resolution (globally and per block)
// - [ ]sqlite persistent storage (store json blobs?)
// - [ ]navigate through previous days and later days
// - [ ]pagination/scrolling for small terminal windows
// - [ ]alternate tabular view
// - [ ]move cursor outside (to the left of) schedule block
// - [ ]check boxes
// - [ ]day analysis (how many tasks checked off?)
// - [ ]color scheme/schemes
// - [ ]update time live (as a bubbletea Cmd)
// - [ ]function to convert timeMsg to timestamp
// completed
// - [x]check current time with bubbleTea command
// - [x]always show 24 hours in the day
// - [x]implement scrolling
// - [x]grey out past time blocks
// - [x]label time blocks with start and end times (based on blocksPerHour)
// - [x]refactor into separate files
// - [x]highlight current time block
// - [x]fix formatting of current time block

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
