package goanalyzer_test

import (
	"goanalyzer"
	"testing"

"golang.org/x/tools/go/analysis/analysistest"
)

// TestLazyScopeAnalyzer is a test for LazyScopeAnalyzer.
func TestLazyScopeAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, goanalyzer.LazyScopeAnalyzer, "lazyscope/a")
	analysistest.Run(t, testdata, goanalyzer.LazyScopeAnalyzer, "lazyscope/fp")
	analysistest.Run(t, testdata, goanalyzer.LazyScopeAnalyzer, "lazyscope/multiref")
}


