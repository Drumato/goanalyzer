package goanalyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

const lazyToplevelDoc = "lazytoplevelanalyzer is ..."

var (
	// LazyToplevelAnalyzer is ...
	LazyToplevelAnalyzer = &analysis.Analyzer{
		Name:     "lazytoplevelanalyzer",
		Doc:      lazyToplevelDoc,
		Run:      detectLazyToplevelDecls,
	}
	idRefCounter map[string]*identifier
)

type identifier struct {
	// 参照回数．identifier を参照する手続きの数を数える
	refCount int
	p        token.Pos
	// どの関数から参照されたか
	refByFn string
}

func detectLazyToplevelDecls(pass *analysis.Pass) (interface{}, error) {
	// トップレベル宣言をチェックし，辞書を構築する
	idRefCounter = make(map[string]*identifier)
	pkgScope := pass.Pkg.Scope()
	correctTopLevelDeclarations(pkgScope)

	// 参照をカウント
	// inspect.Preorderを使うと定義箇所まで調べてしまうので，ast.Inspectを用いる
	for _, f := range pass.Files{
		for _, decl := range f.Decls{
			if fn, ok := decl.(*ast.FuncDecl); ok {
				ast.Inspect(fn, func(n ast.Node) bool {
					if id, ok := n.(*ast.Ident); ok {
						if entry, exist := idRefCounter[id.Name]; exist && entry.refByFn != fn.Name.Name {
							idRefCounter[id.Name].refCount++
							idRefCounter[id.Name].refByFn = fn.Name.Name
						}
						return false
					}
					return true
				})
			}
		}
	}

	// 検出
	for _, id := range idRefCounter {
		if id.refCount == 1 {
			pass.Reportf(id.p, "This identifier is only referenced by a function so should move the declaration to it")
		}
	}
	return nil, nil
}

func correctTopLevelDeclarations(pkgScope *types.Scope) {
	for _, pkgVarName := range pkgScope.Names() {
		pkgVar := pkgScope.Lookup(pkgVarName)
		if !pkgVar.Exported() {
			if _, ok := pkgVar.Type().(*types.Signature); !ok {
				idRefCounter[pkgVar.Name()] = &identifier{p: pkgVar.Pos()}
			}
		}
	}
}