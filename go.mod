module linhx.com/tbmk

go 1.19

replace linhx.com/tbmk/common => ./common

replace linhx.com/tbmk/bookmark => ./bookmark

replace linhx.com/tbmk/views/save => ./views/save

replace linhx.com/tbmk/views/search => ./views/search

require (
	github.com/charmbracelet/bubbletea v0.23.0
	linhx.com/tbmk/bookmark v0.0.0-00010101000000-000000000000
	linhx.com/tbmk/views/save v0.0.0-00010101000000-000000000000
	linhx.com/tbmk/views/search v0.0.0-00010101000000-000000000000
)

require (
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52 v1.0.3 // indirect
	github.com/charmbracelet/bubbles v0.14.0 // indirect
	github.com/charmbracelet/lipgloss v0.5.0 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/gookit/color v1.5.2 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/muesli/ansi v0.0.0-20211018074035-2e021307bc4b // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.13.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sahilm/fuzzy v0.1.0 // indirect
	github.com/sonyarouje/simdb v0.1.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	linhx.com/tbmk/common v0.0.0-00010101000000-000000000000 // indirect
)
