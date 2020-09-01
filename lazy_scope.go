package goanalyzer

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

const lazyScopeDocument = "lazyscopeanalyzer is ..."

var (
	// LazyScopeAnalyzer is ...
	LazyScopeAnalyzer = &analysis.Analyzer{
		Name: "lazyscopeanalyzer",
		Doc:  lazyScopeDocument,
		Run:  detectLazyScopes,
	}
	idRefCounter map[string]*identifier
)

type identifier struct {
	*ast.Ident
	// 参照回数．identifier を参照する手続きの数を数える
	refCount int
	// どの関数から利用されているか
	refBy map[string]bool
}

func detectLazyScopes(pass *analysis.Pass) (interface{}, error) {
	// まずはトップレベル宣言のみチェックする
	for _, f := range pass.Files {
		idRefCounter = make(map[string]*identifier)
		fnDecls := correctTopLevelDecls(f.Decls)

		for fn := range fnDecls{
			for _, stmt := range fn.Body.List{
				switch st := stmt.(type) {
				case *ast.ExprStmt:
					checkInExpr(fn.Name.Name, st.X)
				case *ast.DeclStmt:
					switch decl := st.Decl.(type) {
					case *ast.GenDecl:
						checkDeclStmt(fn.Name.Name, decl)
					default:
					}
				default:
					continue
				}
			}
		}

		for _, id := range idRefCounter {
			if id.refCount == 1 && !id.IsExported() {
				pass.Reportf(id.NamePos, "This identifier is declared in an unnecessarily wide scope")
			}
		}
	}
	return nil, nil
}

func checkDeclStmt(fnName string, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		v, ok := spec.(*ast.ValueSpec)
		if ok {
			checkInExpr(fnName, v.Type)
		}
	}
}

func checkInExpr(fnName string, ex ast.Expr) {
	switch e := ex.(type) {
	case *ast.CallExpr:
		for _, arg := range e.Args {
			checkInExpr(fnName, arg)
		}
	case *ast.BinaryExpr:
		checkInExpr(fnName, e.X)
		checkInExpr(fnName, e.Y)
	case *ast.Ident:
		_, ok := idRefCounter[e.Name]
		if ok {
			if _, ok := idRefCounter[e.Name].refBy[fnName]; !ok {
				idRefCounter[e.Name].refCount++
			}
			idRefCounter[e.Name].refBy[fnName] = true
		}
	default:
	}
}

func correctTopLevelDecls(decls []ast.Decl) (map[*ast.FuncDecl]int) {
	fnDecls := make(map[*ast.FuncDecl]int)

	for _, decl := range decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			analyzeGenDecl(d)
		case *ast.FuncDecl:
			fnDecls[d] = 0
		default:
			continue
		}
	}

	return fnDecls
}

func analyzeGenDecl(decl *ast.GenDecl){
	switch decl.Tok{
	case token.CONST, token.VAR:
		for _, spec := range decl.Specs{
			v := spec.(*ast.ValueSpec)

			for identIdx := range v.Names{
				idRefCounter[v.Names[identIdx].Name] = &identifier{v.Names[identIdx], 0, make(map[string]bool)}
			}
		}
	case token.TYPE:
		for _, spec := range decl.Specs{
			t := spec.(*ast.TypeSpec)
			idRefCounter[t.Name.Name] = &identifier{t.Name, 0, make(map[string]bool)}
		}
	default:
	}
}