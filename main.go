package main

import (
	"encoding/json"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/crosleyzack/bubbles/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	styleDoc = lipgloss.NewStyle().Padding(1)
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var RootCmd = &cobra.Command{}

func init() {
	RootCmd.AddCommand(GetRunCmd())
}

func GetRunCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:     "run",
		Example: "run --file data.json",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			w, h, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				w = 80
				h = 24
			}
			top, right, bottom, left := styleDoc.GetPadding()
			w = w - left - right
			h = h - top - bottom
			content, err := os.ReadFile(file)
			if err != nil {
				log.Fatal("Error when opening file: ", err)
			}
			// Now let's unmarshall the data into `payload`
			var result utils.JsonBlob
			err = json.Unmarshal(content, &result)
			if err != nil {
				log.Fatal("Error during Unmarshal(): ", err)
			}
			model := result.Treeify()
			model.SetHeight(h)
			model.SetWidth(w)
			// model.KeyMap = tree.KeyMap{
			// 	Down: key.NewBinding(
			// 		key.WithKeys("j", "down"),
			// 		key.WithHelp("↓", "down"),
			// 	),
			// 	Up: key.NewBinding(
			// 		key.WithKeys("k", "up"),
			// 		key.WithHelp("↑", "up"),
			// 	),
			// }
			program := tea.NewProgram(utils.NewModel(model))
			_, err = program.Run()
			if err != nil {
				log.Fatal("Error during program start: ", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "JSON file to display")
	return cmd
}
