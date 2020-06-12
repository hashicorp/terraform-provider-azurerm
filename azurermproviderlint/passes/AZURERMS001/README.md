# AZURERMS001

The AZURERMS001 analyzer reports Schema contains case-insensitive validation but missing case diff suppression.

## Flagged Code

```go
schema.Schema{
    ValidateFunc: validation.StringInSlice([]string{}, true),
}

// or

schema.Schema{
    ValidateFunc: validation.StringInSlice([]string{}, true),
    DiffSuppressFunc: nil,
}
```

## Passing Code

```go
schema.Schema{
    ValidateFunc:     validation.StringInSlice([]string{}, true),
    DiffSuppressFunc: suppress.CaseDifference,
}
```

## Ignoring Reports

Singular reports can be ignored by adding the a `//lintignore:AZURERMR001` Go code comment at the end of the offending line or on the line immediately proceding, e.g.

```go
//lintignore:AZURERMS001
schema.Schema{
    ValidateFunc: validation.StringInSlice([]string{}, true),
}
```
