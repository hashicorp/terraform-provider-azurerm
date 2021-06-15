package sdk

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

// combineSchema combines the arguments (user-configurable) and attributes (read-only) schema fields
// into a canonical object - ensuring that each contains the relevant information
//
// whilst this may look overkill, this allows for simpler implementations in other tooling, for example
// when generating documentation
func combineSchema(arguments map[string]*schema.Schema, attributes map[string]*schema.Schema) (*map[string]*schema.Schema, error) {
	out := make(map[string]*schema.Schema)

	for k, v := range arguments {
		if _, alreadyExists := out[k]; alreadyExists {
			return nil, fmt.Errorf("%q already exists in the schema", k)
		}

		if v.Computed && !(v.Optional || v.Required) {
			return nil, fmt.Errorf("%q is a Computed-only field - this should be specified as an Attribute", k)
		}

		out[k] = v
	}

	for k, v := range attributes {
		if _, alreadyExists := out[k]; alreadyExists {
			return nil, fmt.Errorf("%q already exists in the schema", k)
		}

		if v.Optional || v.Required {
			return nil, fmt.Errorf("%q is a user-specifyable field - this should be specified as an Argument", k)
		}

		// every attribute has to be computed
		v.Computed = true
		out[k] = v
	}

	return &out, nil
}

func runArgs(d *schema.ResourceData, meta interface{}, logger Logger) ResourceMetaData {
	client := meta.(*clients.Client)
	metaData := ResourceMetaData{
		Client:                   client,
		Logger:                   logger,
		ResourceData:             d,
		serializationDebugLogger: NullLogger{},
	}

	return metaData
}
