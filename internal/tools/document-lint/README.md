# Introduction 
This tool is to fix the documetn mismatch and format issues of terraform azurerm provider.

What Can It Do:
1. document format fix.
2. check required/optional value of properties
3. check default value of properties
4. check force new value
5. check timeout of create/update/read/delete
6. check properties missed/redundant in document
7. check possible values of proeprties in document

# Getting Started
```bash
# print the usage
go run main.go -h

# check documents and print the error information
go run main.go check

# check and try to fix existing errors
go run main.go fix
```