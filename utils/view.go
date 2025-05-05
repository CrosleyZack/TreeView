package utils

// Taken from https://github.com/savannahostrowski/tree-bubble/blob/main/example/main.go

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/crosleyzack/bubbles/tree"
)

var (
	styleDoc = lipgloss.NewStyle().Padding(1)
)

// NewModel creates a new model with the given tree.
func NewModel(tree tree.Model) model {
	// set top level nodes to expanded
	for _, node := range tree.Nodes() {
		node.Expand = true
	}
	return model{tree: tree}
}

type model struct {
	tree tree.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.tree, cmd = m.tree.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return styleDoc.Render(m.tree.View())
}
