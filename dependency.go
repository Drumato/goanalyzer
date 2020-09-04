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
	// DependencyAnalyzer is ...
	DependencyAnalyzer = &analysis.Analyzer{
		Name: "dependencyanalyzer",
		Doc:  dependencyAnalyzerDocument,
		Run:  analyzeDependency,
	}
)

func prepareStandardPackages() (map[string]bool, error) {
	stdPkgs, err := packages.Load(nil, "std")
	if err != nil {
		return nil, err
	}
	standardPackages := make(map[string]bool)
	for _, stdPkg := range stdPkgs {
		standardPackages[stdPkg.PkgPath] = true
	}
	return standardPackages, nil
}

func analyzeDependency(pass *analysis.Pass) (interface{}, error) {
	// 依存グラフに標準ライブラリを含めるとすごいことになってしまうので，含めない
	standardPackages, err := prepareStandardPackages()
	if err != nil {
		return nil, err
	}

	graph := make(map[string]map[string]bool)

	for _, f := range pass.Files {
		initial, err := packages.Load(&packages.Config{
			Mode:  packages.NeedName | packages.NeedImports,
			Dir:   path.Dir(f.Name.Name),
			Tests: true,
		})
		if err != nil {
			return nil, err
		}

		for _, pkg := range initial {
			recursiveVisitPkgImports(graph, standardPackages, pkg)
		}

		if err := renderDOTProgram(graph); err != nil {
			return nil, err
		}

		pass.Reportf(f.Pos(), "dependency analyze finished")
	}

	return nil, nil
}

func recursiveVisitPkgImports(graph map[string]map[string]bool, standardPackages map[string]bool, pkg *packages.Package) {
	pkgName := pkg.PkgPath

	// 依存パッケージをすべて探索
	for _, importedPkg := range pkg.Imports {
		importPkgName := importedPkg.PkgPath

		// 標準パッケージならば無視する
		if isStd := standardPackages[importPkgName]; isStd {
			continue
		}

		// 枝刈り
		// 既に探索したパッケージならば登録されているので，探索しない
		if _, exist := graph[importPkgName]; exist{
			continue
		}

		if _, ok := graph[pkgName]; !ok {
			graph[pkgName] = make(map[string]bool)
		}
		graph[pkgName][importPkgName] = true

		recursiveVisitPkgImports(graph, standardPackages, importedPkg)
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
