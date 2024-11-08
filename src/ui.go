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
	workdayBlock  lipgloss.Style

	tiTextStyle        lipgloss.Style
	tiPlaceholderStyle lipgloss.Style

	titleStyle lipgloss.Style
}

func (m model) View() string {
	// TODO: switch to strings.Builder
	s := m.styles.titleStyle.Render("Schedule")

	currentBlockIdx := blockIdxFromTime(m.currentTime)

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

		var style lipgloss.Style

		// highlight tasks if selected or current time
		if i < currentBlockIdx {
			style = m.styles.pastBlock.Copy()
		} else {
			style = m.styles.normalBlock.Copy()
		}

		if i >= m.dayStartBlock && i < (m.dayStartBlock+(m.dayLengthHrs*2)) {
			style = m.styles.workdayBlock.Copy()
		}

		if i == currentBlockIdx {
			style = m.styles.currentBlock.Copy()
		}

		if _, ok := m.selected[i]; ok {
			style = m.styles.selectedBlock.Copy()
		}

		s += style.Width(m.blockWidth).
			Height(m.blockHeight).Render(block)
	}

	s += "\nPress q to quit.\n"
	s += m.debugInfo()
	s += fmt.Sprintf("\nvpend: %v\n", vpEnd)

	return s
}

func (m model) debugInfo() string {
	return fmt.Sprintf("\n%s | height: %v | width: %v | hour: %v | cursor: %v | block height: %v | num blocks: %v \n", modes[m.mode], m.height, m.width, m.currentTime, m.cursor, m.blockHeight, m.numBlocks)
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

func makeBlockLabels(numBlocks, blocksPerHour int) []string {
	labels := make([]string, numBlocks)

	time := 0.0
	interval := float64(1) / float64(blocksPerHour)
	for i := 0; i < len(labels); i++ {
		labels[i] = conv24To12(time)
		time += float64(interval)
	}
	return labels
}

func defaultStyles() (s styles) {
	normalBlockBackgroundColor := "8"
	selectedBlockBackgroundColor := "4"
	currentBlockBackgroundColor := "5"
	fadedTextColor := "16"
	workdayBlockColor := "1"

	s.normalBlock = lipgloss.NewStyle().
		Background(lipgloss.Color(normalBlockBackgroundColor)).
		MaxWidth(40).
		Margin(1, 0, 0)

	s.currentBlock = s.normalBlock.Copy().
		Background(lipgloss.Color(currentBlockBackgroundColor))

	s.selectedBlock = s.normalBlock.Copy().
		Background(lipgloss.Color(selectedBlockBackgroundColor))

	s.pastBlock = s.normalBlock.Copy().
		Foreground(lipgloss.Color(fadedTextColor))

	s.workdayBlock = s.normalBlock.Copy().
		Background(lipgloss.Color(workdayBlockColor))

	s.stretchDebug = s.normalBlock.Copy().
		UnsetBorderBottom().
		UnsetBorderTop().
		Background(lipgloss.Color("5"))

	s.tiTextStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(selectedBlockBackgroundColor))

	s.tiPlaceholderStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(selectedBlockBackgroundColor))

	s.titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)

	return s
}

func changeStyles(config configMsg) *styles {
	normalBlockBackgroundColor := config["normalBlockBackgroundColor"]
	selectedBlockBackgroundColor := config["selectedBlockBackgroundColor"]
	currentBlockBackgroundColor := config["currentBlockBackgroundColor"]
	fadedTextColor := config["fadedTextColor"]
	workdayBlockColor := config["workdayBlockColor"]

	s := defaultStyles()

	s.normalBlock = s.normalBlock.UnsetBackground().Background(lipgloss.Color(normalBlockBackgroundColor))
	s.currentBlock = s.currentBlock.UnsetBackground().Background(lipgloss.Color(currentBlockBackgroundColor))
	s.selectedBlock = s.selectedBlock.UnsetBackground().Background(lipgloss.Color(selectedBlockBackgroundColor))
	s.pastBlock = s.pastBlock.UnsetForeground().Foreground(lipgloss.Color(fadedTextColor))
	s.workdayBlock = s.workdayBlock.UnsetBackground().Background(lipgloss.Color(workdayBlockColor))
	return &s
}
