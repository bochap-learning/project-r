package tree

import (
	"testing"
)

func TestValidNewTreeNode(t *testing.T) {
	testCases := []struct {
		description string
		record      [][]string
		want        string
	}{
		{
			"single record - level 1",
			[][]string{{"category 1", "item 1"}},
			"{\"children\":{\"category 1\":{\"children\":{\"item 1\":{\"item\":true}}}}}",
		},
		{
			"single record - level 1, level 2",
			[][]string{{"category 1", "category 2", "item 1"}},
			"{\"children\":{\"category 1\":{\"children\":{\"category 2\":{\"children\":{\"item 1\":{\"item\":true}}}}}}}",
		},
		{
			"single record - level 1, level 2, level 3",
			[][]string{{"category 1", "category 2", "category 3", "item 1"}},
			"{\"children\":{\"category 1\":{\"children\":{\"category 2\":{\"children\":{\"category 3\":{\"children\":{\"item 1\":{\"item\":true}}}}}}}}}",
		},
		{
			"double record - level 1",
			[][]string{{"category 1", "item 1"}, {"category 1", "item 2"}},
			"{\"children\":{\"category 1\":{\"children\":{\"item 1\":{\"item\":true},\"item 2\":{\"item\":true}}}}}",
		},
		{
			"double record - level 1, level 2",
			[][]string{{"category 1", "category 2", "item 1"}, {"category 1", "category 3", "item 1"}},
			"{\"children\":{\"category 1\":{\"children\":{\"category 2\":{\"children\":{\"item 1\":{\"item\":true}}},\"category 3\":{\"children\":{\"item 1\":{\"item\":true}}}}}}}",
		},
		{
			"double record - level 1, level 2, level 3",
			[][]string{{"category 1", "category 2", "category 3", "item 1"}, {"category 1", "category 3", "category 2", "item 5"}},
			"{\"children\":{\"category 1\":{\"children\":{\"category 2\":{\"children\":{\"category 3\":{\"children\":{\"item 1\":{\"item\":true}}}}},\"category 3\":{\"children\":{\"category 2\":{\"children\":{\"item 5\":{\"item\":true}}}}}}}}}",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			node, ok := NewTreeNode(testCase.record)
			if !ok {
				t.Errorf("want ok: true, got ok: false")
			}
			data, err := node.ExportJson()
			if err != nil {
				t.Errorf("want error: nil, got error: %s", err)
			}
			got := string(data)
			if testCase.want != got {
				t.Errorf("want node: %s, got node: %s", testCase.want, got)
			}
		})
	}
}

func TestInvalidNewTreeNode(t *testing.T) {
	testCases := []struct {
		description string
		records     [][]string
	}{
		{"nil level", nil},
		{"empty records", [][]string{}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			got, ok := NewTreeNode(testCase.records)
			if ok {
				t.Errorf("want ok: false, got ok: true")
			}
			if got != nil {
				t.Errorf("want nil tree node, got tree node: %+v", got)
			}
		})
	}
}
