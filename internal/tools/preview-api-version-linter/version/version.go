package version

type Version struct {
	Module  string // example: github.com/hashicorp/go-azure-sdk/resource-manager
	Service string // example: compute
	Version string // example: 2020-02-02-preview
}
