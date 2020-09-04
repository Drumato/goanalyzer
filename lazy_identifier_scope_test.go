package goanalyzer_test

import (
	"github.com/Drumato/goanalyzer"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestLazyIdentifierScopeAnalyzer is a test for LazyIdentifierScopeAnalyzer.
func TestLazyIdentifierScopeAnalyzer(t *testing.T) {
	tests := []string{
		"lazy_ident_scope/a",
		"lazy_ident_scope/nested",
		"lazy_ident_scope/closure",
		"lazy_ident_scope/in_if",
		"lazy_ident_scope/in_for"}
	testdata := analysistest.TestData()

	for _, tt := range tests{
		tt := tt
		t.Run(tt, func (t *testing.T) {
			t.Parallel()
			analysistest.Run(t, testdata, goanalyzer.LazyIdentifierScopeAnalyzer, tt)
		})
	}
}
