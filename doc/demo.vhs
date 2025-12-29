Output doc/assets/demo.gif

Set Shell "bash"
Set Framerate 24
Set FontFamily "JetBrainsMono Nerd Font"
# Doubled values for better quality
# See https://github.com/charmbracelet/vhs/issues/69#issuecomment-1295581303
Set FontSize 28
Set Width 2400
Set Height 1600
Set Padding 40

Hide
Type@1ms "go run ./cmd/lazykiq" Enter
Sleep 1s
Show

# Dashboard: stay 10s, but only record a few frames
Sleep 0.5s
Hide
Sleep 10s
Show
Sleep 0.5s
Hide
Sleep 10s
Show
Sleep 0.5s
Hide
Sleep 10s
Show
Sleep 2s

# Busy
Type "2" Sleep 1.5s

# Queues
Type "3" Sleep 1s
Down Sleep 0.5s Down Sleep 0.5s Enter Sleep 2s
Escape Sleep 0.5s

# Retries
Type "4" Sleep 0.5s
Type "]" Sleep 0.5s Type "]" Sleep 2s

# Scheduled
Type "5" Sleep 0.5s
Type "/" Sleep 0.5s Type "Data" Sleep 0.5s Enter Sleep 1s
Enter Sleep 2s
Escape Sleep 0.5s

# Dead
Type "6" Sleep 1s
Enter Sleep 2s
Escape Sleep 0.5s
