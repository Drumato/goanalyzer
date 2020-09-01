package main

import (
"goanalyzer"
"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(goanalyzer.DependencyAnalyzer) }


