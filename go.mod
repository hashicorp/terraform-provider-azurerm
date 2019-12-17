module github.com/terraform-providers/terraform-provider-azurerm

require (
	github.com/Azure/azure-sdk-for-go v36.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.3
	github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/btubbs/datetime v0.1.0
	github.com/davecgh/go-spew v1.1.1
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-azure-helpers v0.10.0
	github.com/hashicorp/go-getter v1.4.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/terraform v0.12.9
	github.com/hashicorp/terraform-plugin-sdk v1.1.1
	github.com/satori/go.uuid v1.2.0
	github.com/satori/uuid v0.0.0-20160927100844-b061729afc07
	github.com/terraform-providers/terraform-provider-azuread v0.6.1-0.20191007035844-361c0a206ad4
	github.com/tombuildsstuff/giovanni v0.7.1
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/Azure/azure-sdk-for-go => github.com/hashicorp/azure-sdk-for-go v36.2.2-hashi+incompatible

go 1.13
