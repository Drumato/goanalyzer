package goanalyzer_test

import (
	"testing"

	"goanalyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestDependencyAnalyzer is a test for DependencyAnalyzer.
func TestDependencyAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, goanalyzer.DependencyAnalyzer, "dependency/a")
}

