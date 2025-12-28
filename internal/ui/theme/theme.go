package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines all colors used throughout the UI.
type Theme struct {
	// Base colors
	Primary   lipgloss.AdaptiveColor
	Secondary lipgloss.AdaptiveColor

	// Text colors
	Text       lipgloss.AdaptiveColor
	TextMuted  lipgloss.AdaptiveColor
	TextBright lipgloss.AdaptiveColor

	// Background colors
	Bg           lipgloss.AdaptiveColor
	BgAlt        lipgloss.AdaptiveColor
	MetricsBarBg lipgloss.AdaptiveColor

	// Border colors
	Border      lipgloss.AdaptiveColor
	BorderFocus lipgloss.AdaptiveColor

	// Accent colors
	TableSelectedFg lipgloss.AdaptiveColor
	TableSelectedBg lipgloss.AdaptiveColor
	Success         lipgloss.AdaptiveColor
	Error           lipgloss.AdaptiveColor

	// Metrics colors
	MetricsText  lipgloss.AdaptiveColor
	MetricsSepBg lipgloss.AdaptiveColor
}

// DefaultTheme is the adaptive color scheme used by default.
// Use Open Color palette when possible to define colors: https://yeun.github.io/open-color/
var DefaultTheme = Theme{
	// Sidekiq-inspired primary
	Primary: lipgloss.AdaptiveColor{
		Light: "#B2003C",
		Dark:  "#F73D68",
	},
	Secondary: lipgloss.AdaptiveColor{
		Light: "#6B7280", // Gray-500
		Dark:  "#6B7280", // Gray-500
	},

	// Text
	Text: lipgloss.AdaptiveColor{
		Light: "#111827", // Gray-900
		Dark:  "#F9FAFB", // Gray-50
	},
	TextMuted: lipgloss.AdaptiveColor{
		Light: "#6B7280", // Gray-500
		Dark:  "#9CA3AF", // Gray-400
	},
	TextBright: lipgloss.AdaptiveColor{
		Light: "#030712", // Gray-950
		Dark:  "#FFFFFF", // White
	},

	// Backgrounds
	Bg: lipgloss.AdaptiveColor{
		Light: "#FFFFFF", // White
		Dark:  "#111827", // Gray-900
	},
	BgAlt: lipgloss.AdaptiveColor{
		Light: "#F3F4F6", // Gray-100
		Dark:  "#1F2937", // Gray-800
	},
	MetricsBarBg: lipgloss.AdaptiveColor{
		Light: "#1c7ed6", // blue 7
		Dark:  "#4dabf7", // blue 4
	},

	// Borders
	Border: lipgloss.AdaptiveColor{
		Light: "#D1D5DB", // Gray-300
		Dark:  "#374151", // Gray-700
	},
	BorderFocus: lipgloss.AdaptiveColor{
		Light: "#9CA3AF", // Gray-400
		Dark:  "#6B7280", // Gray-500
	},

	// Accents
	TableSelectedFg: lipgloss.AdaptiveColor{
		Light: "229",
		Dark:  "229",
	},
	TableSelectedBg: lipgloss.AdaptiveColor{
		Light: "57",
		Dark:  "57",
	},
	Success: lipgloss.AdaptiveColor{
		Light: "#16A34A",
		Dark:  "#22C55E",
	},
	Error: lipgloss.AdaptiveColor{
		Light: "#FF0000",
		Dark:  "#FF0000",
	},

	// Metrics
	MetricsText: lipgloss.AdaptiveColor{
		Light: "#f8f9fa",
		Dark:  "#212529", //gray 9
	},
	MetricsSepBg: lipgloss.AdaptiveColor{
		Light: "#1971c2", // blue 8
		Dark:  "#339af0", // blue 5
	},
}

// Styles holds all lipgloss styles derived from a theme
type Styles struct {
	// Metrics bar
	MetricsBar   lipgloss.Style
	MetricsFill  lipgloss.Style
	MetricsLabel lipgloss.Style
	MetricsValue lipgloss.Style
	MetricsSep   lipgloss.Style
	MetricLabel  lipgloss.Style
	MetricValue  lipgloss.Style

	// Navbar
	NavBar  lipgloss.Style
	NavItem lipgloss.Style
	NavKey  lipgloss.Style
	NavQuit lipgloss.Style

	// Content
	ViewTitle lipgloss.Style
	ViewText  lipgloss.Style
	ViewMuted lipgloss.Style

	// Table
	TableHeader    lipgloss.Style
	TableSelected  lipgloss.Style
	TableSeparator lipgloss.Style

	// Layout helpers
	BoxPadding  lipgloss.Style
	BorderStyle lipgloss.Style
	FocusBorder lipgloss.Style

	// Charts
	ChartSuccess lipgloss.Style
	ChartFailure lipgloss.Style

	// Errors
	ErrorTitle  lipgloss.Style
	ErrorBorder lipgloss.Style
}

// NewStyles creates a Styles instance from the default adaptive theme.
func NewStyles() Styles {
	t := DefaultTheme
	return Styles{
		// Metrics bar
		MetricsBar: lipgloss.NewStyle().
			Foreground(t.MetricsText).
			Background(t.MetricsBarBg).
			Padding(0, 0),

		MetricsFill: lipgloss.NewStyle().
			Background(t.MetricsBarBg),

		MetricsLabel: lipgloss.NewStyle().
			Foreground(t.MetricsText).
			Background(t.MetricsBarBg),

		MetricsValue: lipgloss.NewStyle().
			Foreground(t.MetricsText).
			Background(t.MetricsBarBg).
			Bold(true),

		MetricsSep: lipgloss.NewStyle().
			Background(t.MetricsSepBg),

		MetricLabel: lipgloss.NewStyle().
			Foreground(t.TextMuted),

		MetricValue: lipgloss.NewStyle().
			Foreground(t.Text).
			Bold(true),

		// Navbar
		NavBar: lipgloss.NewStyle().
			Padding(0, 1),

		NavItem: lipgloss.NewStyle().
			Foreground(t.TextMuted).
			PaddingRight(1),

		NavKey: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Border).
			Padding(0, 1),

		NavQuit: lipgloss.NewStyle().
			Foreground(t.TextMuted).
			PaddingRight(1),

		// Content
		ViewTitle: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),

		ViewText: lipgloss.NewStyle().
			Foreground(t.Text),

		ViewMuted: lipgloss.NewStyle().
			Foreground(t.TextMuted),

		// Table
		TableHeader: lipgloss.NewStyle().
			Foreground(t.Text).
			Bold(true),

		TableSelected: lipgloss.NewStyle().
			Foreground(t.TableSelectedFg).
			Background(t.TableSelectedBg),

		TableSeparator: lipgloss.NewStyle().
			Foreground(t.Border),

		// Layout helpers
		BoxPadding: lipgloss.NewStyle().
			Padding(0, 1),

		BorderStyle: lipgloss.NewStyle().
			Foreground(t.Border),

		FocusBorder: lipgloss.NewStyle().
			Foreground(t.BorderFocus),

		ChartSuccess: lipgloss.NewStyle().
			Foreground(t.Success),

		ChartFailure: lipgloss.NewStyle().
			Foreground(t.Error),

		ErrorTitle: lipgloss.NewStyle().
			Foreground(t.Error).
			Bold(true),

		ErrorBorder: lipgloss.NewStyle().
			Foreground(t.Error),
	}
}
