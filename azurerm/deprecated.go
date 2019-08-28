package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
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

func validateAzureRMTags(v interface{}, k string) (warnings []string, errors []error) {
	return tags.Validate(v, k)
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
func filterTags(tagsMap map[string]*string, tagNames ...string) map[string]*string {
	return tags.Filter(tagsMap, tagNames...)
}

// nolint: deadcode unused
func tagsForceNewSchema() *schema.Schema {
	return tags.ForceNewSchema()
}

func parseAzureResourceID(id string) (*azure.ResourceID, error) {
	return azure.ParseAzureResourceID(id)
}

// nolint: deadcode unused
func azureRMLockByName(name string, resourceType string) {
	locks.ByName(name, resourceType)
}

// nolint: deadcode unused
func azureRMUnlockByName(name string, resourceType string) {
	locks.UnlockByName(name, resourceType)
}
