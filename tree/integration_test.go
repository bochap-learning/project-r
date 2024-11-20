package tree

import (
	"reflect"
	"testing"
)

func TestConcurrentIntegration(t *testing.T) {
	testCases := []struct {
		description string
		source      string
		target      string
	}{
		{"testing concurrent extraction of small input", "../testdata/small_input.csv", "../testdata/small_output.json"},
		{"testing concurrent extraction of large input", "../testdata/large_input.csv", "../testdata/large_output.json"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			source, sourceCloser := constructTransientHierarchyFromFile(testCase.source)
			defer sourceCloser()
			want, wantCloser := constructTreeNodeFromFile(testCase.target)
			defer wantCloser()
			records, err := source.ConcurrentExtract()
			if err != nil {
				t.Errorf("want error: nil, got error: %s", err)
			}
			got, ok := NewTreeNode(records)
			if !ok {
				t.Errorf("want ok: true, got ok: false")
			}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("want content: %+v, got content: %+v", want, got)
			}
		})
	}
}

func TestSynchronousIntegration(t *testing.T) {
	testCases := []struct {
		description string
		source      string
		target      string
	}{
		{"testing concurrent extraction of small input", "../testdata/small_input.csv", "../testdata/small_output.json"},
		{"testing concurrent extraction of large input", "../testdata/large_input.csv", "../testdata/large_output.json"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			source, sourceCloser := constructTransientHierarchyFromFile(testCase.source)
			defer sourceCloser()
			want, wantCloser := constructTreeNodeFromFile(testCase.target)
			defer wantCloser()
			records, err := source.SynchronousExtract()
			if err != nil {
				t.Errorf("want error: nil, got error: %s", err)
			}
			got, ok := NewTreeNode(records)
			if !ok {
				t.Errorf("want ok: true, got ok: false")
			}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("want content: %+v, got content: %+v", want, got)
			}
		})
	}
}
