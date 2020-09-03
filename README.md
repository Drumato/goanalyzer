# go analyzer

Golangのソースコードの静的解析ツール．  
いくつかの機能が集まっている．

## Dependency analyzer(WIP)

パッケージ間の依存関係を解析し，  
DOT言語を生成してGraphvizで可視化する．  

## Scope checker(not implemented)

トップレベルで `var/const/type` 宣言されているが，同一パッケージ内で一つの手続内でしか用いられていないときなど，  
不要にスコープが広くなっているような識別子を検出する．

また，`testdata/src/lazy_ident_scope/a/a.go` の `x` など，  
ブロックスコープについても同様に検出する．

## How to use

```
# using go build
go build ./cmd/goanalyzer
mv goanalyzer.exe ../../bin/
go vet -vettool="$(which goanalyzer)" <go-package>

go test ./...
```
