package goanalyzer_test

import (
	"testing"

	"goanalyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, goanalyzer.DependencyAnalyzer, "a")
}

