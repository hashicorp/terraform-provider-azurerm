module github.com/terraform-providers/terraform-provider-azurerm

require (
	github.com/Azure/azure-sdk-for-go v32.5.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.0
	github.com/Azure/go-autorest/autorest/adal v0.6.0
	github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/btubbs/datetime v0.1.0
	github.com/davecgh/go-spew v1.1.1
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-azure-helpers v0.5.0
	github.com/hashicorp/go-getter v1.3.1-0.20190627223108-da0323b9545e
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/go-version v1.1.0
	github.com/hashicorp/terraform v0.12.6
	github.com/relvacode/iso8601 v0.0.0-20181221151331-e9cae14c704e // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/satori/uuid v0.0.0-20160927100844-b061729afc07
	github.com/sirupsen/logrus v1.2.0 // indirect
	github.com/terraform-providers/terraform-provider-azuread v0.4.1-0.20190610202312-5a179146b9f9
	github.com/tombuildsstuff/giovanni v0.3.2
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net v0.0.0-20190502183928-7f726cade0ab
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/tombuildsstuff/giovanni => github.com/tcz001/giovanni v0.5.0
