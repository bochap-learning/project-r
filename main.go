package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bochap-learning/r-project/tree"
)

func synchronousExtract(readCloser io.ReadCloser) (func() ([][]string, error), error) {
	reader := bufio.NewReader(readCloser)
	transientTree, err := tree.NewTransientHierarchy(reader)
	if err != nil {
		return nil, err
	}
	return transientTree.SynchronousExtract, nil
}

func concurrentExtract(readCloser io.ReadCloser) (func() ([][]string, error), error) {
	reader := bufio.NewReader(readCloser)
	transientTree, err := tree.NewTransientHierarchy(reader)
	if err != nil {
		return nil, err
	}
	return transientTree.ConcurrentExtract, nil
}

func handlerFactory(extractFactory func(readCloser io.ReadCloser) (func() ([][]string, error), error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "text/csv" {
			http.Error(w, "Invalid content type", http.StatusUnsupportedMediaType)
			return
		}

		extract, err := extractFactory(r.Body)
		if err != nil {
			http.Error(w, "Invalid csv header", http.StatusBadRequest)
			return
		}

		content, err := extract()
		if err != nil {
			http.Error(w, "Invalid csv content", http.StatusBadRequest)
			return
		}

		tree, ok := tree.NewTreeNode(content)
		if !ok {
			http.Error(w, "Invalid csv content", http.StatusBadRequest)
			return
		}
		data, err := tree.ExportJson()
		if err != nil {
			http.Error(w, "Invalid csv content", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func main() {
	http.HandleFunc("/", handlerFactory(synchronousExtract))
	http.HandleFunc("/concurrent", handlerFactory(concurrentExtract))
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
