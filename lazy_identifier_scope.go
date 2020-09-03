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
	identifierDefs = make(map[string]*defInfo)
)

type defInfo struct {
	base types.Object
	// 定義されたスコープ
	scope *types.Scope
	// 参照されたスコープ群
	refs []*types.Scope
}

func detectLazyIdentifierScope(pass *analysis.Pass) (interface{}, error) {
	// 識別子について探索
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	inspect.Preorder([]ast.Node{new(ast.Ident)}, func(n ast.Node) {
		//識別子の使用位置がどのスコープに含まれているか
		if id, ok := n.(*ast.Ident); ok {
			// 定義位置の取得
			def, defined := pass.TypesInfo.Defs[id]
			if def != nil && defined {
				_, isFn := def.Type().(*types.Signature)
				if isFn {
					return
				}

				identifierDefs[id.Name] = &defInfo{base: def, scope: def.Parent(), refs: make([]*types.Scope, 0)}
				return
			}

			_, defined = identifierDefs[id.Name]
			if defined {
				identifierDefs[id.Name].refs = append(identifierDefs[id.Name].refs, pass.Pkg.Scope().Innermost(id.Pos()))
			}
		}
	})

	// 参照スコープが定義箇所より内側かつ一つのスコープからの参照であれば検出とする
	for _, id := range identifierDefs {
		if len(id.refs) == 1 && id.refs[0] != id.scope {
			pass.Reportf(id.base.Pos(), "This identifier is only referenced in a scope so should move the declaration to it")
		}
	}

	return nil, nil
}
