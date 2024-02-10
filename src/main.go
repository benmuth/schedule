package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO (rough priority order)
// soon
// -[ ] "select mode": expand/contract a task to cover more/fewer blocks
// -[ ] grey out past time blocks
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
// -[x] refactor into separate files
// -[x] highlight current time block
// -[x] fix formatting of current time block

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
