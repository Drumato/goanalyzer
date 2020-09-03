package goanalyzer

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const lazyIdentifierScopeDoc = "lazyidentifierscopeanalyzer is ..."

var (
	// LazyIdentifierScopeAnalyzer is ...
	LazyIdentifierScopeAnalyzer = &analysis.Analyzer{
		Name:     "lazyidentifierscopeanalyzer",
		Doc:      lazyIdentifierScopeDoc,
		Run:      detectLazyIdentifierScope,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	defsIdentifiers = make(map[string]*localVar)
)

type localVar struct {
	base types.Object
	// 定義されたスコープ
	scope *types.Scope
	// 参照されたスコープ群
	refs []*types.Scope
}

func detectLazyIdentifierScope(pass *analysis.Pass) (interface{}, error) {
	// トップレベル以下のスコープに対し，再帰関数を適用する．
	// これにより，各スコープで定義された識別子情報を構築できる．
	pkgScope := pass.Pkg.Scope()

	for fnScopeIdx := 0; fnScopeIdx < pkgScope.NumChildren(); fnScopeIdx++ {
		fnScope := pkgScope.Child(fnScopeIdx)

		recursiveVisitScope(pass, fnScope)
	}

	// 識別子の使用場所を探索．
	// 参照スコープの情報を埋める
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	inspect.Preorder([]ast.Node{new(ast.Ident)}, func(n ast.Node) {

		// TODO: ここで識別子の使用位置がどのスコープに含まれているか取りたい，が取れない．
		if id, ok := n.(*ast.Ident); ok {
			_, exist := defsIdentifiers[id.Name]
			if exist {
				// fmt.Println(entry.base.Parent())
			}
		}
	})

	// TODO: 参照スコープが定義箇所より内側かつ一つのスコープからの参照であれば検出とする

	return nil, nil
}

func recursiveVisitScope(pass *analysis.Pass, outer *types.Scope){
	// 外側のスコープで定義された識別子を取得する
	for _, idName := range outer.Names(){
		obj := outer.Lookup(idName)
		defsIdentifiers[idName] = &localVar{base: obj, scope: outer, refs: make([]*types.Scope, 0)}
	}

	for innerScopeIdx := 0; innerScopeIdx < outer.NumChildren(); innerScopeIdx++ {
		innerScope := outer.Child(innerScopeIdx)
		recursiveVisitScope(pass, innerScope)
	}
}