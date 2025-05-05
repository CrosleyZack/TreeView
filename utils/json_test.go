package utils

import (
	"reflect"
	"strings"
	"testing"

	tree "github.com/savannahostrowski/tree-bubble"
)

func TestGetTypedEntry(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want TypedEntry
	}{
		{
			name: "string",
			v:    "hello",
			want: TypedEntry{Type: entryTypeString, Value: "hello"},
		},
		{
			name: "int",
			v:    42,
			want: TypedEntry{Type: entryTypeInt, Value: 42},
		},
		{
			name: "float",
			v:    3.14,
			want: TypedEntry{Type: entryTypeFloat, Value: 3.14},
		},
		{
			name: "boolean",
			v:    true,
			want: TypedEntry{Type: entryTypeBoolean, Value: true},
		},
		{
			name: "array",
			v:    []interface{}{"a", 1, "c"},
			want: TypedEntry{Type: entryTypeArray, Value: []interface{}{"a", 1, "c"}},
		},
		{
			name: "map",
			v:    map[string]interface{}{"key": "value", "number": 123, "nested": map[string]interface{}{"nestedKey": "nestedValue"}},
			want: TypedEntry{Type: entryTypeMap, Value: map[string]interface{}{"key": "value", "number": 123, "nested": map[string]interface{}{"nestedKey": "nestedValue"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTypedEntry(tt.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTypedEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypedEntryString(t *testing.T) {
	tests := []struct {
		name string
		e    TypedEntry
		want string
	}{
		{
			name: "string",
			e:    TypedEntry{Type: entryTypeString, Value: "hello"},
			want: "hello",
		},
		{
			name: "int",
			e:    TypedEntry{Type: entryTypeInt, Value: 42},
			want: "42",
		},
		{
			name: "float",
			e:    TypedEntry{Type: entryTypeFloat, Value: 3.14},
			want: "3.140",
		},
		{
			name: "boolean",
			e:    TypedEntry{Type: entryTypeBoolean, Value: true},
			want: "true",
		},
		{
			name: "array",
			e:    TypedEntry{Type: entryTypeArray, Value: []interface{}{"a", 1, "c"}},
			want: "a 1 c",
		},
		{
			name: "map",
			e:    TypedEntry{Type: entryTypeMap, Value: map[string]interface{}{"key": "value", "number": 123, "nested": map[string]interface{}{"nestedKey": "nestedValue"}}},
			// TODO flakey
			want: "value 123 nestedValue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("typedEntry.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypedEntryTreeify(t *testing.T) {
	tests := []struct {
		name string
		e    TypedEntry
		want tree.Node
	}{
		{
			name: "string",
			e:    TypedEntry{Type: entryTypeString, Value: "hello"},
			want: tree.Node{
				Value:    "",
				Desc:     "hello",
				Children: []tree.Node{},
			},
		},
		{
			name: "int",
			e:    TypedEntry{Type: entryTypeInt, Value: 42},
			want: tree.Node{
				Value:    "",
				Desc:     "42",
				Children: []tree.Node{},
			},
		},
		{
			name: "float",
			e:    TypedEntry{Type: entryTypeFloat, Value: 3.14},
			want: tree.Node{
				Value:    "",
				Desc:     "3.140",
				Children: []tree.Node{},
			},
		},
		{
			name: "boolean",
			e:    TypedEntry{Type: entryTypeBoolean, Value: true},
			want: tree.Node{
				Value:    "",
				Desc:     "true",
				Children: []tree.Node{},
			},
		},
		// TODO Debug
		// {
		// 	name: "array",
		// 	e:    typedEntry{Type: entryTypeArray, Value: []any{"a", 1, "c"}},
		// 	want: tree.Node{
		// 		Value: "",
		// 		Desc:  "a 1 c",
		// 		Children: []tree.Node{
		// 			{
		// 				Value: "0",
		// 				Desc:  "a",
		// 			},
		// 			{
		// 				Value: "1",
		// 				Desc:  "1",
		// 			},
		// 			{
		// 				Value: "2",
		// 				Desc:  "c",
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "map",
		// 	e:    typedEntry{Type: entryTypeMap, Value: map[string]any{"key": "value", "number": 123, "nested": map[string]any{"nestedKey": "nestedValue"}}},
		// 	want: tree.Node{
		// 		Desc: "value 123 nestedValue",
		// 		Children: []tree.Node{
		// 			{
		// 				Value: "key",
		// 				Desc:  "value",
		// 			},
		// 			{
		// 				Value: "number",
		// 				Desc:  "123",
		// 			},
		// 			{
		// 				Value: "nested",
		// 				Desc:  "",
		// 				Children: []tree.Node{
		// 					{
		// 						Value: "nestedKey",
		// 						Desc:  "nestedValue",
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Treeify(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("typedEntry.Treeify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDictGet(t *testing.T) {
	tests := []struct {
		name string
		d    JsonBlob
		k    string
		want TypedEntry
	}{
		{
			name: "bad key",
			d:    JsonBlob{"key": "value"},
			k:    "foo",
			want: TypedEntry{Type: entryTypeUnknown, Value: nil},
		},
		{
			name: "existing key with string value",
			d:    JsonBlob{"key": "value"},
			k:    "key",
			want: TypedEntry{Type: entryTypeString, Value: "value"},
		},
		{
			name: "existing key with int value",
			d:    JsonBlob{"number": 42},
			k:    "number",
			want: TypedEntry{Type: entryTypeInt, Value: 42},
		},
		{
			name: "existing key with float value",
			d:    JsonBlob{"pi": 3.14},
			k:    "pi",
			want: TypedEntry{Type: entryTypeFloat, Value: 3.14},
		},
		{
			name: "existing key with boolean value",
			d:    JsonBlob{"flag": true},
			k:    "flag",
			want: TypedEntry{Type: entryTypeBoolean, Value: true},
		},
		{
			name: "existing key with array value",
			d:    JsonBlob{"list": []interface{}{"a", 1, "c"}},
			k:    "list",
			want: TypedEntry{Type: entryTypeArray, Value: []interface{}{"a", 1, "c"}},
		},
		{
			name: "existing key with map value",
			d:    JsonBlob{"nested": map[string]interface{}{"key": "value"}},
			k:    "nested",
			want: TypedEntry{Type: entryTypeMap, Value: map[string]interface{}{"key": "value"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Get(tt.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dict.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDictString(t *testing.T) {
	tests := []struct {
		name string
		d    JsonBlob
		want string
	}{
		{
			name: "empty dict",
			d:    JsonBlob{},
			want: "",
		},
		{
			name: "single key-value pair",
			d:    JsonBlob{"key": "value"},
			want: "key: value",
		},
		{
			name: "multiple key-value pairs",
			d:    JsonBlob{"key1": "value1", "key2": 42, "key3": true},
			// TODO flakey
			want: "key1: value1 key2: 42 key3: true",
		},
		{
			name: "nested map",
			d:    JsonBlob{"nested": map[string]interface{}{"key": "value"}, "another": map[string]interface{}{"nested": "foo"}},
			want: "nested: value another: foo",
		},
		{
			name: "array value",
			d:    JsonBlob{"list": []interface{}{"a", 1, "c"}},
			want: "list: a 1 c",
		},
		{
			name: "exceeding max string length",
			d:    JsonBlob{"key1": strings.Repeat("a", MaxStringLength), "key2": "value2"},
			want: "key1: " + strings.Repeat("a", MaxStringLength),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("dict.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
