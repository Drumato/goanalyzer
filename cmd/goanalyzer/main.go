package main

import (
	"github.com/Drumato/goanalyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	analyzers := []*analysis.Analyzer{
		goanalyzer.DependencyAnalyzer,
		goanalyzer.LazyToplevelAnalyzer,
		goanalyzer.LazyIdentifierScopeAnalyzer,
		goanalyzer.UndefinedUnitTestAnalyzer,
	}

	unitchecker.Main(analyzers...)
}
