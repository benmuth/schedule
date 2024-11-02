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

func (m model) View() string {
	// TODO: switch to strings.Builder
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)

	s := titleStyle.Render("Schedule")

	currentBlockIdx := m.findCurrentTimeBlock()

	// TODO: maybe change this to iterating through spans, since spans
	// define what makes a visible block

	vpEnd := m.calcVPEnd()
	for i := m.vpStart; i < vpEnd; i++ {
		cursorIndicator := " " // no cursor
		if m.cursor == i {
			cursorIndicator = ">" // cursor!
		}

		var block string
		// fill the block
		block += fmt.Sprintf("%s ", m.blockLabels[i])

		if m.mode == insertMode && m.cursor == i {
			block += fmt.Sprintf("%s\n", m.textInput.View())
		} else {
			block += fmt.Sprintf("%s %s\n", cursorIndicator, m.tasks[i])
		}

		// highlight tasks if selected or current time
		if _, ok := m.selected[i]; ok {
			s += m.styles.selectedBlock.Width(m.blockWidth).Height(m.blockHeight).Render(block)
		} else if i == currentBlockIdx {
			s += m.styles.currentBlock.Width(m.blockWidth).Height(m.blockHeight).Render(block)
		} else if i < currentBlockIdx {
			s += m.styles.pastBlock.Width(m.blockWidth).Height(m.blockHeight).Render(block)
		} else {
			s += m.styles.normalBlock.Width(m.blockWidth).Height(m.blockHeight).Render(block)
		}
	}

	s += "\nPress q to quit.\n"
	s += m.debugInfo()
	s += fmt.Sprintf("\nvpend: %v\n", vpEnd)

	return s
}

func (m model) debugInfo() string {
	hr := m.findCurrentTimeBlock()
	return fmt.Sprintf("\n%s | height: %v | width: %v | hour: %v | cursor: %v | block height: %v | num blocks: %v \n", modes[m.mode], m.height, m.width, hr, m.cursor, m.blockHeight, m.numBlocks)
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

	idx := hr - dayStartTime

	return idx
}

func defaultStyles() (s styles) {
	normalBlockBackgroundColor := "8"
	selectedBlockBackgroundColor := "4"
	currentBlockBackgroundColor := "5"
	fadedTextColor := "16"

	s.normalBlock = lipgloss.NewStyle().
		Background(lipgloss.Color(normalBlockBackgroundColor)).
		MaxWidth(40).
		Margin(1, 0, 0)

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
