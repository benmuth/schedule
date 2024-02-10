package main

import (
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type tickMsg time.Time

const (
	normalMode = iota
	insertMode
	selectMode
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

func (m model) resize() {
	height := int(math.Floor(float64(m.height)/float64(m.numBlocks)) * float64(0.6))
	if height < 2 {
		height = 2
	}
	width := int(math.Floor(float64(m.width) - float64(float64(m.width)/float64(10))))
	if width < 20 {
		width = 20
	}
	normalBlock = normalBlock.Width(width).Height(height)
	currentBlock = currentBlock.Width(width).Height(height)
	selectedBlock = selectedBlock.Width(width).Height(height)
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
