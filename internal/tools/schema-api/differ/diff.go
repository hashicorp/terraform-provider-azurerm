// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package differ

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
	schema_rules "github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/schema-rules"
)

type Differ struct {
	base    *providerschema.ProviderWrapper
	current *providerschema.ProviderWrapper
}

func (d *Differ) Diff(fileName string, providerName string) []string {
	if err := d.loadFromProvider(providerschema.LoadData(), providerName); err != nil {
		return []string{err.Error()}
	}

	if err := d.loadFromFile(fileName); err != nil {
		return []string{err.Error()}
	}

	if d.base.ProviderName != d.current.ProviderName {
		return []string{fmt.Sprintf("provider name mismatch, expected %q, got %q", d.base.ProviderName, d.current.ProviderName)}
	}

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
				// New property, could be breaking - Required etc
				baseItem = providerschema.SchemaJSON{}
			}
			if errs := compareNodeResource(baseItem, propertySchema, propertyName); errs != nil {
				violations = append(violations, errs...)
			}
		}
	}

	for dataSource, ds := range d.current.ProviderSchema.DataSourcesMap {
		_, ok := d.base.ProviderSchema.DataSourcesMap[dataSource]
		if !ok {
			// New data source, no breaking changes to worry about
			continue
		}
		for propertyName, propertySchema := range ds.Schema {
			// Get the same from the base (released) json
			baseItem, ok := d.base.ProviderSchema.DataSourcesMap[dataSource].Schema[propertyName]
			if !ok {
				// New property, could be breaking - Required etc
				baseItem = providerschema.SchemaJSON{}
			}
			if errs := compareNodeDataSource(baseItem, propertySchema, propertyName); errs != nil {
				violations = append(violations, errs...)
			}
		}
	}

	return violations
}

func compareNodeResource(base providerschema.SchemaJSON, current providerschema.SchemaJSON, nodeName string) (errs []string) {
	if nodeIsBlock(base) {
		newBaseRaw := base.Elem.(providerschema.ResourceJSON).Schema
		newCurrent := current.Elem.(*providerschema.ResourceJSON).Schema
		for k, newBase := range newBaseRaw {
			errs = append(errs, compareNodeResource(newBase, newCurrent[k], k)...)
		}
	}

	for _, v := range schema_rules.BreakingChangeRules {
		if err := v.Check(base, current, nodeName); err != nil {
			errs = append(errs, *err)
		}
	}

	return
}

func compareNodeDataSource(base providerschema.SchemaJSON, current providerschema.SchemaJSON, nodeName string) (errs []string) {
	if nodeIsBlock(base) {
		newBaseRaw := base.Elem.(providerschema.ResourceJSON).Schema
		newCurrent := current.Elem.(*providerschema.ResourceJSON).Schema
		for k, newBase := range newBaseRaw {
			errs = append(errs, compareNodeDataSource(newBase, newCurrent[k], k)...)
		}
	}

	for _, v := range schema_rules.BreakingChangeRulesDataSource {
		if err := v.Check(base, current, nodeName); err != nil {
			errs = append(errs, *err)
		}
	}

	return
}

func nodeIsBlock(input providerschema.SchemaJSON) bool {
	if input.Type == providerschema.SchemaTypeList || input.Type == providerschema.SchemaTypeSet {
		if _, ok := input.Elem.(providerschema.ResourceJSON); ok {
			return true
		}
	}

	return false
}
