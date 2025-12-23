package sdk

import "regexp"

type SdkType struct {
	// The go module name of the SDK, this also correspond to the path under the vendor folder
	// example: github.com/hashicorp/go-azure-sdk/resource-manager
	Module string

	// The regex that should be used to extract service and version value while directories are recursed. It needs to have
	// exactly two capture groups, first for the service and second for the version.
	//
	// For example: `resource-manager/([^/]+)/([^/]+)` will capture the service "containerregistry" and the version
	// "2023-11-01-preview" for the input string:
	// vendor/github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-11-01-preview/tokens
	ServiceAndVersionRegex *regexp.Regexp
}

var (
	// Example: vendor/github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-11-01-preview/tokens
	GO_AZURE_SDK = SdkType{
		Module:                 "github.com/hashicorp/go-azure-sdk/resource-manager",
		ServiceAndVersionRegex: regexp.MustCompile(`(?i)resource-manager/([^/]+)/([^/]+-preview)`),
	}

	// Example: vendor/github.com/jackofallops/kermit/sdk/synapse/2020-08-01-preview
	KERMIT_SDK = SdkType{
		Module:                 "github.com/jackofallops/kermit/sdk",
		ServiceAndVersionRegex: regexp.MustCompile(`(?i)sdk/([^/]+)/([^/]+-preview)`),
	}

	// Example:	vendor/github.com/jackofallops/giovanni/storage/2023-11-03-preview/file/shares
	GIOVANNI_SDK = SdkType{
		Module:                 "github.com/jackofallops/giovanni",
		ServiceAndVersionRegex: regexp.MustCompile(`(?i)giovanni/([^/]+)/([^/]+-preview)`),
	}

	// Example: vendor/github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview
	AZURE_SDK_FOR_GO_TRACK_1 = SdkType{
		Module:                 "github.com/Azure/azure-sdk-for-go",
		ServiceAndVersionRegex: regexp.MustCompile(`(?i)services/preview/([^/]+)/mgmt/([^/]+)`),
	}
)

var SdkTypes = []SdkType{
	GO_AZURE_SDK,
	KERMIT_SDK,
	GIOVANNI_SDK,
	AZURE_SDK_FOR_GO_TRACK_1,
}
