package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type tickMsg time.Time

const (
	normalMode = iota
	insertMode
	selectMode
	stretchMode
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.mode != insertMode {
				return m, tea.Quit
			}

		case "up", "k":
			if m.mode == stretchMode {
				m.cursor = m.stretch(-1)
			} else if m.mode != insertMode {
				m.cursor = m.moveCursor(-1)
			}

		case "down", "j":
			if m.mode == stretchMode {
				m.cursor = m.stretch(1)
			} else if m.mode != insertMode {
				m.cursor = m.moveCursor(1)
			}

		// selects a block to move around
		case "enter":
			if m.mode != insertMode {
				m.mode = selectMode
				m.toggleSelectedBlock()
			}

		// insert mode (edit activity text)
		case "i":
			// NOTE: we should have a single source of truth for what defines
			// insert mode. Is it focused text input or mode == insertMode? We
			// can then derive further logic from there.
			if !m.textInput.Focused() || !(m.mode == insertMode) {
				m.textInput.SetValue(m.tasks[m.cursor])
				m.textInput.Focus()
				// m.insertMode = true
				m.mode = insertMode
			}
			delete(m.selected, m.cursor)

		// stretch mode (select a block and change its length)
		case "v":
			if !(m.mode == insertMode || m.mode == stretchMode) {
				m.mode = stretchMode
			} else {
				m.mode = normalMode
			}
			// implementation notes:
			// - have some sort of data in model that represents a contiguous
			// span of time, probably just a start and end index.
			// - stretch mode should have its own color
			// - adjacent blocks that are part of the same stretch should have
			// no gap in between them
			// - only the first block of a stretch should display the time. The
			// end time will be shown by the block after the stretch
			// - "selecting" a stretched block moves it just as a normal block
			// - when changing the span of a stretch, going above the "anchor"
			// block just extends the stretch upwards

		// return to "normal mode"
		case "esc":
			if m.textInput.Focused() || m.mode == insertMode {
				m.textInput.Blur()
				m.tasks[m.cursor] = m.textInput.Value()
				// m.insertMode = false
				m.textInput.Reset()
			}
			if m.blockIsSelected() {
				m.toggleSelectedBlock()
			}
			m.mode = normalMode
		}

	case errMsg:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.logger.Info("New window size", "width", fmt.Sprintf("%d", msg.Width), "height", fmt.Sprintf("%d", msg.Height))
		m.width = msg.Width
		m.height = msg.Height
		// m.resize()
		// normalTask.Width(m.width - (m.width / 10))
	}

	m.assertInvariants()

	// Return the updated model to the Bubble Tea runtime
	return m, cmd
}

// func (m model) resize() {
// 	// height := int(math.Floor(float64(*m.height)/float64(m.numBlocks)) * float64(0.6))
// 	// if height < 2 {
// 	// height = 2
// 	// }
// 	*m.height = 2
// 	width := int(math.Floor(float64(*m.width) - float64(float64(*m.width)/float64(10))))
// 	if width < 20 {
// 		width = 20
// 	}
// 	m.styles.normalBlock = m.styles.normalBlock.Width(width).Height(*m.height)
// 	m.styles.currentBlock = m.styles.currentBlock.Width(width).Height(*m.height)
// 	m.styles.selectedBlock = m.styles.selectedBlock.Width(width).Height(*m.height)
// 	m.styles.pastBlock = m.styles.pastBlock.Width(width).Height(*m.height)
// }

// stretch extends the current span in the given direction
// TODO: make this non destructive somehow. maybe push adjacent spans into
// an overflow buffer on either side of the visible span buffer
func (m model) stretch(dir int) (finalCursor int) {
	initialNumber := m.spans[m.cursor]
	start, end := m.getSpanEnds(m.cursor)

	if dir < 0 {
		if start > 0 {
			start += dir
		}
		finalCursor = start
	}
	if dir > 0 {
		if end < len(m.spans)-1 {
			end += dir
		}
		finalCursor = end
	}

	if start >= end {
		panic(fmt.Sprintf("startAnchor greater than or equal to endAnchor: start %v, end %v\n", start, end))
	}

	for i := range m.spans {
		if i >= start && i <= end {
			m.spans[i] = initialNumber
		}
	}

	return
}

func (m model) getSpanEnds(idx int) (int, int) {
	initialNumber := m.spans[idx]
	startAnchor := -1
	endAnchor := -1

	for i, number := range m.spans {
		if number == initialNumber {
			if startAnchor < 0 {
				startAnchor = i
			}
			endAnchor = i
		}
	}
	return startAnchor, endAnchor
}

func (m model) moveCursor(amount int) int {
	initial := m.cursor
	if m.cursor >= 0 && m.cursor <= len(m.tasks)-1 {
		final := initial + amount
		if final < 0 {
			final = 0
		}
		if final > len(m.tasks)-1 {
			final = len(m.tasks) - 1
		}
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

func (m model) blockIsSelected() bool {
	_, ok := m.selected[m.cursor]
	return ok
}

func (m model) toggleSelectedBlock() {
	if m.blockIsSelected() {
		delete(m.selected, m.cursor)
	} else {
		m.selected[m.cursor] = struct{}{}
	}
}

// HACK: I don't know if this is the best way to do this, but it works for now
// Enforce simple invariants about the model on each update loop
func (m model) assertInvariants() {
	if len(m.selected) > 1 {
		panic(fmt.Sprintf("too many elements selected! want 1 have %v", len(m.selected)))
	}

	if m.blockIsSelected() && m.mode == insertMode {
		panic(fmt.Sprintf("selected block while editing! selected block at index %v. cursor at %v", m.selected, m.cursor))
	}
}
