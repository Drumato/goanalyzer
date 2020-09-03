package goanalyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const lazyToplevelDoc = "lazytoplevelanalyzer is ..."

var (
	// LazyToplevelAnalyzer is ...
	LazyToplevelAnalyzer = &analysis.Analyzer{
		Name: "lazytoplevelanalyzer",
		Doc:  lazyToplevelDoc,
		Run:  detectLazyToplevelDecls,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	idRefCounter map[string]*topLevelDecl
)

type topLevelDecl struct {
	// 参照回数．identifier を参照する手続きの数を数える
	refCount int
	p        token.Pos
	// どのスコープから参照されたか
	refBy *types.Scope
}

func detectLazyToplevelDecls(pass *analysis.Pass) (interface{}, error) {
	idRefCounter = make(map[string]*topLevelDecl)
	pkgScope := pass.Pkg.Scope()

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	inspect.Preorder([]ast.Node{new(ast.Ident)}, func(n ast.Node) {
		//識別子の使用位置がどのスコープに含まれているか
		if id, ok := n.(*ast.Ident); ok {
			if id.IsExported() {
				return
			}

			// 定義位置の取得
			def := pkgScope.Lookup(id.Name)
			if def != nil {
				use, used := pass.TypesInfo.Uses[id]
				_, isFn := def.Type().(*types.Signature)
				if isFn {
					return
				}
				if use == nil || !used {
					// 定義情報の格納
					idRefCounter[id.Name] = &topLevelDecl{p: def.Pos()}
					return
				}

				// 参照カウントの更新
				useScope := pass.Pkg.Scope().Innermost(id.Pos())
				if entry, exist := idRefCounter[id.Name]; exist && entry.refBy != useScope {
					idRefCounter[id.Name].refCount++
					idRefCounter[id.Name].refBy = useScope
				}
			}
		}
	})

	// 検出
	for _, id := range idRefCounter {
		if id.refCount == 1 {
			pass.Reportf(id.p, "This identifier is only referenced by a function so should move the declaration to it")
		}
	}
	return nil, nil
}