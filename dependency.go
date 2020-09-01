package goanalyzer

import (
	"fmt"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"path"
	"time"
)

const dependencyAnalyzerDocument = "dependencyanalyzer is ..."

var (

    standardPackages = make(map[string]bool)
	// DependencyAnalyzer is ...
	DependencyAnalyzer = &analysis.Analyzer{
		Name: "dependencyanalyzer",
		Doc:  dependencyAnalyzerDocument,
		Run:  analyzeDependency,
	}
)

func analyzeDependency(pass *analysis.Pass) (interface{}, error) {
	graph      := make(map[string]map[string]bool)
	stdPkgs, err := packages.Load(nil, "std")

	if err != nil {
		return nil, err
	}

	for _, stdPkg := range stdPkgs {
		standardPackages[stdPkg.PkgPath] = true
	}

	for _, f := range pass.Files {
		initial, err := packages.Load(&packages.Config{
			Mode:  packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax,
			Dir:   path.Dir(f.Name.Name),
			Tests: true,
		})
		if err != nil {
			return nil, err
		}

		for _, pkg := range initial {
			recursiveVisitPkgImports(graph, pkg)
		}

		if err := renderDOTProgram(graph); err != nil {
			return nil, err
		}

		pass.Reportf(f.Pos(), "analyze finished")
	}

	return nil, nil
}

func recursiveVisitPkgImports(graph map[string]map[string]bool, pkg *packages.Package) {
	pkgName := pkg.PkgPath
	if isStd := standardPackages[pkgName]; isStd {
		return
	}

	for _, importedPkg := range pkg.Imports {
		importPkgName := importedPkg.PkgPath
		if isStd := standardPackages[importPkgName]; isStd {
			continue
		}

		if _, ok := graph[pkgName]; !ok {
			graph[pkgName] = make(map[string]bool)
		}
		graph[pkgName][importPkgName] = true

		recursiveVisitPkgImports(graph, importedPkg)
	}
	return
}

func renderDOTProgram(graph map[string]map[string]bool) error {
	program := "digraph G { \n"

	for src, dsts := range graph {
		program += fmt.Sprintf("    \"%s\"[];\n", src)

		for dst := range dsts {
			program += fmt.Sprintf("    \"%s\" -> \"%s\";\n", src, dst)
		}
	}

	program += "} \n"

	if err := ioutil.WriteFile(fmt.Sprintf("%s.dot", time.Now().Format("2006_01_02")), []byte(program), 0644); err != nil {
		return err
	}

	return nil
}
