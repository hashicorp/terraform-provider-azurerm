## Website Categories

This tool generates files necessary for the website from the service registrations which already exist in code.

This is run via go:generate whenever the "SupportedServices" array is changed so that this is kept up-to-date.

## Example Usage

```
go run main.go -path=/path/to/output/file
```

## Arguments

* `help` - Show help?

* `path` - The Relative Path to the `allowed-subcategories` file used for the website

