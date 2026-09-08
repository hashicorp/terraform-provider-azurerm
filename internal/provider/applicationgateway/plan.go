// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package applicationgateway

import (
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func preserveUnchangedBlocks(prior, planned cty.Value, resource *schema.Resource) cty.Value {
	attributes := planned.AsValueMap()
	for name, field := range resource.Schema {
		if field.Computed && !field.Optional {
			continue
		}
		attributes[name] = preserveUnchangedCollection(prior.GetAttr(name), attributes[name], field)
	}
	return cty.ObjectVal(attributes)
}

func preserveUnchangedCollection(prior, planned cty.Value, field *schema.Schema) cty.Value {
	if prior.IsNull() || planned.IsNull() || !prior.IsKnown() || !planned.IsKnown() || prior.RawEquals(planned) {
		return planned
	}
	block, ok := field.Elem.(*schema.Resource)
	if !ok {
		return planned
	}
	switch field.Type {
	case schema.TypeSet:
		if block.Schema["name"] != nil {
			return preserveSetBlocksByName(prior, planned, block)
		}
	case schema.TypeList:
		return preserveListBlocksByPosition(prior, planned, block)
	}
	return planned
}

func preserveSetBlocksByName(prior, planned cty.Value, block *schema.Resource) cty.Value {
	priorByName, valid := indexBlocksByName(prior)
	if !valid {
		return planned
	}
	if _, valid := indexBlocksByName(planned); !valid {
		return planned
	}

	values := planned.AsValueSlice()
	changed := false
	for i, value := range values {
		if previous, exists := priorByName[value.GetAttr("name").AsString()]; exists {
			values[i] = preserveUnchangedBlock(previous, value, block)
			changed = changed || !values[i].RawEquals(value)
		}
	}
	if changed {
		return cty.SetVal(values)
	}
	return planned
}

func preserveListBlocksByPosition(prior, planned cty.Value, block *schema.Resource) cty.Value {
	values := planned.AsValueSlice()
	previous := prior.AsValueSlice()
	changed := false
	for i := 0; i < min(len(values), len(previous)); i++ {
		preserved := preserveUnchangedBlock(previous[i], values[i], block)
		changed = changed || !preserved.RawEquals(values[i])
		values[i] = preserved
	}
	if changed {
		return cty.ListVal(values)
	}
	return planned
}

func indexBlocksByName(collection cty.Value) (map[string]cty.Value, bool) {
	blocks := make(map[string]cty.Value)
	for it := collection.ElementIterator(); it.Next(); {
		_, value := it.Element()
		if value.IsNull() || !value.IsKnown() {
			return nil, false
		}
		name := value.GetAttr("name")
		if name.IsNull() || !name.IsKnown() {
			return nil, false
		}
		if _, duplicate := blocks[name.AsString()]; duplicate {
			return nil, false
		}
		blocks[name.AsString()] = value
	}
	return blocks, true
}

func preserveUnchangedBlock(prior, planned cty.Value, block *schema.Resource) cty.Value {
	if prior.IsNull() || planned.IsNull() || !prior.IsWhollyKnown() || !planned.IsKnown() || prior.RawEquals(planned) {
		return planned
	}
	for name, field := range block.Schema {
		previous, next := prior.GetAttr(name), planned.GetAttr(name)
		switch {
		case previous.RawEquals(next), equivalentEmptyValues(previous, next):
			continue
		case canPreserveAbsentComputedString(name, field, prior, next):
			continue
		}
		if !preserveUnchangedCollection(previous, next, field).RawEquals(previous) {
			return planned
		}
	}
	return prior
}

func canPreserveAbsentComputedString(name string, field *schema.Schema, prior cty.Value, planned cty.Value) bool {
	if !field.Computed || field.Optional || field.Type != schema.TypeString || planned.IsKnown() || !isAbsentString(prior.GetAttr(name)) {
		return false
	}
	return !hasConfiguredReference(name, prior)
}

func hasConfiguredReference(attribute string, block cty.Value) bool {
	reference := strings.TrimSuffix(attribute, "_id") + "_name"
	if attribute == "id" {
		reference = "name"
	}
	return block.Type().HasAttribute(reference) && !isAbsentString(block.GetAttr(reference))
}

func equivalentEmptyValues(prior, planned cty.Value) bool {
	if !prior.IsKnown() || !planned.IsKnown() {
		return false
	}
	return (prior.IsNull() && isZeroValue(planned)) || (planned.IsNull() && isZeroValue(prior))
}

func isZeroValue(value cty.Value) bool {
	if value.IsNull() {
		return true
	}
	switch {
	case value.Type() == cty.String:
		return value.RawEquals(cty.StringVal(""))
	case value.Type() == cty.Bool:
		return value.RawEquals(cty.False)
	case value.Type() == cty.Number:
		return value.RawEquals(cty.NumberIntVal(0))
	case value.Type().IsCollectionType():
		return value.LengthInt() == 0
	}
	return false
}

func isAbsentString(value cty.Value) bool {
	return value.IsKnown() && (value.IsNull() || value.RawEquals(cty.StringVal("")))
}
