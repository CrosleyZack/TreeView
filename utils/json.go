package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/crosleyzack/bubbles/tree"
)

const (
	MaxStringLength = 100
)

type EntryType int

const (
	entryTypeUnknown EntryType = iota
	entryTypeString
	entryTypeInt
	entryTypeFloat
	entryTypeBoolean
	entryTypeArray
	entryTypeMap
)

type JsonBlob map[string]any

func (d JsonBlob) Get(k string) TypedEntry {
	var entry TypedEntry
	if v, ok := d[k]; ok {
		entry = getTypedEntry(v)
	}
	return entry
}

func (d JsonBlob) String() string {
	ret := strings.Builder{}
	first := true
	for k := range d {
		if ret.Len() > MaxStringLength {
			break
		}
		ret.WriteString(spacerToken(first))
		ret.WriteString(fmt.Sprintf("%s: %s", k, d.Get(k).String()))
		first = false
	}
	return ret.String()
}

func (d JsonBlob) Treeify() *tree.Model {
	nodes := make([]*tree.Node, 0)
	for k, v := range d {
		node := getTypedEntry(v).Treeify()
		node.Value = k
		node.Expand = true
		nodes = append(nodes, node)
	}
	return tree.New(nodes, 1, 1)
}

type TypedEntry struct {
	Type  EntryType `json:"type"`
	Value any       `json:"value"`
	Key   string    `json:"key"`
}

func getTypedEntry(v any) TypedEntry {
	switch v := v.(type) {
	case string:
		return TypedEntry{Type: entryTypeString, Value: v}
	case float64:
		return TypedEntry{Type: entryTypeFloat, Value: v}
	case int:
		return TypedEntry{Type: entryTypeInt, Value: v}
	case bool:
		return TypedEntry{Type: entryTypeBoolean, Value: v}
	case []any:
		return TypedEntry{Type: entryTypeArray, Value: v}
	case map[string]any:
		return TypedEntry{Type: entryTypeMap, Value: v}
	default:
		return TypedEntry{}
	}
}

func (e TypedEntry) String() string {
	ret := strings.Builder{}
	first := true
	stack := NewQueue[TypedEntry]()
	for entry := e; stack != nil; stack, entry = stack.Pop() {
		if ret.Len() > MaxStringLength {
			break
		}
		switch entry.Type {
		case entryTypeString:
			ret.WriteString(spacerToken(first))
			ret.WriteString(fmt.Sprintf("%s", entry.Value))
			first = false
		case entryTypeBoolean:
			ret.WriteString(spacerToken(first))
			ret.WriteString(fmt.Sprintf("%t", entry.Value))
			first = false
		case entryTypeInt:
			ret.WriteString(spacerToken(first))
			ret.WriteString(fmt.Sprintf("%d", entry.Value))
			first = false
		case entryTypeFloat:
			ret.WriteString(spacerToken(first))
			ret.WriteString(fmt.Sprintf("%.3f", entry.Value))
			first = false
		case entryTypeArray:
			for _, item := range entry.Value.([]interface{}) {
				stack = stack.Push(getTypedEntry(item))
			}
		case entryTypeMap:
			for _, v := range entry.Value.(map[string]interface{}) {
				stack = stack.Push(getTypedEntry(v))
			}
		default:
			return ""
		}
	}
	return ret.String()
}

func (e TypedEntry) Treeify() *tree.Node {
	node := tree.Node{
		Desc:     e.String(),
		Expand:   false,
		Children: make([]*tree.Node, 0),
	}
	switch e.Type {
	case entryTypeArray:
		for i, item := range e.Value.([]interface{}) {
			child := getTypedEntry(item).Treeify()
			child.Value = strconv.FormatUint(uint64(i), 10)
			node.Children = append(node.Children, child)
		}
	case entryTypeMap:
		for k, v := range e.Value.(map[string]interface{}) {
			child := getTypedEntry(v).Treeify()
			child.Value = k
			node.Children = append(node.Children, child)
		}
	}
	return &node
}

// spacerToken returns a space if not the first element
func spacerToken(first bool) string {
	if first {
		return ""
	}
	return " "
}
