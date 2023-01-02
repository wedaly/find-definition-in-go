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

You should see this output:
```
"lookupAndPrintGoDef" is defined at main.go:39:6
```
