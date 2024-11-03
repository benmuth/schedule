package main

import (
	"fmt"
	"math"

	// "time"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

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
				panic("Not implemented")
			} else if m.mode != insertMode {
				initial := m.cursor
				m.cursor = m.moveCursor(-1)
				if _, ok := m.selected[initial]; ok {
					m.moveSelectedBlock(initial, m.cursor)
				}
				m.vpStart = m.adjustVPStart()
			}

		case "down", "j":
			if m.mode == stretchMode {
				panic("Not implemented")
			} else if m.mode != insertMode {
				initial := m.cursor
				m.cursor = m.moveCursor(1)
				if _, ok := m.selected[initial]; ok {
					m.moveSelectedBlock(initial, m.cursor)
				}
				m.vpStart = m.adjustVPStart()
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

		m.vpRange = m.height / 3

	// TODO: figure out if initialization like this is necessary (look at bubbletea examples)
	// m.ready = true
	// } else {
	// 	m.width = msg.Width
	// 	m.height = msg.Height
	// }

	// m.resize()
	// normalTask.Width(m.width - (m.width / 10))
	case configMsg:
		m.styles = changeStyles(msg)
	case timeMsg:
		m.currentTime = msg
	}

	blockHeight := int(math.Floor(float64(m.height)/float64(m.numBlocks)) * float64(0.1))
	if blockHeight < 1 {
		blockHeight = 1
	}
	m.blockHeight = blockHeight

	blockWidth := int(math.Floor(float64(m.width) - float64(float64(m.width)/10.0)))
	if blockWidth < 20 {
		blockWidth = 20
	}
	m.blockWidth = blockWidth

	m.assertInvariants()

	return m, cmd
}

func (m model) moveCursor(amount int) int {
	newPos := m.cursor + amount

	if newPos < 0 {
		return 0
	}

	if newPos > len(m.tasks)-1 {
		return len(m.tasks) - 1
	}

	return newPos
}

func (m model) adjustVPStart() int {
	if m.cursor < m.vpStart {
		return m.cursor
	}

	if m.cursor >= m.calcVPEnd()-1 {
		m.logger.Error(fmt.Sprintf("cursor: %v, vpStart: %v, vpEnd: %v, vpRange: %v", m.cursor, m.vpStart, m.calcVPEnd(), m.vpRange))
		return max(m.cursor-m.vpRange+1, 0)
	}
	// no change
	return m.vpStart
}

func (m model) calcVPEnd() int {
	return min(m.vpStart+m.vpRange, m.numBlocks)
}

func (m model) moveSelectedBlock(initial, final int) {
	swapBlocks(m.tasks, initial, final)
	delete(m.selected, initial)
	m.selected[final] = struct{}{}
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

// HACK: Probably not the best way to do this, but it works for now
// Enforce simple invariants about the model on each update loop
func (m model) assertInvariants() {
	if len(m.selected) > 1 {
		panic(fmt.Sprintf("too many elements selected! want 1 have %v", len(m.selected)))
	}

	if m.blockIsSelected() && m.mode == insertMode {
		panic(fmt.Sprintf("selected block while editing! selected block at index %v. cursor at %v", m.selected, m.cursor))
	}

	if m.vpRange <= 0 {
		panic(fmt.Sprintf("vpRange too small: %v", m.vpRange))
	}

	if m.calcVPEnd() < 0 || m.calcVPEnd() > m.numBlocks {
		panic(fmt.Sprintf("invalid vp end: %v", m.calcVPEnd()))
	}
}
