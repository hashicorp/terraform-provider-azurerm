# Framework Adoption

## Resources

Example:
```go
package someservicepackage
// TODO

```

## Data Sources

Example:

```go
package someservicepackage
// TODO

```

## Provider Functions (Core >= 1.8)

See `internal/provider/function` for live/shipped examples.

## Ephemeral Resources (Core >= 1.10)

Example:

```go
package someazureservice

import (
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type MyEphemeralResource struct {
	sdk.EphemeralResourceMetadata
}

var _ sdk.EphemeralResource = MyEphemeralResource{}

func (m MyEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "azurerm_my_ephemeral_resource"
}

func (m MyEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ...
		},
		Blocks: map[string]schema.Block{
			// ...
		},
	}
}

func (m MyEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	m.Defaults(req, resp)
}

func (m MyEphemeralResource) Open(ctx context.Context, request ephemeral.OpenRequest, openResponse *ephemeral.OpenResponse) {
	client := m.Client.SomeAzureService.FooClient

}

```