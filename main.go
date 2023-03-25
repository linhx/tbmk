package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"linhx.com/tbmk/bookmark"
	saveView "linhx.com/tbmk/views/save"
	searchView "linhx.com/tbmk/views/search"
)

func NewCancellationSignal() (func(), func()) {
	canceled := false

	cancel := func() {
		canceled = true
	}
	exit := func() {
		if canceled {
			os.Exit(1)
		}
	}

	return cancel, exit
}

func main() {
	_, exit := NewCancellationSignal()
	defer exit()
	saveCmd := flag.NewFlagSet("save", flag.ExitOnError)
	saveCommand := saveCmd.String("command", "", "command")

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchQuery := searchCmd.String("query", "", "query")

	if len(os.Args) < 2 {
		fmt.Println("expected 'save' or 'search' subcommands")
		exit()
	}

	bmk, err := bookmark.NewBookmark()
	switch os.Args[1] {
	case "save":
		saveCmd.Parse(os.Args[2:])
		p := tea.NewProgram(saveView.InitialModel(*bmk, *saveCommand), tea.WithOutput(os.Stderr))
		_, err := p.Run()
		if err != nil {
			log.Fatal(err)
		}
	case "search":
		searchCmd.Parse(os.Args[2:])
		p := tea.NewProgram(searchView.InitialModel(*bmk, *searchQuery), tea.WithOutput(os.Stderr))
		m, err := p.Run()
		if err != nil {
			log.Fatal(err)
		}
		selectedCommand := m.(searchView.Model).SelectedItem.Command
		if len(selectedCommand) > 0 {
			fmt.Print(selectedCommand)
		} else {
			fmt.Print(*searchQuery)
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}
