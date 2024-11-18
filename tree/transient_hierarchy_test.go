package tree

import (
	"reflect"
	"testing"
)

func TestValidSchema(t *testing.T) {
	testCases := []struct {
		description string
		header      string
		want        map[string]int
	}{
		{
			"level 1", "level_1, item_id",
			map[string]int{"level_1": 0, "item_id": 1},
		},
		{
			"level 1 reversed order", "item_id, level_1",
			map[string]int{"level_1": 1, "item_id": 0},
		},
		{
			"level 2", "level_1, level_2, item_id",
			map[string]int{"level_1": 0, "level_2": 1, "item_id": 2},
		},
		{
			"level 2 reversed order", "item_id, level_1, level_2",
			map[string]int{"level_1": 1, "level_2": 2, "item_id": 0},
		},
		{
			"level 3", "level_1, level_2, level_3, item_id",
			map[string]int{"level_1": 0, "level_2": 1, "level_3": 2, "item_id": 3},
		},
		{
			"level 3 reversed order", "item_id, level_1, level_2, level_3",
			map[string]int{"level_1": 1, "level_2": 2, "level_3": 3, "item_id": 0},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			got, ok := generateSchemaColumns(testCase.header)
			if !ok {
				t.Errorf("header: %s, want ok: true, got ok: false", testCase.header)
			}

			if !reflect.DeepEqual(testCase.want, got) {
				t.Errorf("header: %s, want: %+v, got: %+v", testCase.header, testCase.want, got)
			}
		})
	}
}

func TestInvalidSchema(t *testing.T) {
	testCases := []struct {
		description string
		header      string
	}{
		{"empty header", ""},
		{"item_id only", "item_id"},
		{"missing item_id single level", "level_1"},
		{"missing level_1", "level_2,item_id"},
		{"missing item_id multiple levels", "level_1,level_2"},
		{"duplicated item_id", "item_id,level_1,item_id"},
		{"duplicated levels", "level_1,level_1,item_id"},
		{"invalid header name", "random,level_1,item_id"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {

			got, ok := generateSchemaColumns(testCase.header)
			if ok != false {
				t.Errorf("header: %s, want ok: false, got ok: true", testCase.header)
			}
			if got != nil {
				t.Errorf("header: %s, want: nil, got: %+v", testCase.header, got)
			}
		})
	}
}

func TestInvalidRetrieveValue(t *testing.T) {
	columns, _ := generateSchemaColumns("level_1,item_id")
	gotValue, gotOk := retrieveValue(columns, "random", []string{"category 1", "item 1"})
	if gotOk {
		t.Errorf("want ok: false, got ok: true")
	}
	if gotValue != "" {
		t.Errorf("want value: '', got value: %s", gotValue)
	}
}

func TestValidRetrieveValue(t *testing.T) {
	testCases := []struct {
		description string
		header      string
		key         string
		fields      []string
		wantValue   string
	}{
		{"level string", "level_1,item_id", "level_1", []string{"category 1", "item 1"}, "category 1"},
		{"level string with space", "level_1,item_id", "level_1", []string{" category 1 ", "item 1"}, "category 1"},
		{"item string", "level_1,item_id", "item_id", []string{"category 1", "item 1"}, "item 1"},
		{"item string with space", "level_1,item_id", "item_id", []string{"category 1", " item 1 "}, "item 1"},
		{"empty value", "level_1,item_id", "level_1", []string{"", "item 1"}, ""},
		{"empty value with space", "level_1,item_id", "level_1", []string{"  ", "item 1"}, ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			columns, _ := generateSchemaColumns("level_1,item_id")
			gotValue, gotOk := retrieveValue(columns, testCase.key, testCase.fields)
			if !gotOk {
				t.Errorf("want ok: true, got ok: false")
			}
			if testCase.wantValue != gotValue {
				t.Errorf("want value: %s, got value: %s", testCase.wantValue, gotValue)
			}
		})
	}
}

func TestValidExtractRecord(t *testing.T) {
	testCases := []struct {
		description string
		header      string
		data        string
		want        []string
	}{
		{"level 1", "level_1, item_id", "Category 1,item 1", []string{"Category 1", "item 1"}},
		{
			"level 1 + level 2",
			"level_1, level_2, item_id",
			"Category 1,Category 2,item 2",
			[]string{"Category 1", "Category 2", "item 2"},
		},
		{
			"level 1 + level 2 + level 3",
			"level_1, level_2, level_3,item_id",
			"Category 1,Category 2, Category 3, item 2",
			[]string{"Category 1", "Category 2", "Category 3", "item 2"},
		},
		{"level 1 reversed", "item_id,level_1", "item 1,Category 1", []string{"Category 1", "item 1"}},
		{
			"level 1 + level 2 reversed",
			"item_id,level_1, level_2",
			"item 3,Category 1,Category 2",
			[]string{"Category 1", "Category 2", "item 3"},
		},
		{
			"level 1 + level 2 + level 3 reversed",
			"item_id,level_1, level_2, level_3",
			"item 4,Category 1,Category 2, Category 3",
			[]string{"Category 1", "Category 2", "Category 3", "item 4"},
		},
		{"skip level 2", "level_1, level_2, item_id", "Category 1,,item 5", []string{"Category 1", "item 5"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			columns, _ := generateSchemaColumns(testCase.header)
			got, ok := extractRecord(columns, testCase.data)
			if !ok {
				t.Errorf("want ok: true, got ok: false")
			}
			if !reflect.DeepEqual(testCase.want, got) {
				t.Errorf("want record: %+v, got record: %+v", testCase.want, got)
			}
		})
	}
}

func TestInvalidExtractRecord(t *testing.T) {
	testCases := []struct {
		description string
		header      string
		data        string
	}{
		{"empty record", "level_1,item_id", ""},
		{"missing item", "level_1,item_id", "category 1"},
		{"missing level", "level_1,item_id", "item_1"},
		{"mismatch header and record", "level_1,level_2,item_id", "category 1,item_2"},
		{"empty level 1", "level_1,item_id", ",item_1"},
		{"empty item", "level_1,item_id", "category 1,"},
		{"skip level 1", "level_1,level_2,item_id", ",category 2,item_1"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			columns, _ := generateSchemaColumns(testCase.header)
			got, ok := extractRecord(columns, testCase.data)
			if ok {
				t.Errorf("data: %+v, want ok: false, got ok: true", testCase.data)
			}
			if nil != got {
				t.Errorf("data: %+v, want record: nil, got record: %+v", testCase.data, got)
			}
		})
	}
}
