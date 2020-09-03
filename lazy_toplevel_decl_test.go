package goanalyzer_test

import (
	"goanalyzer"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestLazyScopeAnalyzer is a test for LazyScopeAnalyzer.
func TestLazyScopeAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, goanalyzer.LazyToplevelAnalyzer, "lazy_toplevel_decl/a")
	analysistest.Run(t, testdata, goanalyzer.LazyToplevelAnalyzer, "lazy_toplevel_decl/fp")
	analysistest.Run(t, testdata, goanalyzer.LazyToplevelAnalyzer, "lazy_toplevel_decl/multiref")
	analysistest.Run(t, testdata, goanalyzer.LazyToplevelAnalyzer, "lazy_toplevel_decl/samepkg")
}
