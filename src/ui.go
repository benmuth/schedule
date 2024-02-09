package main

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

// styles
var (
	normalTask   = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true)
	selectedTask = normalTask.Copy().Background(lipgloss.Color(taskBackgroundColor))
)

func (m model) View() string {
	s := "Schedule"

	for i, task := range m.tasks {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		var block string
		// fill the block
		block += fmt.Sprintf("%s\n", m.blockLabels[i])
		if m.mode == insertMode && m.cursor == i {
			block += fmt.Sprintf("%s\n", m.textInput.View())
		} else {
			block += fmt.Sprintf("%s %s\n", cursor, task)
		}

		if _, ok := m.selected[i]; ok {
			s += selectedTask.Render(block)
		} else {
			s += normalTask.Render(block)
		}
	}

	s += "\nPress q to quit.\n"
	s += m.debugInfo()
	return s
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
	normalTask.Width(width).Padding(0, 1, 0).Height(height).Margin(1, 1, 0)
	selectedTask = normalTask.Copy().Background(lipgloss.Color(taskBackgroundColor))
}

func (m model) showMode() string {
	switch m.mode {
	case normalMode:
		return "NOR"
	case insertMode:
		return "INS"
	case selectMode:
		return "SEL"
	}
	return ""
}

func (m model) assertInvariants() {
	if len(m.selected) > 1 {
		panic(fmt.Sprintf("too many elements selected! want 1 have %v", len(m.selected)))
	}

	if m.hasSelectedBlock() && m.mode == insertMode {
		panic(fmt.Sprintf("selected block while editing! selected block at index %v. cursor at %v", m.selected, m.cursor))
	}
}

func (m model) debugInfo() string {
	return fmt.Sprintf("\n%s | height: %v | width: %v \n", m.showMode(), m.height, m.width)
}

// conv24To12 converts a 24 hour timestamp into a 12 hour clock time string.
// time24 represents the hour of the day and must be between 0.0 and 24.0.
func conv24To12(time24 float64) string {
	integer, fraction := math.Modf(time24)
	hrs := int(integer) % 12
	mins := math.Floor(fraction * 60)

	var period string
	if time24 < 12 {
		period = "am"
	} else {
		period = "pm"
	}
	return fmt.Sprintf("%v:%02v %s", hrs, mins, period)
}

func makeBlockLabels(numBlocks int) []string {
	labels := make([]string, numBlocks)

	time := float64(dayStartTime)
	interval := float64(1) / float64(blocksPerHour)
	for i := 0; i < len(labels); i++ {
		labels[i] = conv24To12(time)
		time += float64(interval)
	}
	return labels
}
