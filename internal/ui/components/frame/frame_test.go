package frame

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestFrameLineCountAndWidth(t *testing.T) {
	box := New(
		WithSize(10, 4),
		WithTitle("T"),
		WithTitlePadding(1),
		WithContent("hi"),
	)

	view := box.View()
	lines := strings.Split(view, "\n")
	if len(lines) != 4 {
		t.Fatalf("want 4 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if lipgloss.Width(line) != 10 {
			t.Fatalf("line %d: want width 10, got %d", i, lipgloss.Width(line))
		}
	}
}

func TestFrameMinHeight(t *testing.T) {
	box := New(
		WithSize(10, 2),
		WithMinHeight(5),
		WithTitle("T"),
		WithTitlePadding(1),
		WithContent("hi"),
	)

	view := box.View()
	lines := strings.Split(view, "\n")
	if len(lines) != 5 {
		t.Fatalf("want 5 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if lipgloss.Width(line) != 10 {
			t.Fatalf("line %d: want width 10, got %d", i, lipgloss.Width(line))
		}
	}
}

func TestFrameFocusStyles(t *testing.T) {
	focusedBorder := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	styles := Styles{
		Focused: StyleState{
			Title:  lipgloss.NewStyle(),
			Border: focusedBorder,
		},
		Blurred: StyleState{
			Title:  lipgloss.NewStyle(),
			Border: lipgloss.NewStyle(),
		},
	}

	focused := New(
		WithStyles(styles),
		WithFocused(true),
		WithTitle("T"),
		WithSize(8, 3),
	)
	unfocused := New(
		WithStyles(styles),
		WithFocused(false),
		WithTitle("T"),
		WithSize(8, 3),
	)

	if !strings.Contains(focused.View(), "\x1b[") {
		t.Fatalf("expected focused view to contain ANSI sequences")
	}
	if strings.Contains(unfocused.View(), "\x1b[") {
		t.Fatalf("expected unfocused view to avoid ANSI sequences")
	}
}
