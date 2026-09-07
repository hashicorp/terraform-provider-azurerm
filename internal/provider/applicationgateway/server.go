// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

// Package applicationgateway normalizes SDK plan artifacts for Application
// Gateway nested blocks without changing planning for other AzureRM resources.
package applicationgateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const resourceName = "azurerm_application_gateway"

// Wrap preserves absent computed outputs when the SDK changes only their
// representation while reconstructing a set. All other RPCs are delegated.
func Wrap(server tfprotov5.ProviderServer, resource *schema.Resource) tfprotov5.ProviderServer {
	if resource == nil {
		return server
	}
	return &planServer{ProviderServer: server, resourceType: resource.CoreConfigSchema().ImpliedType(), resource: resource}
}

type planServer struct {
	tfprotov5.ProviderServer
	resourceType cty.Type
	resource     *schema.Resource
}

func (s *planServer) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	resp, err := s.ProviderServer.PlanResourceChange(ctx, req)
	if err != nil || resp == nil || req.TypeName != resourceName || req.PriorState == nil || resp.PlannedState == nil || len(resp.RequiresReplace) > 0 || resp.Deferred != nil {
		return resp, err
	}
	for _, diagnostic := range resp.Diagnostics {
		if diagnostic.Severity == tfprotov5.DiagnosticSeverityError {
			return resp, nil
		}
	}
	prior, err := decode(req.PriorState, s.resourceType)
	if err != nil {
		return nil, fmt.Errorf("decoding Application Gateway prior state: %w", err)
	}
	planned, err := decode(resp.PlannedState, s.resourceType)
	if err != nil {
		return nil, fmt.Errorf("decoding Application Gateway planned state: %w", err)
	}
	if prior.IsNull() || planned.IsNull() || !prior.IsKnown() || !planned.IsKnown() {
		return resp, nil
	}
	// Restrict normalization to updates of the same gateway.
	if !prior.GetAttr("id").IsKnown() || !prior.GetAttr("id").RawEquals(planned.GetAttr("id")) {
		return resp, nil
	}
	attrs := planned.AsValueMap()
	changed := false
	for key, field := range s.resource.Schema {
		if field.Computed && !field.Optional {
			continue
		}
		if _, ok := field.Elem.(*schema.Resource); !ok || (field.Type != schema.TypeSet && field.Type != schema.TypeList) {
			continue
		}
		value := normalizeCollection(prior.GetAttr(key), attrs[key], field)
		if !value.RawEquals(attrs[key]) {
			attrs[key] = value
			changed = true
		}
	}
	if !changed {
		return resp, nil
	}
	encoded, err := msgpack.Marshal(cty.ObjectVal(attrs), s.resourceType)
	if err != nil {
		return nil, fmt.Errorf("encoding Application Gateway planned state: %w", err)
	}
	resp.PlannedState = &tfprotov5.DynamicValue{MsgPack: encoded}
	return resp, nil
}

func decode(value *tfprotov5.DynamicValue, typ cty.Type) (cty.Value, error) {
	if len(value.MsgPack) > 0 {
		return msgpack.Unmarshal(value.MsgPack, typ)
	}
	return ctyjson.Unmarshal(value.JSON, typ)
}

// normalizeCollection matches named set members by name and list members by
// position. The entire candidate block must match its prior value before any
// computed output is restored; a matching name alone is never sufficient.
func normalizeCollection(prior, planned cty.Value, field *schema.Schema) cty.Value {
	if prior.IsNull() || planned.IsNull() || !prior.IsKnown() || !planned.IsKnown() || prior.RawEquals(planned) {
		return planned
	}
	resource, ok := field.Elem.(*schema.Resource)
	if !ok {
		return planned
	}
	values := planned.AsValueSlice()
	if len(values) == 0 {
		return planned
	}
	changed := false
	switch field.Type {
	case schema.TypeSet:
		if _, named := resource.Schema["name"]; !named {
			return planned
		}
		before, valid := namedMembers(prior)
		if !valid {
			return planned
		}
		if _, valid := namedMembers(planned); !valid {
			return planned
		}
		for i, value := range values {
			if value.IsNull() || !value.IsKnown() {
				continue
			}
			name := value.GetAttr("name")
			if name.IsNull() || !name.IsKnown() {
				continue
			}
			if old, exists := before[name.AsString()]; exists {
				values[i] = normalizeBlock(old, value, resource)
				changed = changed || !values[i].RawEquals(value)
			}
		}
		if changed {
			return cty.SetVal(values)
		}
	case schema.TypeList:
		for i, value := range values {
			if i >= prior.LengthInt() {
				continue
			}
			values[i] = normalizeBlock(prior.Index(cty.NumberIntVal(int64(i))), value, resource)
			changed = changed || !values[i].RawEquals(value)
		}
		if changed {
			return cty.ListVal(values)
		}
	}
	return planned
}

func namedMembers(collection cty.Value) (map[string]cty.Value, bool) {
	result := map[string]cty.Value{}
	for it := collection.ElementIterator(); it.Next(); {
		_, value := it.Element()
		if value.IsNull() || !value.IsKnown() {
			return nil, false
		}
		name := value.GetAttr("name")
		if name.IsNull() || !name.IsKnown() {
			return nil, false
		}
		if _, duplicate := result[name.AsString()]; duplicate {
			return nil, false
		}
		result[name.AsString()] = value
	}
	return result, true
}

func normalizeBlock(prior, planned cty.Value, resource *schema.Resource) cty.Value {
	if prior.IsNull() || planned.IsNull() || !prior.IsWhollyKnown() || !planned.IsKnown() || prior.RawEquals(planned) {
		return planned
	}
	candidate := planned.AsValueMap()
	for key, field := range resource.Schema {
		old, next := prior.GetAttr(key), candidate[key]
		switch {
		case equivalentEmpty(old, next):
			// The legacy SDK treats null and zero values as equivalent, but the
			// resulting set elements differ to Terraform Core.
			candidate[key] = old
		case field.Computed && !field.Optional && field.Type == schema.TypeString && absentString(old) && !next.IsKnown():
			// Do not restore a missing ID when the corresponding named object is
			// configured: that output is expected to become populated.
			reference := strings.TrimSuffix(key, "_id") + "_name"
			if key == "id" {
				reference = "name"
			}
			if _, exists := resource.Schema[reference]; exists && !absentString(prior.GetAttr(reference)) {
				continue
			}
			candidate[key] = old
		case field.Type == schema.TypeSet || field.Type == schema.TypeList:
			candidate[key] = normalizeCollection(old, next, field)
		}
	}
	if cty.ObjectVal(candidate).RawEquals(prior) {
		return prior
	}
	return planned
}

func equivalentEmpty(a, b cty.Value) bool {
	if !a.IsKnown() || !b.IsKnown() || a.RawEquals(b) {
		return false
	}
	if a.IsNull() {
		return zeroValue(b)
	}
	if b.IsNull() {
		return zeroValue(a)
	}
	return false
}

func zeroValue(value cty.Value) bool {
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

func absentString(value cty.Value) bool {
	return value.IsKnown() && (value.IsNull() || value.RawEquals(cty.StringVal("")))
}
