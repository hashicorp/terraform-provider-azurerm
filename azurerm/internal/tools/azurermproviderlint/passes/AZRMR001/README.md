# AZRMR001

The AZRMR001 analyzer reports for error message that uses "Error" prefix.

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

Singular reports can be ignored by adding the a `//lintignore:AZRMR001` Go code comment at the end of the offending line or on the line immediately proceding, e.g.

```go
//lintignore:AZRMR001
fmt.Errorf("Error something failed")
//lintignore:AZRMR001
errors.New("error something failed")
```
