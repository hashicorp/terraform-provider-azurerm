package parse

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

type FlexibleServerId struct {
	ResourceGroup string
	Name          string
}

func FlexibleServerID(input string) (*FlexibleServerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing postgresqlflexible Server ID %q: %+v", input, err)
	}

	flexibleServer := FlexibleServerId{
		ResourceGroup: id.ResourceGroup,
	}
	if flexibleServer.Name, err = id.PopSegment("flexibleServers"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &flexibleServer, nil
}

func FlattenAndSetPropertyTags(d *schema.ResourceData, tagMap map[string]*string) error {
	flattened := tags.Flatten(tagMap)
	if err := d.Set("properties_tags", flattened); err != nil {
		return fmt.Errorf("setting `properties_tags`: %s", err)
	}

	return nil
}
