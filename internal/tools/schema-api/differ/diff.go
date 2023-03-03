package differ

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/schema-rules"
)

type Differ struct {
	base    *providerjson.ProviderWrapper
	current *providerjson.ProviderWrapper
}

func (d *Differ) Diff(fileName string, providerName string) []string {
	if err := d.loadFromProvider(providerjson.LoadData(), providerName); err != nil {
		return []string{err.Error()}
	}

	if err := d.loadFromFile(fileName); err != nil {
		return []string{err.Error()}
	}

	if d.base.ProviderName != d.current.ProviderName {
		return []string{fmt.Sprintf("provider name mismatch, expected %q, got %q", d.base.ProviderName, d.current.ProviderName)}
	}

	// TODO - Walk the resources
	violations := make([]string, 0)
	for resource, rs := range d.current.ProviderSchema.ResourcesMap {
		_, ok := d.base.ProviderSchema.ResourcesMap[resource]
		if !ok {
			// New resource, no breaking changes to worry about
			continue
		}
		for propertyName, propertySchema := range rs.Schema {
			// Get the same from the base (released) json
			baseItem, ok := d.base.ProviderSchema.ResourcesMap[resource].Schema[propertyName]
			if !ok {
				// TODO - New property, Could be Required which would be breaking, need to account for this
				continue
			}
			if errs := compareNode(baseItem, propertySchema, propertyName); errs != nil {
				violations = append(violations, errs...)
			}
		}
	}

	// TODO - walk the data sources

	return violations
}

func compareNode(base providerjson.SchemaJSON, current providerjson.SchemaJSON, nodeName string) (errs []string) {
	if nodeIsBlock(base) {
		newBaseRaw := base.Elem.(map[string]interface{})["schema"].(map[string]interface{})
		newCurrent := current.Elem.(*providerjson.ResourceJSON).Schema
		for k, v := range newBaseRaw {
			newBase := providerjson.SchemaFromMap(v.(map[string]interface{}))
			errs = append(errs, compareNode(newBase, newCurrent[k], k)...)
		}
	}

	for _, v := range schema_rules.BreakingChangeRules {
		if err := v.Check(base, current, nodeName); err != nil {
			errs = append(errs, *err)
		}
	}

	return
}

func nodeIsBlock(input providerjson.SchemaJSON) bool {
	if input.Type == "TypeList" || input.Type == "TypeSet" {
		if elem, ok := input.Elem.(map[string]interface{}); ok {
			if _, ok := elem["schema"]; ok {
				return true
			}
		}
	}
	return false
}
