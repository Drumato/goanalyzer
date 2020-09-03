package goanalyzer

import (
	"bytes"
	"fmt"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
)

const undefinedUnitTestDoc = "undefinedunittestanalyzer is ..."


var (
	// UndefinedUnitTestAnalyzer is ...
	UndefinedUnitTestAnalyzer = &analysis.Analyzer{
		Name:     "undefinedunittestanalyzer",
		Doc:      undefinedUnitTestDoc,
		Run:      findUndefinedUnitTest,
	}
)

func findUndefinedUnitTest(pass *analysis.Pass) (interface{}, error) {
	pkgScope := pass.Pkg.Scope()

	for _, pkgScopeDeclName := range pkgScope.Names(){
		pkgScopeDecl := pkgScope.Lookup(pkgScopeDeclName)
		if pkgScopeDeclName == "main" {
			continue
		}

		if _, isFn := pkgScopeDecl.Type().(*types.Signature); !isFn{
			continue

		}

		// テスト関数の場合も同様に無視する
		if strings.HasPrefix(pkgScopeDeclName, "Test"){
			continue
		}

		testFnName := fmt.Sprintf("Test%s", toCapitalize(pkgScopeDeclName))
		if testFn := pkgScopeDecl.Parent().Lookup(testFnName); testFn == nil {
			pass.Reportf(pkgScopeDecl.Pos(), "This function's unit test is not defined")
		}
	}

	return nil, nil
}

func toCapitalize(s string) string {
	return fmt.Sprintf("%s%s", bytes.ToUpper([]byte{s[0]}), s[1:])
}