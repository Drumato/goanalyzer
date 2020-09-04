# go analyzer

Golangのソースコードの静的解析ツール．  
いくつかの機能が集まっている．

## Dependency analyzer

パッケージ間の依存関係を解析し，  
DOT言語を生成してGraphvizで可視化する．  

## Scope checker

トップレベルで `var/const/type` 宣言されているが，同一パッケージ内で一つの手続内でしか用いられていないときなど，  
不要にスコープが広くなっているような識別子を検出する．

また，`testdata/src/lazy_ident_scope/a/a.go` の `x` など，  
ブロックスコープについても同様に検出する．

## undefined unit test

ユニットテストが定義されていない関数を検出する．   

## How to use

### as a binary

```
# using go build
go build -o goanalyzer ./cmd/goanalyzer
mv goanalyzer "$GOPATH/bin"
go vet -vettool="$(which goanalyzer)" <go-package>

# test
go test ./...
```
