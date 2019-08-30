package azurerm

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

// NOTE: these methods are deprecated, but provided to ease compatibility for open PR's

// nolint: deadcode unused
var requireResourcesToBeImported = features.ShouldResourcesBeImported()

// nolint: deadcode unused
func flattenAndSetTags(d *schema.ResourceData, tagMap map[string]*string) {
	// we intentionally ignore the error here, since this method doesn't expose it
	_ = tags.FlattenAndSet(d, tagMap)
}

// nolint: deadcode unused
func expandTags(tagsMap map[string]interface{}) map[string]*string {
	return tags.Expand(tagsMap)
}

// nolint: deadcode unused
func tagsForDataSourceSchema() *schema.Schema {
	return tags.SchemaDataSource()
}

// nolint: deadcode unused
func tagsSchema() *schema.Schema {
	return tags.Schema()
}

// nolint: deadcode unused
func tagsForceNewSchema() *schema.Schema {
	return tags.ForceNewSchema()
}

// nolint: deadcode unused
func parseAzureResourceID(id string) (*azure.ResourceID, error) {
	return azure.ParseAzureResourceID(id)
}

func evaluateSchemaValidateFunc(i interface{}, k string, validateFunc schema.SchemaValidateFunc) (bool, error) { // nolint: unparam
	_, errors := validateFunc(i, k)

	errorStrings := []string{}
	for _, e := range errors {
		errorStrings = append(errorStrings, e.Error())
	}

	if len(errors) > 0 {
		return false, fmt.Errorf(strings.Join(errorStrings, "\n"))
	}

	return true, nil
}
