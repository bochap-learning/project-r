package tree

import (
	"bufio"
	"os"
	"testing"
)

func constructEnvironment(iFile string) *TransientHierarchy {
	input, _ := os.Open(iFile)
	inputReader := bufio.NewReader(input)
	testTarget, _ := NewTransientHierarchy(inputReader)
	return testTarget
}

func BenchmarkConcurrentExtractLargeInput(b *testing.B) {
	for range b.N {
		testTarget := constructEnvironment("../testdata/large_input.csv")
		testTarget.ConcurrentExtract()
	}
}

func BenchmarkSynchronousExtractLargeInput(b *testing.B) {
	for range b.N {
		testTarget := constructEnvironment("../testdata/large_input.csv")
		testTarget.SynchronousExtract()
	}
}

func BenchmarkConcurrentExtractSmallInput(b *testing.B) {
	for range b.N {
		testTarget := constructEnvironment("../testdata/small_input.csv")
		testTarget.ConcurrentExtract()
	}
}

func BenchmarkSynchronousExtractSmallInput(b *testing.B) {
	for range b.N {
		testTarget := constructEnvironment("../testdata/small_input.csv")
		testTarget.SynchronousExtract()
	}
}
