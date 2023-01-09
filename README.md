# find-definition-in-go

Example code for finding a Go definition. For a full walkthrough, please see my blog post [implement "find definition" in 77 lines of go](https://dev-nonsense.com/posts/find-definition-in-go/).

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
