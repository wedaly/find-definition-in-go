# find-definition-in-go

Example code for finding a Go definition.

## Usage

Build the command-line tool:
```
go build .
```

Lookup a definition:
```
./find-definition-in-go main.go 32 8
```

The output should look like this:
```
"lookupAndPrintGoDef" is defined at /Users/will/go/src/github.com/wedaly/find-definition-in-go/main.go:39:6
```
