package main

import (
	"goanalyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	analyzers := []*analysis.Analyzer{
		goanalyzer.DependencyAnalyzer,
		goanalyzer.LazyToplevelAnalyzer,
	}
	unitchecker.Main(analyzers...)
}
