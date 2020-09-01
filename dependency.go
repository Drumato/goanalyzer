package goanalyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const doc = "dependencyanalyzer is ..."

var (
// DependencyAnalyzer is ...
DependencyAnalyzer = &analysis.Analyzer{
	Name: "dependencyanalyzer",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{},
}
graph = make(map[string][]*ast.File)

)

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		err := recursiveVisit(f)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(len(graph))
	return nil, nil
}

func recursiveVisit(file *ast.File) error {

	for _, spec := range file.Imports {
		packageName, err  := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return err
		}
		fileSet := token.NewFileSet()
		importedFile, err := findSourceFile(fileSet, packageName)
		if err != nil {
			// ディレクトリだとして再度チェック
			pkgs, err := findPackageDir(fileSet, packageName)
			if err != nil {
				return err
			}

			for _, pkg := range pkgs{
				for _, f := range pkg.Files {
					recursiveVisit(f)
				}
			}

			return nil
		}

		// ソースファイルを探索する
		if _, ok := graph[file.Name.Name] ; !ok {
			graph[file.Name.Name] = make([]*ast.File, 0)
		}

		graph[file.Name.Name] = append(graph[file.Name.Name], importedFile)
	}

	return nil
}

func findSourceFile(set *token.FileSet, pkgName string) (*ast.File, error) {
	importedPath, err := concatWithGOPATH(pkgName)
	if err != nil {
		// GOROOT(標準ライブラリ)からも探索する
		importedPath, err = concatWithGOROOT(pkgName)
		if err != nil {
			return nil, err
		}
	}

	importedFile, err := parser.ParseFile(set, importedPath, nil, 0)
	if err != nil {
		return nil, err
	}
	return importedFile, nil
}

func findPackageDir(set *token.FileSet, pkgName string) (map[string]*ast.Package, error) {
	importedPath, err := concatWithGOPATH(pkgName)
	if err != nil {
		// GOROOT(標準ライブラリ)からも探索する
		importedPath, err = concatWithGOROOT(pkgName)
		if err != nil {
			return nil, err
		}
	}

	packages, err := parser.ParseDir(set, importedPath, nil, 0)
	if err != nil {
		return nil, err
	}

	return packages, nil
}

func concatWithGOPATH(base string) (string, error) {
	key, ok := os.LookupEnv("GOPATH")
	if !ok {
		return "", fmt.Errorf("GOPATH not found")
	}

	return path.Join(key, "src", base), nil
}

func concatWithGOROOT(base string) (string, error) {
	out, err := exec.Command("which", "go").Output()
	if err != nil {
		return "", err
	}

	// hoge/bin/goのようになっている
	outputStr := path.Dir(path.Dir(string(strings.TrimSpace(string(out)))))
	return path.Join(outputStr, "src", base), nil
}