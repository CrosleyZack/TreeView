package tree

// copied from https://github.com/savannahostrowski/tree-bubble/blob/main/tree.go

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
)

const (
	bottomLeft string = " └──"

	white  = lipgloss.Color("#ffffff")
	black  = lipgloss.Color("#000000")
	purple = lipgloss.Color("#bd93f9")
)

type Styles struct {
	Shapes     lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	Help       lipgloss.Style
}

func defaultStyles() Styles {
	return Styles{
		Shapes:     lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(purple),
		Selected:   lipgloss.NewStyle().Margin(0, 0, 0, 0).Background(purple),
		Unselected: lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}),
		Help:       lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}),
	}
}

type Node struct {
	Value string
	// Desc is used to store the shorthand for the collapsed values
	Desc     string
	Children []*Node
	Expand   bool
}

type Model struct {
	KeyMap KeyMap
	Styles Styles

	width  int
	height int
	nodes  []*Node
	cursor int

	currentNode *Node

	Help     help.Model
	showHelp bool

	AdditionalShortHelpKeys func() []key.Binding
}

func New(nodes []*Node, width int, height int) *Model {
	return &Model{
		KeyMap: DefaultKeyMap(),
		Styles: defaultStyles(),

		width:  width,
		height: height,
		nodes:  nodes,

		showHelp: true,
		Help:     help.New(),
	}
}

// KeyMap holds the key bindings for the table.
type KeyMap struct {
	Bottom      key.Binding
	Top         key.Binding
	SectionDown key.Binding
	SectionUp   key.Binding
	Down        key.Binding
	Up          key.Binding
	Quit        key.Binding
	Collapse    key.Binding

	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
}

// DefaultKeyMap is the default key bindings for the table.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Bottom: key.NewBinding(
			key.WithKeys("bottom"),
			key.WithHelp("end", "bottom"),
		),
		Top: key.NewBinding(
			key.WithKeys("top"),
			key.WithHelp("home", "top"),
		),
		SectionDown: key.NewBinding(
			key.WithKeys("secdown"),
			key.WithHelp("secdown", "section down"),
		),
		SectionUp: key.NewBinding(
			key.WithKeys("secup"),
			key.WithHelp("secup", "section up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓", "down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑", "up"),
		),
		Collapse: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("tab", "collapse"),
		),

		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (m Model) Nodes() []*Node {
	return m.nodes
}

func (m *Model) SetNodes(nodes []*Node) {
	m.nodes = nodes
}

func (m *Model) NumberOfNodes() int {
	count := 0

	var countNodes func([]*Node)
	countNodes = func(nodes []*Node) {
		for _, node := range nodes {
			count++
			if node.Children != nil && node.Expand {
				// Recursively count the children, if expanded
				countNodes(node.Children)
			}
		}
	}

	countNodes(m.nodes)

	return count

}

func (m Model) Width() int {
	return m.width
}

func (m Model) Height() int {
	return m.height
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model) SetWidth(newWidth int) {
	m.SetSize(newWidth, m.height)
}

func (m *Model) SetHeight(newHeight int) {
	m.SetSize(m.width, newHeight)
}

func (m Model) Cursor() int {
	return m.cursor
}

func (m *Model) SetCursor(cursor int) {
	m.cursor = cursor
}

func (m *Model) SetShowHelp() bool {
	return m.showHelp
}

func (m *Model) NavUp() {
	m.cursor--

	if m.cursor < 0 {
		m.cursor = 0
		return
	}

}

func (m *Model) NavDown() {
	m.cursor++

	if m.cursor >= m.NumberOfNodes() {
		m.cursor = m.NumberOfNodes() - 1
		return
	}
}

func (m *Model) InvertCollaped() {
	if m.currentNode.Children != nil {
		m.currentNode.Expand = !m.currentNode.Expand
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Up):
			m.NavUp()
		case key.Matches(msg, m.KeyMap.Down):
			m.NavDown()
		case key.Matches(msg, m.KeyMap.Collapse):
			m.InvertCollaped()
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}

	return m, nil
}

func (m *Model) View() string {
	availableHeight := m.height
	var sections []string

	nodes := m.Nodes()

	var help string
	if m.showHelp {
		help = m.helpView()
		availableHeight -= lipgloss.Height(help)
	}

	count := 0 // This is used to keep track of the index of the node we are on (important because we are using a recursive function)
	sections = append(sections, lipgloss.NewStyle().Height(availableHeight).Render(m.renderTree(m.nodes, 0, &count)), help)

	if len(nodes) == 0 {
		return "No data"
	}
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) getDisplayRange(maxRows int) (int, int) {
	rowsAbove := m.height / 2
	rowsBelow := m.height / 2
	if m.cursor < rowsAbove {
		rowsAbove = m.cursor
		rowsBelow = m.height - m.cursor
	}
	if m.cursor+rowsBelow > maxRows {
		rowsBelow = maxRows - m.cursor
		rowsAbove = m.height - rowsBelow
	}
	return m.cursor - rowsAbove, m.cursor + rowsBelow
}

func (m *Model) renderTree(remainingNodes []*Node, indent int, count *int) string {
	var b strings.Builder

	minRow, maxRow := m.getDisplayRange(m.NumberOfNodes())

	for _, node := range remainingNodes {

		var str string

		// If we aren't at the root, we add the arrow shape to the string
		if indent > 0 {
			shape := strings.Repeat(" ", (indent-1)*2) + m.Styles.Shapes.Render(bottomLeft) + " "
			str += shape
		}

		// Generate the correct index for the node
		idx := *count
		*count++

		// Format the string with fixed width for the value and description fields
		valueWidth := 10
		descWidth := 20
		valueStr := strings.ReplaceAll(fmt.Sprintf("%-*s", valueWidth, node.Value), "\n", " ")
		descStr := strings.ReplaceAll(fmt.Sprintf("%-*s", descWidth, node.Desc), "\n", " ")

		// If we are at the cursor, we add the selected style to the string
		if m.cursor == idx {
			m.currentNode = node
			str += fmt.Sprintf("%s\t\t%s\n", m.Styles.Selected.Render(valueStr), m.Styles.Selected.Render(descStr))
		} else if idx >= minRow && idx <= maxRow {
			str += fmt.Sprintf("%s\t\t%s\n", m.Styles.Unselected.Render(valueStr), m.Styles.Unselected.Render(descStr))
		} else {
			logrus.Debugf("Skipping node %d: %s", idx, node.Value)
		}

		b.WriteString(str)

		if node.Children != nil && node.Expand {
			childStr := m.renderTree(node.Children, indent+1, count)
			b.WriteString(childStr)
		}
	}

	return b.String()
}

func (m *Model) helpView() string {
	return m.Styles.Help.Render(m.Help.View(m))
}

func (m *Model) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.Up,
		m.KeyMap.Down,
		m.KeyMap.Collapse,
	}

	if m.AdditionalShortHelpKeys != nil {
		kb = append(kb, m.AdditionalShortHelpKeys()...)
	}

	return append(kb,
		m.KeyMap.Quit,
	)
}

func (m *Model) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.KeyMap.Up,
		m.KeyMap.Down,
		m.KeyMap.Collapse,
	}}

	return append(kb,
		[]key.Binding{
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		})
}
