package tree

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

// generateSchemaColumns generates a map of column names to their indices in the header.
// It takes a header string as input, which is expected to be a comma-separated list of column names.
// The function returns a map where keys are column names and values are their corresponding indices.
// It returns an error if the header is invalid.  A valid header must contain "level_1" and "item_id" columns,
// and all columns must be valid.  The function also checks for duplicate column names.
func generateSchemaColumns(header string) (map[string]int, bool) {
	if len(header) == 0 {
		return nil, false
	}
	a := strings.Split(header, ",")
	if len(a) < 2 {
		return nil, false
	}
	result := make(map[string]int)
	for i, h := range a {
		k := strings.TrimSpace(h)
		if _, ok := allowedColumns[k]; !ok {
			return nil, false
		}
		result[k] = i
	}
	if _, ok := result[level1]; !ok {
		return nil, false
	}
	if _, ok := result[itemId]; !ok {
		return nil, false
	}

	if len(result) != len(a) {
		return nil, false
	}
	return result, true
}

// generateLevels generates a function that yields level labels.
// The generated function takes a yield function as an argument.  The yield function is called for each level,
// passing the level index and label.  The yield function should return true to continue generating levels,
// and false to stop.  The level index starts at 0.  The level label is formatted as "level_{index+1}".
// For example, generateLevels(3) would generate a function that yields:
// (0, "level_1"), (1, "level_2"), (2, "level_3")
func generateLevels(end int) func(func(int, string) bool) {
	return func(yield func(int, string) bool) {
		for i := 0; i < end; i++ {
			if !yield(i, fmt.Sprintf("level_%d", i+1)) {
				return
			}
		}
	}
}

// retrieveValue retrieves a value from a record based on its key and column indices.
// It takes a map of column names to indices, a key, and a slice of fields as input.
// It returns the value associated with the key and a boolean indicating whether the key was found.
func retrieveValue(columns map[string]int, key string, fields []string) (string, bool) {
	index, ok := columns[key]
	if !ok {
		return "", false
	}
	return strings.TrimSpace(fields[index]), true
}

// extractRecord extracts a record from a data string based on the given schema columns.
// It takes a map of column names to indices and a data string as input.
// The data string is expected to be a comma-separated list of fields.
// The function returns a slice of strings representing the extracted record and a boolean indicating success.
// It returns false in the boolean if the data is invalid.  A valid record must contain "item_id" column,
// and all columns must be valid.  The function also checks for schema and data column count parity.
func extractRecord(columns map[string]int, data string) ([]string, bool) {
	fields := strings.Split(data, ",")
	if len(fields) <= 1 {
		return nil, false
	}

	if len(columns) != len(fields) {
		return nil, false
	}

	item, ok := retrieveValue(columns, itemId, fields)
	if !ok {
		return nil, false
	}

	if len(item) == 0 {
		return nil, false
	}

	record := make([]string, 0, len(columns))
	previous := "level_0"
	// subtract 1 from len(s.columns) to remove item column
	for _, levelLabel := range generateLevels(len(columns) - 1) {
		level, ok := retrieveValue(columns, levelLabel, fields)
		if !ok {
			return nil, false
		}
		if len(level) > 0 && len(previous) == 0 {
			return nil, false
		}
		previous = level
		if len(level) == 0 {
			continue
		}
		record = append(record, level)
	}

	if len(record) == 0 {
		return nil, false
	}
	record = append(record, item)
	return record, true
}

// TransientHierarchy is a struct that facilitates the extraction of data from a CSV reader, supporting both synchronous and concurrent operations.
// The CSV data must adhere to a specific format: a header row with column names is mandatory, and this header must define "level_1" and "item_id" columns.
// Any columns beyond these two are interpreted as valid hierarchy level columns. This structure enables the construction of a hierarchical representation from the CSV data,
// where each row corresponds to an item and its associated levels within the hierarchy.
type TransientHierarchy struct {
	reader  *bufio.Reader
	columns map[string]int
}

// NewTransientHierarchy creates a new TransientHierarchy from a bufio.Reader.  It reads the header from the reader,
// validates the header, and returns a new TransientHierarchy object.  It returns an error if the header is invalid.
// A valid header must contain "level_1" and "item_id" columns, and all columns must be valid.
func NewTransientHierarchy(reader *bufio.Reader) (*TransientHierarchy, error) {
	header, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	columns, ok := generateSchemaColumns(header)
	if !ok {
		return nil, err
	}
	return &TransientHierarchy{reader: reader, columns: columns}, nil
}

// SynchronousExtract extracts all records from the CSV reader synchronously.  It reads each line from the reader,
// extracts the record using extractRecord, and appends it to the content slice.  It returns the content slice
// and an error if any errors occur during the process.
func (t *TransientHierarchy) SynchronousExtract() ([][]string, error) {
	var content [][]string
	for {
		line, err := t.reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		record, ok := extractRecord(t.columns, line)
		if !ok {
			return nil, fmt.Errorf("invalid data: %s", line)
		}
		content = append(content, record)
	}

	return content, nil
}

// ConcurrentExtract extracts all records from the CSV reader concurrently.  It reads each line from the reader,
// launches a goroutine to extract the record using extractRecord, and appends it to the content slice.
// It uses a WaitGroup to wait for all goroutines to complete and a mutex to protect the content slice from race conditions.
// It returns the content slice and an error if any errors occur during the process.
// Errors are sent through an error channel to prevent blocking.
func (t *TransientHierarchy) ConcurrentExtract() ([][]string, error) {
	var (
		content [][]string
		wg      sync.WaitGroup
		mu      sync.Mutex            // Mutex for safe access to content
		errChan = make(chan error, 1) // Channel to signal errors
	)

	for {
		line, err := t.reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			record, ok := extractRecord(t.columns, line)
			if !ok {
				select {
				case errChan <- fmt.Errorf("invalid data: %s", line): // Send error
				default: // Prevent blocking if another error was already sent
				}
				return
			}
			mu.Lock()
			content = append(content, record)
			mu.Unlock()
		}(line)
	}

	wg.Wait()
	close(errChan)

	// Check for errors from goroutines
	if err := <-errChan; err != nil {
		return nil, err
	}

	return content, nil
}
