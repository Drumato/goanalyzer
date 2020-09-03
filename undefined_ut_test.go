package goanalyzer_test

import (
	"github.com/Drumato/goanalyzer"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

// TestUndefinedUnitTestAnalyzer is a test for UndefinedUnitTestAnalyzer.
func TestUndefinedUnitTestAnalyzer(t *testing.T) {
	tests := []string{"undefined_unit_test/a"}
	testdata := analysistest.TestData()

	for _, tt := range tests{
		tt := tt
		t.Run(tt, func (t *testing.T) {
			t.Parallel()
			analysistest.Run(t, testdata, goanalyzer.UndefinedUnitTestAnalyzer, tt)
		})
	}
}
