package metrics

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func testStyles() Styles {
	return Styles{
		Bar:       lipgloss.NewStyle(),
		Fill:      lipgloss.NewStyle(),
		Label:     lipgloss.NewStyle(),
		Value:     lipgloss.NewStyle(),
		Separator: lipgloss.NewStyle(),
	}
}

func testData() Data {
	return Data{
		Processed: 1,
		Failed:    22,
		Busy:      333,
		Enqueued:  4444,
		Retries:   55555,
		Scheduled: 6,
		Dead:      7777777,
	}
}

func TestViewSnapshots(t *testing.T) {
	data := testData()
	cases := []struct {
		name     string
		width    int
		expected string
	}{
		{
			name:     "truncate",
			width:    99,
			expected: " Processed: 1   Failed: 22   Busy: 333   Enqueued: 4.4K   Retries: 55.6K   Scheduled: 6   Dead: 7.8",
		},
		{
			name:     "min-spacing",
			width:    101,
			expected: " Processed: 1   Failed: 22   Busy: 333   Enqueued: 4.4K   Retries: 55.6K   Scheduled: 6   Dead: 7.8M ",
		},
		{
			name:     "trim-padding",
			width:    116,
			expected: " Processed: 1    Failed: 22      Busy: 333        Enqueued: 4.4K   Retries: 55.6K   Scheduled: 6     Dead: 7.8M     ",
		},
		{
			name:     "equal-width",
			width:    118,
			expected: " Processed: 1     Failed: 22       Busy: 333        Enqueued: 4.4K   Retries: 55.6K   Scheduled: 6     Dead: 7.8M     ",
		},
		{
			name:     "extra-distributed",
			width:    120,
			expected: " Processed: 1      Failed: 22        Busy: 333        Enqueued: 4.4K   Retries: 55.6K   Scheduled: 6     Dead: 7.8M     ",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := New(
				WithStyles(testStyles()),
				WithWidth(tc.width),
				WithData(data),
			)
			if got := m.View(); got != tc.expected {
				t.Fatalf("unexpected output:\nexpected %q\ngot      %q", tc.expected, got)
			}
		})
	}
}
