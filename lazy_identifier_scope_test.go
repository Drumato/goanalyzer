package goanalyzer_test

import (
"goanalyzer"
"testing"

"golang.org/x/tools/go/analysis/analysistest"
)

// TestLazyIdentifierScopeAnalyzer is a test for LazyIdentifierScopeAnalyzer.
func TestLazyIdentifierScopeAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, goanalyzer.LazyIdentifierScopeAnalyzer, "lazy_ident_scope/a")
}
