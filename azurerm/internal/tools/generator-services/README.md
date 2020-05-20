## Generator: Services

Each Service Definition contains metadata (such as the Display Name & Website Categories) required in other parts of the codebase, for easier grouping.

This generator takes that metadata and uses it to generate two things:

1. Website Categories - which validates the categories used in the website exist, required for website deployments to happen.
2. Service Definitions - generates the list of services used to run the Acceptance Tests

This is run via go:generate whenever the "SupportedServices" array is changed so that this is kept up-to-date.

## Example Usage

```
go run main.go -path=../../path/to/root-directory
```

## Arguments

* `help` - Show help?

* `path` - The Relative Path to the root of the repository
