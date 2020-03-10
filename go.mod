module github.com/terraform-providers/terraform-provider-azurerm

require (
	contrib.go.opencensus.io/exporter/ocagent v0.5.0 // indirect
	github.com/Azure/azure-sdk-for-go v38.1.0+incompatible
	github.com/Azure/go-autorest v13.0.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.9.3
	github.com/Azure/go-autorest/autorest/azure/auth v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/btubbs/datetime v0.1.0
	github.com/davecgh/go-spew v1.1.1
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-azure-helpers v0.10.0
	github.com/hashicorp/go-getter v1.4.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/terraform-plugin-sdk v1.6.0
	github.com/pkg/errors v0.8.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/satori/uuid v0.0.0-20160927100844-b061729afc07
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
	github.com/spf13/cobra v0.0.5 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/terraform-providers/terraform-provider-azuread v0.6.1-0.20191007035844-361c0a206ad4
	github.com/tombuildsstuff/giovanni v0.9.0
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	golang.org/x/net v0.0.0-20191009170851-d66e71096ffb
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/Azure/go-autorest/autorest => github.com/tombuildsstuff/go-autorest/autorest v0.9.3-hashi-auth

go 1.13
