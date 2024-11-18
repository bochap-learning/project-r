package tree

import (
	"testing"
)

func BenchmarkConcurrentExtractLargeInput(b *testing.B) {
	for range b.N {
		testTarget, closer := constructTransientHierarchyFromFile("../testdata/large_input.csv")
		defer closer()
		testTarget.ConcurrentExtract()
	}
}

func BenchmarkSynchronousExtractLargeInput(b *testing.B) {
	for range b.N {
		testTarget, closer := constructTransientHierarchyFromFile("../testdata/large_input.csv")
		defer closer()
		testTarget.SynchronousExtract()
	}
}

func BenchmarkConcurrentExtractSmallInput(b *testing.B) {
	for range b.N {
		testTarget, closer := constructTransientHierarchyFromFile("../testdata/small_input.csv")
		defer closer()
		testTarget.ConcurrentExtract()
	}
}

func BenchmarkSynchronousExtractSmallInput(b *testing.B) {
	for range b.N {
		testTarget, closer := constructTransientHierarchyFromFile("../testdata/small_input.csv")
		defer closer()
		testTarget.SynchronousExtract()
	}
}
