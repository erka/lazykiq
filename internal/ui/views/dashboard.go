package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Dashboard is the main overview view
type Dashboard struct {
	width  int
	height int
	styles Styles
}

// NewDashboard creates a new Dashboard view
func NewDashboard() *Dashboard {
	return &Dashboard{}
}

// Init implements View
func (d *Dashboard) Init() tea.Cmd {
	return nil
}

// Update implements View
func (d *Dashboard) Update(msg tea.Msg) (View, tea.Cmd) {
	return d, nil
}

// View implements View
func (d *Dashboard) View() string {
	content := d.styles.Text.Render("Overview of Sidekiq status will appear here.") + "\n\n" +
		d.styles.Muted.Render("Press 1-6 to switch views, t to toggle theme")

	return d.renderBorderedBox(d.Name(), content, d.width, d.height)
}

// Name implements View
func (d *Dashboard) Name() string {
	return "Dashboard"
}

// ShortHelp implements View
func (d *Dashboard) ShortHelp() []key.Binding {
	return nil
}

// SetSize implements View
func (d *Dashboard) SetSize(width, height int) View {
	d.width = width
	d.height = height
	return d
}

// SetStyles implements View
func (d *Dashboard) SetStyles(styles Styles) View {
	d.styles = styles
	return d
}

// renderBorderedBox renders content in a box with title on the top border.
func (d *Dashboard) renderBorderedBox(title, content string, width, height int) string {
	if width < 4 {
		width = 4
	}
	if height < 3 {
		height = 3
	}

	border := lipgloss.RoundedBorder()
	titleText := " " + d.styles.Title.Render(title) + " "
	titleWidth := lipgloss.Width(titleText)

	innerWidth := width - 2
	contentHeight := height - 2

	topLeft := d.styles.BorderStyle.Render(border.TopLeft)
	topRight := d.styles.BorderStyle.Render(border.TopRight)
	hBar := d.styles.BorderStyle.Render(border.Top)

	remainingWidth := width - 2 - titleWidth
	leftPad := 1
	rightPad := remainingWidth - leftPad
	if rightPad < 0 {
		rightPad = 0
	}

	topBorder := topLeft + strings.Repeat(hBar, leftPad) + titleText + strings.Repeat(hBar, rightPad) + topRight

	vBar := d.styles.BorderStyle.Render(border.Left)
	vBarRight := d.styles.BorderStyle.Render(border.Right)

	lines := strings.Split(content, "\n")
	middleLines := make([]string, 0, contentHeight)
	for i := 0; i < contentHeight; i++ {
		line := ""
		if i < len(lines) {
			line = lines[i]
		}
		line = " " + line + " "
		lineWidth := lipgloss.Width(line)
		if lineWidth < innerWidth {
			line += strings.Repeat(" ", innerWidth-lineWidth)
		}
		middleLines = append(middleLines, vBar+line+vBarRight)
	}

	bottomLeft := d.styles.BorderStyle.Render(border.BottomLeft)
	bottomRight := d.styles.BorderStyle.Render(border.BottomRight)
	bottomBorder := bottomLeft + strings.Repeat(hBar, innerWidth) + bottomRight

	return topBorder + "\n" + strings.Join(middleLines, "\n") + "\n" + bottomBorder
}
