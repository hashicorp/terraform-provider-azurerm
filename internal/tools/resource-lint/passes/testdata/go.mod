module testdata

go 1.25.3

require (
	github.com/hashicorp/go-azure-helpers/lang/pointer v0.0.0
	github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines v0.0.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.0
)

replace (
	github.com/hashicorp/go-azure-helpers/lang/pointer => ./src/github.com/hashicorp/go-azure-helpers/lang/pointer
	github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines => ./src/github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines
	github.com/hashicorp/terraform-plugin-sdk/v2 => ./src/github.com/hashicorp/terraform-plugin-sdk/v2
)
