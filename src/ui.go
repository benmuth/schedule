package main

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

const (
	selectedBlockBackgroundColor = "4"
	currentBlockBackgroundColor  = "5"
)

// styles
var (
	normalBlock  = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true).Background(lipgloss.Color("0")).MaxWidth(40).Padding(0, 1, 0).Margin(1, 1, 0)
	currentBlock = normalBlock.Copy().Background(lipgloss.Color(selectedBlockBackgroundColor))
	// currentBlock = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, false, true).Background(lipgloss.Color("5")).MaxWidth(40)
	// selectedBlock = normalBlock.Copy().Background(lipgloss.Color(currentBlockBackgroundColor))
	selectedBlock = normalBlock.Copy().Background(lipgloss.Color(currentBlockBackgroundColor))
	// test             = selectedBlock
)

func (m model) View() string {
	s := "Schedule"

	currentBlockIdx := m.findCurrentTimeBlock()

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

		// highlight tasks if selected or current time
		if _, ok := m.selected[i]; ok {
			s += selectedBlock.Render(block)
		} else if i == currentBlockIdx {
			// s += currentTimeBlock.Render(block)
			// s += selectedBlock.Render(block)
			s += currentBlock.Render(block)
		} else {
			s += normalBlock.Render(block)
		}

	}

	s += "\nPress q to quit.\n"
	s += m.debugInfo()
	return s
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
	hr := m.findCurrentTimeBlock()
	return fmt.Sprintf("\n%s | height: %v | width: %v | hour: %v \n", m.showMode(), m.height, m.width, hr)
}

// conv24To12 converts a 24 hour timestamp into a 12 hour clock time string.
// time24 represents the hour of the day and must be between 0.0 and 24.0.
func conv24To12(time24 float64) string {
	integer, fraction := math.Modf(time24)
	hrs := int(integer) % 12
	if hrs == 0 {
		hrs = 12
	}
	mins := math.Floor(fraction * 60)

	var period string
	if time24 < 12 {
		period = "am"
	} else {
		period = "pm"
	}
	return fmt.Sprintf("%v:%02v %s", hrs, mins, period)
}

func makeBlockLabels(numBlocks, startTime, blocksPerHour int) []string {
	labels := make([]string, numBlocks)

	time := float64(startTime)
	interval := float64(1) / float64(blocksPerHour)
	for i := 0; i < len(labels); i++ {
		labels[i] = conv24To12(time)
		time += float64(interval)
	}
	return labels
}

func (m model) findCurrentTimeBlock() int {
	// hr, _, _ := m.currentTime.Clock()
	hr := 15 // dummy to stand in for currentTime.Clock()

	idx := hr - dayStartTime

	return idx
}
