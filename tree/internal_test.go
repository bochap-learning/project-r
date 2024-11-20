package tree

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"sort"
)

const wellFormedPayload = `item_id,level_1,level_2
1,A,C
2,A,C
5,B,E
6,B,
`

const illegalSkipLevelPayload = `item_id,level_1,level_2
1,A,C
2,A,C
3,,D
4,A,D
`

func sortBySliceOfSliceOfStrings(data [][]string) {
	sort.Slice(data, func(i, j int) bool {
		// Concatenate all strings in each inner slice
		str1 := ""
		for _, s := range data[i] {
			str1 += s
		}
		str2 := ""
		for _, s := range data[j] {
			str2 += s
		}

		// Compare the concatenated strings
		return str1 < str2
	})
}

// openFile takes a file path and returns an io.ReadCloser that
// automatically closes the underlying file when reading is finished or
// an error occurs.
func openFile(filePath string) (io.ReadCloser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &autoClosingReader{file}, nil
}

// autoClosingReader is a custom io.ReadCloser that wraps an io.Reader
// and ensures the underlying Reader is closed when Close() is called.
type autoClosingReader struct {
	io.Reader
}

// Close closes the underlying Reader.
func (r *autoClosingReader) Close() error {
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func constructTransientHierarchyFromFile(file string) (TransientHierarchy, func() error) {
	input, err := openFile(file)
	if err != nil {
		panic(err)
	}
	inputReader := bufio.NewReader(input)
	testTarget, _ := NewTransientHierarchy(inputReader)
	return testTarget, func() error {
		return input.Close()
	}
}

func constructTreeNodeFromFile(file string) (TreeNode, func() error) {
	var result TreeNode
	input, err := openFile(file)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(input)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result, func() error {
		return input.Close()
	}
}
