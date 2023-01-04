/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/PNCommand/dstm/tui/component"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		T03()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func T01() {
	editor, err := component.NewWorldEditor("./json/zh-CN/forest.gen.zh-CN.json", true, true)
	if err != nil {
		fmt.Println("========")
		fmt.Println(err)
		os.Exit(1)
	}
	model, err := tea.NewProgram(editor).Run()
	if err != nil {
		fmt.Println("========")
		fmt.Println(err)
		os.Exit(1)
	}
	for _, group := range model.(*component.WorldEditor).ConfigGroups {
		for _, item := range group.Items {
			if item.Index == nil {
				continue
			}
			fmt.Println(item.Display, item.OptsDisplay[*item.Index])
		}
	}
}

func T02() {
	cfg, err := ini.Load("./_test_data/cluster.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	sections := cfg.Sections()
	for _, sec := range sections {
		keys := sec.Keys()
		for _, key := range keys {
			fmt.Println(sec.Name(), ":", key.Value())
		}
	}
}

func T03() {
	editor, err := component.NewIniEditor(false, "ini editor")
	if err != nil {
		fmt.Println("========")
		fmt.Println(err)
		os.Exit(1)
	}
	model, err := tea.NewProgram(editor).Run()
	if err != nil {
		fmt.Println("========")
		fmt.Println(err)
		os.Exit(1)
	}
	for _, group := range model.(*component.IniEditor).IniGroups {
		for _, item := range group.Items {
			fmt.Println(item.Key + ":" + item.Value + ",")
		}
	}
}
