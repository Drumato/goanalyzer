package goanalyzer_test

import (
	"goanalyzer"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestLazyToplevelAnalyzer is a test for LazyToplevelAnalyzer.
func TestLazyToplevelAnalyzer(t *testing.T) {
	tests := []string{"lazy_toplevel_decl/a", "lazy_toplevel_decl/fp", "lazy_toplevel_decl/multiref", "lazy_toplevel_decl/samepkg"}
	testdata := analysistest.TestData()

	for _, tt := range tests{
		tt := tt
		t.Run(tt, func (t *testing.T) {
			t.Parallel()
			analysistest.Run(t, testdata, goanalyzer.LazyToplevelAnalyzer, tt)
		})
	}
}
