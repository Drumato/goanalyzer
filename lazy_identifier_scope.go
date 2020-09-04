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
)

type defInfo struct {
	base *ast.Ident
	// 定義されたスコープ
	scope *types.Scope
	// 参照されたスコープ群
	refs []*types.Scope
}

func detectLazyIdentifierScope(pass *analysis.Pass) (interface{}, error) {
	identifierDefs := make(map[types.Object]*defInfo)
	pkgScope := pass.Pkg.Scope()
	// 識別子定義について探索
	visitAllDecls(pass, identifierDefs, pkgScope)

	// 参照しているスコープをすべて調べる
	countReferenceScopes(pass, pkgScope, identifierDefs)

	// 検出
	for _, id := range identifierDefs {
		onlyRefByChildScope := true
		// 参照するスコープがすべて定義スコープの内側かチェック
		for _, child := range id.refs{
			onlyRefByChildScope = onlyRefByChildScope && !child.Contains(id.base.Pos())
		}
		if onlyRefByChildScope {
			pass.Reportf(id.base.Pos(), "This identifier is only referenced in a scope so should move the declaration to it")
		}
	}

	return nil, nil
}

func countReferenceScopes(pass *analysis.Pass, pkgScope *types.Scope, defs map[types.Object]*defInfo) {
	for useId, use := range pass.TypesInfo.Uses{
		_, ok := defs[use]
		if !ok {
			continue
		}
		usedScope := pkgScope.Innermost(useId.Pos())
		defs[use].refs = append(defs[use].refs, usedScope)
	}
}

func visitAllDecls(pass *analysis.Pass, identifiers map[types.Object]*defInfo, pkgScope *types.Scope) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	declFilters := []ast.Node{
		(*ast.DeclStmt)(nil),
		(*ast.AssignStmt)(nil),
	}

	inspect.Preorder(declFilters, func(n ast.Node) {
		switch st := n.(type){
		case *ast.DeclStmt:
			visitDeclStmt(pass, st, pkgScope, identifiers)
		case *ast.AssignStmt:
			visitAssignStmt(pass, st, pkgScope, identifiers)
		default:
			return
		}
	})

}

func visitAssignStmt(pass *analysis.Pass, st *ast.AssignStmt, pkgScope *types.Scope, identifiers map[types.Object]*defInfo) {
	for _, l := range st.Lhs{
		id, ok := l.(*ast.Ident)
		if !ok{
			return
		}

		def, defined := pass.TypesInfo.Defs[id]
		if def == nil || !defined {
			return
		}
		if _, ok := identifiers[def]; !ok && pkgScope != def.Parent() {
			identifiers[def] = &defInfo{base: id, scope: def.Parent(), refs: make([]*types.Scope, 0)}
		}
	}
}

func visitDeclStmt(pass *analysis.Pass, st *ast.DeclStmt, pkgScope *types.Scope, identifiers map[types.Object]*defInfo) {
	switch d := st.Decl.(type) {
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			switch sp := spec.(type){
			case *ast.ValueSpec:
				visitValueSpec(pass, sp, pkgScope, identifiers)
			}
		}
	case *ast.FuncDecl:
	}
}

func visitValueSpec(pass *analysis.Pass, sp *ast.ValueSpec, pkgScope *types.Scope, identifiers map[types.Object]*defInfo) {
	for _, id := range sp.Names{
		def, defined := pass.TypesInfo.Defs[id]
		if def == nil || !defined {
			return
		}
		if _, ok := identifiers[def]; !ok && pkgScope != def.Parent(){
			identifiers[def] = &defInfo{base: id, scope: def.Parent(), refs: make([]*types.Scope, 0)}
		}
	}
}