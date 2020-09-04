package goanalyzer

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

const undefinedUnitTestDoc = "undefinedunittestanalyzer is ..."


var (
	// UndefinedUnitTestAnalyzer is ...
	UndefinedUnitTestAnalyzer = &analysis.Analyzer{
		Name:     "undefinedunittestanalyzer",
		Doc:      undefinedUnitTestDoc,
		Run:      findUndefinedUnitTest,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
)

func findUndefinedUnitTest(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	testDefs := make(map[string]bool)

	// テスト関数を収集
	for id, v := range pass.TypesInfo.Defs{
		if v == nil {
			continue
		}

		if _, isFn := v.Type().(*types.Signature); !isFn{
			continue
		}

		if !id.IsExported(){
			continue
		}

		if !strings.HasPrefix(id.Name, "Test"){
			continue
		}

		testDefs[id.Name] = true
	}

	// 非テスト関数に対し，対応するテスト関数が存在するかチェック
	inspect.Preorder([]ast.Node{new(ast.Ident)}, func(n ast.Node){
		id, ok := n.(*ast.Ident)
		if !ok {
			return
		}

		def, defined := pass.TypesInfo.Defs[id]
		if def == nil || !defined {
			return
		}

		if _, isFn := def.Type().(*types.Signature); !isFn{
			return
		}

		if strings.HasPrefix(id.Name, "Test"){
			return
		}

		if id.IsExported(){
			return
		}

		if id.Name == "main" || id.Name == "init" {
			return
		}

		testFnName := fmt.Sprintf("Test%s", id.Name)
		if _, ok := testDefs[testFnName]; !ok {
			pass.Reportf(id.Pos(), "This function's unit test is not defined")
		}
	})

	return nil, nil
}