package main

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

// Normal, insert, select, stretch modes
var modes = []string{"NOR", "INS", "SEL", "STR"}

type styles struct {
	normalBlock   lipgloss.Style
	currentBlock  lipgloss.Style
	selectedBlock lipgloss.Style
	pastBlock     lipgloss.Style
	stretchDebug  lipgloss.Style

	tiTextStyle        lipgloss.Style
	tiPlaceholderStyle lipgloss.Style
}

// func (m model) View() string {
// 	s := "Schedule"
// 	var block string
// 	for i := range m.tasks {
// 		block = fmt.Sprint(m.spans[i])
// 		s += m.styles.normalBlock.Render(block)
// 	}
// 	return s
// }

func (m model) View() string {
	s := "Schedule"

	currentBlockIdx := m.findCurrentTimeBlock()

	// TODO: maybe change this to iterating through spans, since spans
	// define what makes a visible block
	for i, task := range m.tasks {
		cursorIndicator := " " // no cursor
		if m.cursor == i {
			cursorIndicator = ">" // cursor!
		}

		var block string
		// fill the block
		block += fmt.Sprintf("%s\n", m.blockLabels[i])
		// block += fmt.Sprintf("%v\n", m.spans[i])

		if m.mode == insertMode && m.cursor == i {
			block += fmt.Sprintf("%s\n", m.textInput.View())
		} else {
			block += fmt.Sprintf("%s %s\n", cursorIndicator, task)
		}

		blockHeight := int(math.Floor(float64(m.height)/float64(m.numBlocks)) * float64(0.6))
		if blockHeight < 2 {
			blockHeight = 2
		}
		// *m.height = 2
		blockWidth := int(math.Floor(float64(m.width) - float64(float64(m.width)/10.0)))
		if blockWidth < 20 {
			blockWidth = 20
		}

		// width := 10
		// height := 1

		// m.logger.Info("UI", "length", len(s))
		// s += m.styles.normalBlock.Render(block)
		// m.logger.Info("UI", "length", len(s))

		// highlight tasks if selected or current time
		if _, ok := m.selected[i]; ok {
			s += m.styles.selectedBlock.Width(blockWidth).Height(blockHeight).Render(block)
		} else if i == currentBlockIdx {
			s += m.styles.currentBlock.Width(blockWidth).Height(blockHeight).Render(block)
		} else if i < currentBlockIdx {
			s += m.styles.pastBlock.Width(blockWidth).Height(blockHeight).Render(block)
		} else {
			s += m.styles.normalBlock.Width(blockWidth).Height(blockHeight).Render(block)
		}

	}

	s += "\nPress q to quit.\n"
	s += m.debugInfo()
	return s
}

func (m model) debugInfo() string {
	hr := m.findCurrentTimeBlock()
	return fmt.Sprintf("\n%s | height: %v | width: %v | hour: %v | cursor: %v\n", modes[m.mode], m.height, m.width, hr, m.cursor)
	// return fmt.Sprintf("\n%s | hour: %v | cursor: %v | tasks len: %v\n", modes[m.mode], hr, m.cursor, len(m.tasks))
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
	hr, _, _ := m.currentTime.Clock()
	// hr := 12 // dummy to stand in for currentTime.Clock()

	idx := hr - dayStartTime

	return idx
}

func defaultStyles() (s styles) {
	normalBlockBackgroundColor := "0"
	// selectedBlockBackgroundColor := "4"
	selectedBlockBackgroundColor := "8"
	// currentBlockBackgroundColor := "5"
	currentBlockBackgroundColor := "20"

	// fadedTextColor := "16"
	fadedTextColor := "1"

	s.normalBlock = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true).
		Background(lipgloss.Color(normalBlockBackgroundColor)).
		MaxWidth(40).
		Padding(0, 1, 0).
		Margin(1, 1, 0)

	s.currentBlock = s.normalBlock.Copy().
		Background(lipgloss.Color(selectedBlockBackgroundColor))

	s.selectedBlock = s.normalBlock.Copy().
		Background(lipgloss.Color(currentBlockBackgroundColor))

	s.pastBlock = s.normalBlock.Copy().
		Foreground(lipgloss.Color(fadedTextColor))

	s.stretchDebug = s.normalBlock.Copy().
		UnsetBorderBottom().
		UnsetBorderTop().
		Background(lipgloss.Color("5"))

	s.tiTextStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(selectedBlockBackgroundColor))

	s.tiPlaceholderStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(selectedBlockBackgroundColor))

	return s
}
