package tree

import "encoding/json"

// TreeNode represents a node in the hierarchical tree structure.
// It contains a boolean indicating whether it's an item node and a map of child nodes, keyed by their names.
// The Item field is true if the node represents an item, and false otherwise.
// The Children field is a map of child nodes, where the keys are the names of the child nodes and the values are pointers to the child nodes.
// If a node has no children, the Children field will be nil.
type TreeNode struct {
	Item     bool                 `json:"item,omitempty"`
	Children map[string]*TreeNode `json:"children,omitempty"`
}

// NewTreeNode creates a new TreeNode from a slice of slices of strings.  Each inner slice represents a record,
// and each string in the inner slice represents a level in the hierarchy.  The last string in each inner slice
// represents the item ID.  The function returns a pointer to the root node of the tree and a boolean indicating
// whether the creation was successful.  It returns false if the input is invalid.  A valid input must contain
// at least one record, and each record must contain at least one level and an item ID.
func NewTreeNode(records [][]string) (*TreeNode, bool) {
	if len(records) == 0 {
		return nil, false
	}
	root := &TreeNode{Children: make(map[string]*TreeNode)}

	for _, record := range records {
		current := root
		for _, field := range record {
			if current.Children == nil {
				current.Children = make(map[string]*TreeNode)
			}
			if _, ok := current.Children[field]; !ok {
				current.Children[field] = &TreeNode{}
			}
			current = current.Children[field]
		}
		current.Item = true
	}
	return root, true
}

func (t *TreeNode) ExportJson() ([]byte, error) {
	return json.Marshal(t)
}
