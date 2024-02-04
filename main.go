package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh/"

type model struct {
	status int
	err    error
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

// checkServer is a Cmd. Cmds are functions that perform IO then return a Msg. All
// IO should be done through commands

// Cmds are functions that take no arguments and return a Msg
// if you need to pass something to a command, use a wrapper function that returns a command, like

// func checkSomeUrl(url string) tea.Cmd {
//     return func() tea.Msg {
//         c := &http.Client{Timeout: 10 * time.Second}
//         res, err := c.Get(url)
//         if err != nil {
//             return errMsg{err}
//         }
//         return statusMsg(res.StatusCode)
//     }
// }

// only the returned function will act asynchronously, so don't do slow stuff and block the main
// goroutine in the outer function

// make a request to a server then return a Msg
func checkServer() tea.Msg {
	// Create an HTTP client and make a GET request.
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		// There was an error making our request. Wrap the error we received
		// in a message and return it.
		return errMsg{err}
	}
	// We received a response from the server. Return the HTTP status code
	// as a message.
	return statusMsg(res.StatusCode)
}

type statusMsg int

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() tea.Cmd {
	return checkServer
}

// Cmds run asynchronously in a goroutine, and their returned Msgs are sent to
// the Update function, where they're handled based on their type
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		// The server returned a status message. Save it to our model. Also
		// tell the Bubble Tea runtime we want to exit because we have nothing
		// else to do. We'll still be able to render a final view with our
		// status message.
		m.status = int(msg)
		return m, tea.Quit

	case errMsg:
		// There was an error. Note it in the model. And tell the runtime
		// we're done and want to quit.
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		// Ctrl+c exits. Even with short running programs it's good to have
		// a quit key, just in case your logic is off. Users will be very
		// annoyed if they can't exit.
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	// If we happen to get any other messages, don't do anything.
	return m, nil
}

func (m model) View() string {
	// If there's an error, print it out and don't do anything else.
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// Tell the user we're doing something.
	s := fmt.Sprintf("Checking %s ... ", url)

	// When the server responds with a status, add it to the current line.
	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	// Send off whatever we came up with above for rendering.
	return "\n" + s + "\n\n"
}

// package main

// import (
// 	"fmt"
// 	"os"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// func main() {
// 	p := tea.NewProgram(initialModel())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Alas, there's been an error: %v", err)
// 		os.Exit(1)
// 	}
// }

// type model struct {
// 	choices  []string         // items on the to-do list
// 	cursor   int              // which to-do list item our cursor is pointing at
// 	selected map[int]struct{} // which to-do items are selected
// }

// func initialModel() model {
// 	return model{
// 		// Our to-do list is a grocery list
// 		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

// 		// A map which indicates which choices are selected. We're using
// 		// the  map like a mathematical set. The keys refer to the indexes
// 		// of the `choices` slice, above.
// 		selected: make(map[int]struct{}),
// 	}
// }

// // models have 3 methods
// // init
// // update
// // view

// func (m model) Init() tea.Cmd {
// 	// Just return `nil`, which means "no I/O right now, please."
// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	// Is it a key press?
// 	case tea.KeyMsg:

// 		// Cool, what was the actual key pressed?
// 		switch msg.String() {

// 		// These keys should exit the program.
// 		case "ctrl+c", "q":
// 			return m, tea.Quit

// 		// The "up" and "k" keys move the cursor up
// 		case "up", "k":
// 			if m.cursor > 0 {
// 				m.cursor--
// 			}

// 		// The "down" and "j" keys move the cursor down
// 		case "down", "j":
// 			if m.cursor < len(m.choices)-1 {
// 				m.cursor++
// 			}

// 		// The "enter" key and the spacebar (a literal space) toggle
// 		// the selected state for the item that the cursor is pointing at.
// 		case "enter", " ":
// 			_, ok := m.selected[m.cursor]
// 			if ok {
// 				delete(m.selected, m.cursor)
// 			} else {
// 				m.selected[m.cursor] = struct{}{}
// 			}
// 		}
// 	}

// 	// Return the updated model to the Bubble Tea runtime for processing.
// 	// Note that we're not returning a command.
// 	return m, nil
// }

// func (m model) View() string {
// 	// The header
// 	s := "What should we buy at the market?\n\n"

// 	// Iterate over our choices
// 	for i, choice := range m.choices {

// 		// Is the cursor pointing at this choice?
// 		cursor := " " // no cursor
// 		if m.cursor == i {
// 			cursor = ">" // cursor!
// 		}

// 		// Is this choice selected?
// 		checked := " " // not selected
// 		if _, ok := m.selected[i]; ok {
// 			checked = "x" // selected!
// 		}

// 		// Render the row
// 		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
// 	}

// 	// The footer
// 	s += "\nPress q to quit.\n"

// 	// Send the UI for rendering
// 	return s
// }
