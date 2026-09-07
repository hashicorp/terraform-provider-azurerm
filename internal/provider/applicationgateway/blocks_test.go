// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package applicationgateway

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
)

// Exercise every repeatable named block from the production gateway schema,
// including nested lists/sets. Azure client configuration is deliberately not
// needed to test how the SDK plans these blocks.
func TestGatewayBlocks(t *testing.T) {
	gateway := (network.Registration{}).SupportedResources()[resourceName]
	for key, field := range gateway.Schema {
		block, ok := field.Elem.(*schema.Resource)
		if !ok || (field.Computed && !field.Optional) || (field.Type != schema.TypeSet && field.Type != schema.TypeList) || field.MaxItems == 1 || block.Schema["name"] == nil {
			continue
		}
		actions := []string{"none", "rename", "add", "remove"}
		editKey := map[string]string{"backend_http_settings": "request_timeout", "request_routing_rule": "priority"}[key]
		if editKey != "" {
			actions = append(actions, "edit")
		}
		for _, empty := range []bool{false, true} {
			for _, action := range actions {
				t.Run(fmt.Sprintf("%s/empty=%t/%s", key, empty, action), func(t *testing.T) {
					t.Parallel()
					resource := &schema.Resource{Schema: map[string]*schema.Schema{key: field}}
					provider := &schema.Provider{ResourcesMap: map[string]*schema.Resource{resourceName: resource}}
					typ := resource.CoreConfigSchema().ImpliedType().AttributeType(key)
					before := []cty.Value{blockFixture(block, typ.ElementType(), "unchanged", empty), blockFixture(block, typ.ElementType(), "edited", empty)}
					after := append([]cty.Value(nil), before...)
					switch action {
					case "edit":
						attrs := after[1].AsValueMap()
						attrs[editKey] = cty.NumberIntVal(121)
						after[1] = cty.ObjectVal(attrs)
					case "rename":
						attrs := after[1].AsValueMap()
						attrs["name"] = cty.StringVal("renamed")
						after[1] = cty.ObjectVal(attrs)
					case "add":
						after = append(after, blockFixture(block, typ.ElementType(), "new", empty))
					case "remove":
						after = after[:1]
					}
					config, proposed := make([]cty.Value, len(after)), make([]cty.Value, len(after))
					for i, value := range after {
						config[i] = blockConfig(value, block)
						proposed[i] = config[i]
						for _, old := range before {
							if value.RawEquals(old) {
								proposed[i] = value
							}
						}
					}
					wrap := func(values []cty.Value, id cty.Value) *tfprotov5.DynamicValue {
						return testDynamic(t, cty.ObjectVal(map[string]cty.Value{"id": id, key: collectionValue(typ, values)}))
					}
					response, err := Wrap(provider.GRPCProvider(), resource).PlanResourceChange(t.Context(), &tfprotov5.PlanResourceChangeRequest{
						TypeName: resourceName, PriorState: wrap(before, cty.StringVal("gateway")), Config: wrap(config, cty.NullVal(cty.String)), ProposedNewState: wrap(proposed, cty.StringVal("gateway")),
					})
					planned := testPlanned(t, resource, response, err).GetAttr(key)
					if planned.LengthInt() != len(after) {
						t.Fatalf("got %d blocks, want %d", planned.LengthInt(), len(after))
					}
					for it := planned.ElementIterator(); it.Next(); {
						_, got := it.Element()
						name := got.GetAttr("name").AsString()
						found := false
						for _, want := range after {
							if name != want.GetAttr("name").AsString() {
								continue
							}
							found = true
							if action == "edit" && name == "edited" && !got.GetAttr(editKey).RawEquals(want.GetAttr(editKey)) {
								t.Errorf("lost configured %s change", editKey)
							}
							if name == "unchanged" || action == "none" || (name == "edited" && action == "add") {
								for attr, value := range got.AsValueMap() {
									if !value.RawEquals(want.GetAttr(attr)) {
										t.Errorf("untouched %s.%s.%s differs", key, name, attr)
									}
								}
							}
						}
						if !found {
							t.Errorf("unexpected block name %q", name)
						}
					}
				})
			}
		}
	}
}

func blockFixture(resource *schema.Resource, typ cty.Type, name string, empty bool) cty.Value {
	attrs := map[string]cty.Value{}
	for key, field := range resource.Schema {
		valueType := typ.AttributeType(key)
		value := cty.NullVal(valueType)
		if nested, ok := field.Elem.(*schema.Resource); ok && (field.Type == schema.TypeList || field.Type == schema.TypeSet) {
			value = collectionValue(valueType, []cty.Value{blockFixture(nested, valueType.ElementType(), name+"-child", empty)})
		} else if field.Default != nil {
			switch v := field.Default.(type) {
			case bool:
				value = cty.BoolVal(v)
			case int:
				value = cty.NumberIntVal(int64(v))
			case string:
				value = cty.StringVal(v)
			}
		} else if field.Required {
			switch field.Type {
			case schema.TypeString:
				if field.StateFunc != nil {
					value = cty.StringVal(field.StateFunc("example"))
				} else {
					value = cty.StringVal("example")
				}
			case schema.TypeInt:
				value = cty.NumberIntVal(1)
			case schema.TypeBool:
				value = cty.False
			case schema.TypeList, schema.TypeSet:
				if valueType.ElementType() == cty.String {
					value = collectionValue(valueType, []cty.Value{cty.StringVal("example")})
				}
			}
		}
		if key == "name" {
			value = cty.StringVal(name)
		}
		attrs[key] = value
	}
	for key, field := range resource.Schema {
		if !field.Computed || field.Optional || field.Type != schema.TypeString {
			continue
		}
		if empty {
			attrs[key] = cty.StringVal("")
		}
		reference := strings.TrimSuffix(key, "_id") + "_name"
		if key == "id" {
			reference = "name"
		}
		if value, ok := attrs[reference]; ok && !absentString(value) {
			attrs[key] = cty.StringVal("/fake/" + name + "/" + key)
		}
	}
	return cty.ObjectVal(attrs)
}

func blockConfig(value cty.Value, resource *schema.Resource) cty.Value {
	attrs := value.AsValueMap()
	for key, field := range resource.Schema {
		if field.Computed && !field.Optional {
			attrs[key] = cty.NullVal(attrs[key].Type())
			continue
		}
		if nested, ok := field.Elem.(*schema.Resource); ok && !attrs[key].IsNull() && (field.Type == schema.TypeSet || field.Type == schema.TypeList) {
			values := attrs[key].AsValueSlice()
			for i, v := range values {
				values[i] = blockConfig(v, nested)
			}
			attrs[key] = collectionValue(attrs[key].Type(), values)
		}
	}
	return cty.ObjectVal(attrs)
}

func collectionValue(typ cty.Type, values []cty.Value) cty.Value {
	if typ.IsSetType() {
		if len(values) == 0 {
			return cty.SetValEmpty(typ.ElementType())
		}
		return cty.SetVal(values)
	}
	if len(values) == 0 {
		return cty.ListValEmpty(typ.ElementType())
	}
	return cty.ListVal(values)
}
