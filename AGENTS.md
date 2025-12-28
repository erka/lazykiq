# Lazykiq - Agents Context

## What

Bubble Tea TUI for Sidekiq monitoring. Go 1.25.

## Structure

```text
cmd/lazykiq/main.go          - entry point
internal/
  sidekiq/
    client.go                - Redis client, GetStats(), uses go-redis with VoidLogger
  ui/
    app.go                   - main model, renderBorderedBox() for titled border
    keys.go                  - KeyMap struct, DefaultKeyMap()
    theme/theme.go           - Theme (Dark/Light), Styles, NewStyles()
    components/
      metrics.go             - top bar: Processed|Failed|Busy|Enqueued|Retries|Scheduled|Dead
      navbar.go              - bottom bar: view keys + quit + theme
      error_popup.go         - centered overlay for connection errors
      table/table.go         - reusable scrollable table with selection
    format/format.go         - Duration, Bytes, Args, Number formatters
    views/
      view.go                - View interface, views.Styles
      dashboard.go, queues.go, busy.go, retries.go, scheduled.go, dead.go
```

## Patterns

- Views implement `views.View` interface
- Components take `*theme.Styles`, have SetStyles()
- Theme uses AdaptiveColor; no runtime toggle
- All color values must live in `theme.DefaultTheme`; no inline colors outside it
- Border title: renderBorderedBox() in app.go
- No backgrounds on metrics/navbar (transparent)
- NO EMOJIS in UI - keep text clean and professional
- Shared components: no lipgloss.NewStyle() calls - pass all styles via struct
- Table in `components/table/` subpackage to avoid import cycle (components ↔ views)
- Table: last column not truncated/padded to allow horizontal scroll of variable content

## Component Pattern (bubbles-style)

Follow the charmbracelet/bubbles pattern for reusable components:

```go
// 1. Styles struct - exported, holds all styles
type Styles struct {
    Title  lipgloss.Style
    Border lipgloss.Style
}

// 2. DefaultStyles() - returns sensible defaults
func DefaultStyles() Styles {
    return Styles{
        Title:  lipgloss.NewStyle().Bold(true),
        Border: lipgloss.NewStyle(),
    }
}

// 3. Model struct - holds all state
type Model struct {
    styles  Styles  // unexported, use SetStyles()
    width   int
    height  int
    // ... other state
}

// 4. Option type for functional options
type Option func(*Model)

// 5. New() constructor with functional options
func New(opts ...Option) Model {
    m := Model{
        styles: DefaultStyles(),
        // ... defaults
    }
    for _, opt := range opts {
        opt(&m)
    }
    return m
}

// 6. WithXxx() option functions
func WithStyles(s Styles) Option {
    return func(m *Model) { m.styles = s }
}

func WithSize(w, h int) Option {
    return func(m *Model) { m.width, m.height = w, h }
}

// 7. SetXxx() methods - pointer receiver, for post-creation updates
func (m *Model) SetStyles(s Styles) { m.styles = s }
func (m *Model) SetSize(w, h int)   { m.width, m.height = w, h }

// 8. Getter methods - value receiver
func (m Model) Width() int  { return m.width }
func (m Model) Height() int { return m.height }

// 9. Update() - value receiver, returns new Model + Cmd (if interactive)
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // handle messages
    return m, nil
}

// 10. View() - value receiver, renders to string
func (m Model) View() string {
    // render
}
```

Key principles:

- **Value receivers** for `Update()`, `View()`, getters (immutable operations)
- **Pointer receivers** for `SetXxx()` methods (mutations)
- **Unexported state** accessed via `SetXxx()`/getters
- **`DefaultStyles()`** so components work without explicit styling
- **Functional options** for clean initialization
- **No `lipgloss.NewStyle()`** in render methods - styles passed in

## Data Flow

- 5-second ticker fetches Sidekiq stats from Redis
- MetricsUpdateMsg updates metrics bar
- connectionErrorMsg shows error popup overlay
- Error popup auto-clears when Redis reconnects

## Keys

1-6: views, t: theme, q: quit, tab/shift+tab: reserved

## Dependencies

bubbletea, lipgloss, bubbles/key, go-redis/v9

## Gotchas

- Horizontal scroll: apply offset to plain text BEFORE lipgloss styling. Slicing ANSI-escaped strings breaks escape sequences.
- Scroll state: clamp xOffset/yOffset when data or dimensions change (new data may have different max width)
- Manual vertical scroll (line slicing) is simpler than bubbles/viewport for tables with selection
- Filtered sorted-set scans use ZSCAN; always sort matches by score to preserve chronological order (dead: newest first; retry/scheduled: earliest first).
- Textinput placeholder rendering needs Width set; otherwise only the first placeholder rune appears.
- When an input component is focused, the app must route key events to the view before global shortcuts to avoid stealing keys (e.g., theme toggle).
- Height calculations: app.go renders metrics bar (top) + view content + navbar (bottom). Views must output exactly the same number of lines consistently. If view outputs too many lines, metrics bar gets pushed off screen. Specific issues:
  - Title is part of the border line, not a separate line (so -2 for borders, not -3)
  - Views with header areas outside the main box (Busy, Queues) get extra height (+3 instead of +2). When showing alternative content (like job detail), must output the same total lines as normal view - add empty lines at top if needed to match the header area
