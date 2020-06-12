# AZURERMR001

The AZURERMR001 analyzer reports for error message that uses "Error" prefix.

## Flagged Code

```go
fmt.Errorf("Error something failed")
errors.New("error something failed")
```

## Passing Code

```go
fmt.Errorf("something failed")
errors.New("something failed")
```

## Ignoring Reports

Singular reports can be ignored by adding the a `//lintignore:AZURERMR001` Go code comment at the end of the offending line or on the line immediately proceding, e.g.

```go
//lintignore:AZURERMR001
fmt.Errorf("Error something failed")
//lintignore:AZURERMR001
errors.New("error something failed")
```
