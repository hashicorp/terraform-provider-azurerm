// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package differ

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
	schema_rules "github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/schema-rules"
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
				baseItem = providerjson.SchemaJSON{}
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
				baseItem = providerjson.SchemaJSON{}
			}
			if errs := compareNodeDataSource(baseItem, propertySchema, propertyName); errs != nil {
				violations = append(violations, errs...)
			}
		}
	}

	return violations
}

func compareNodeResource(base providerjson.SchemaJSON, current providerjson.SchemaJSON, nodeName string) (errs []string) {
	if nodeIsBlock(base) {
		newBaseRaw := base.Elem.(providerjson.ResourceJSON).Schema
		newCurrent := current.Elem.(*providerjson.ResourceJSON).Schema
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

func compareNodeDataSource(base providerjson.SchemaJSON, current providerjson.SchemaJSON, nodeName string) (errs []string) {
	if nodeIsBlock(base) {
		newBaseRaw := base.Elem.(providerjson.ResourceJSON).Schema
		newCurrent := current.Elem.(*providerjson.ResourceJSON).Schema
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

func nodeIsBlock(input providerjson.SchemaJSON) bool {
	if input.Type == providerjson.SchemaTypeList || input.Type == providerjson.SchemaTypeSet {
		if _, ok := input.Elem.(providerjson.ResourceJSON); ok {
			return true
		}
	}

	return false
}
