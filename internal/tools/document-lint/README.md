# Introduction 
This tool detects and fixes inconsistencies in the AzureRM Terraform Provider resource documentation.

## The following can be checked/fixed:
1. Formatting of documentation.
2. The Required/Optional value of properties.
3. The Default value of properties.
4. The ForceNew value of properties.
5. The TimeOut value of create/update/read/delete functions.
6. Properties that are present in the schema but missing in the documentation and vice versa.
7. The list of PossibleValues.

# Getting Started
```bash
# print the usage
go run main.go -h

# check documents and print the error information
go run main.go check

# check and try to fix existing errors
go run main.go fix
```