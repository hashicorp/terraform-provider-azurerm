# Guide: New Service Package

There's a few steps involved in adding a new Service Package.

1. Create a new directory within `./internal/services` with the Service Name (see [naming](reference-naming.md)).

2. Create an empty Client within the Service Package (`./internal/services/{name}/client/client.go`):

```go
package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
}

func NewClient(o *common.ClientOptions) *Client {
	return &Client{}
}
```

3. Create an empty Registration within the Service Package (`./internal/services/{name}/registration.go`) which implements the `TypedServiceRegistration` interface:

```go
package {name}

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration = Registration{}
)

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "App Service"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"App Service",
	}
}
```

4. [Register the Service Registration](https://github.com/hashicorp/terraform-provider-azurerm/blob/2ff15cca48adc7315f67d8b653409e621963ca64/internal/provider/services.go#L109-L118).
5. Define and Register the Client for this Service Package.
   * [Add this to the Client struct](https://github.com/hashicorp/terraform-provider-azurerm/blob/2ff15cca48adc7315f67d8b653409e621963ca64/internal/clients/client.go#L118-L120).
   * [Call the Register function](https://github.com/hashicorp/terraform-provider-azurerm/blob/2ff15cca48adc7315f67d8b653409e621963ca64/internal/clients/client.go#L221-L234).
6. Re-run the generation to ensure the generated files are up to date (`make generate`).

At this point the Service Package should be registered, and you can [build a new Data Source](guide-new-data-source.md) or [a new Resource](guide-new-resource.md) as required.
