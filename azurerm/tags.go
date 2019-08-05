package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

// TODO: deprecate and move these into deprecated.go

// NOTE: these methods are deprecated, but provided to ease compatibility for open PR's

func tagsSchema() *schema.Schema {
	return tags.Schema()
}

func tagsForceNewSchema() *schema.Schema {
	return tags.ForceNewSchema()
}

func tagsForDataSourceSchema() *schema.Schema {
	return tags.DataSourceSchema()
}

func tagValueToString(v interface{}) (string, error) {
	return tags.TagValueToString(v)
}

func validateAzureRMTags(v interface{}, k string) (warnings []string, errors []error) {
	return tags.Validate(v, k)
}

func expandTags(tagsMap map[string]interface{}) map[string]*string {
	return tags.Expand(tagsMap)
}

func filterTags(tagsMap map[string]*string, tagNames ...string) map[string]*string {
	return tags.Filter(tagsMap, tagNames...)
}

func flattenAndSetTags(d *schema.ResourceData, tagMap map[string]*string) {
	// we intentionally ignore the error here, since this method doesn't expose it
	tags.FlattenAndSet(d, tagMap)
}
