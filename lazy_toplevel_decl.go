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
)

type topLevelDecl struct {
	// 参照回数．identifier を参照する手続きの数を数える
	refCount int
	p        token.Pos
	// どのスコープから参照されたか
	refBy *types.Scope
}

func detectLazyToplevelDecls(pass *analysis.Pass) (interface{}, error) {
	idRefCounter := make(map[types.Object]*topLevelDecl)
	pkgScope := pass.Pkg.Scope()

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	inspect.Preorder([]ast.Node{new(ast.Ident)}, func(n ast.Node) {
		//識別子の使用位置がどのスコープに含まれているか
		if id, ok := n.(*ast.Ident); ok {
			if id.IsExported() {
				return
			}

			use, used := pass.TypesInfo.Uses[id]
			if use == nil || !used {
				// 定義情報の登録
				appendDefinedInformation(pkgScope, idRefCounter, id)
			}

			// 参照カウントの更新
			countReference(pkgScope, idRefCounter,use, id)
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

func appendDefinedInformation(pkgScope *types.Scope, idRefCounter map[types.Object]*topLevelDecl, id *ast.Ident) {
	def := pkgScope.Lookup(id.Name)
	if def == nil{
		return
	}

	_, isFn := def.Type().(*types.Signature)
	if isFn {
		return
	}

	// 定義情報の格納
	idRefCounter[def] = &topLevelDecl{p: def.Pos()}
}

func countReference(pkgScope *types.Scope, idRefCounter map[types.Object]*topLevelDecl, use types.Object, id *ast.Ident) {
	usedScope := pkgScope.Innermost(id.Pos())
	if entry, exist := idRefCounter[use]; exist && entry.refBy != usedScope {
		idRefCounter[use].refCount++
		idRefCounter[use].refBy = usedScope
	}
}