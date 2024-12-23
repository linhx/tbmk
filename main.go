/*
 * This file is part of tbmk.
 *
 * tbmk is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * tbmk is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with tbmk.  If not, see <https://www.gnu.org/licenses/>.
 *
 * Copyright (C) 2024 Nguyen Dinh Linh
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/viper"
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

func getAppDir() string {
	if os.Getenv("APP_ENV") == "dev" {
		return "."
	}
	path, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("Cannot get executable: %w", err))
	}
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		panic(fmt.Errorf("Cannot get real path from symlink: %w", err))
	}
	return filepath.Dir(realPath)
}

func main() {
	lipgloss.SetColorProfile(termenv.TrueColor)
	_, exit := NewCancellationSignal()
	defer exit()
	viper.AddConfigPath(getAppDir())
	viper.SetConfigName("config")
	viper.SetDefault("tbmk.dataDir", "./data")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	saveCmd := flag.NewFlagSet("save", flag.ExitOnError)
	saveCommand := saveCmd.String("command", "", "command")

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchQuery := searchCmd.String("query", "", "query")

	if len(os.Args) < 2 {
		fmt.Println("expected 'save' or 'search' subcommands")
		exit()
		return
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
		selectedCommand := m.(searchView.Model).OutputCommand
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
