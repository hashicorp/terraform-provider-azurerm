// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package applicationgateway

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestPlanResourceChange(t *testing.T) {
	for _, profile := range []string{"null", "empty", "populated"} {
		for _, action := range []string{"none", "edit", "add", "remove", "rename", "errors", "profile"} {
			t.Run(profile+"/"+action, func(t *testing.T) {
				t.Parallel()
				p, resource := testProvider()
				before := []cty.Value{testListener(t, resource, "private", profile), testListener(t, resource, "orbit", profile)}
				after := append([]cty.Value(nil), before...)
				switch action {
				case "edit":
					attrs := after[1].AsValueMap()
					attrs["host_names"] = cty.SetVal([]cty.Value{cty.StringVal("*.orbit.example.com"), cty.StringVal("test.orbit.example.com")})
					after[1] = cty.ObjectVal(attrs)
				case "add":
					after = append(after, testListener(t, resource, "new", profile))
				case "remove":
					after = after[:1]
				case "rename":
					attrs := after[1].AsValueMap()
					attrs["name"] = cty.StringVal("renamed")
					after[1] = cty.ObjectVal(attrs)
				case "errors":
					attrs := after[1].AsValueMap()
					errors := attrs["custom_error_configuration"].AsValueSlice()
					errorAttrs := errors[0].AsValueMap()
					errorAttrs["custom_error_page_url"] = cty.StringVal("https://example.com/new-error.html")
					attrs["custom_error_configuration"] = cty.ListVal([]cty.Value{cty.ObjectVal(errorAttrs)})
					after[1] = cty.ObjectVal(attrs)
				case "profile":
					attrs := after[1].AsValueMap()
					attrs["ssl_profile_name"] = cty.StringVal("new-profile")
					after[1] = cty.ObjectVal(attrs)
				}
				req := testRequest(t, resource, before, after)
				response, err := Wrap(p.GRPCProvider(), resource).PlanResourceChange(t.Context(), req)
				planned := testPlanned(t, resource, response, err).GetAttr("http_listener")
				if planned.LengthInt() != len(after) {
					t.Fatalf("got %d listeners, want %d", planned.LengthInt(), len(after))
				}
				for it := planned.ElementIterator(); it.Next(); {
					_, listener := it.Element()
					name := listener.GetAttr("name").AsString()
					var expected cty.Value
					for _, want := range after {
						if want.GetAttr("name").AsString() == name {
							expected = want
						}
					}
					if expected == cty.NilVal {
						t.Fatalf("unexpected listener %q", name)
					}
					if name == "private" || action == "none" || (name == "orbit" && action == "add") {
						if !listener.RawEquals(expected) {
							t.Errorf("untouched listener %q changed:\n got: %#v\nwant: %#v", name, listener, expected)
						}
					} else {
						for _, key := range []string{"host_names", "ssl_profile_name"} {
							if !listener.GetAttr(key).RawEquals(expected.GetAttr(key)) && (!absentString(listener.GetAttr(key)) || !absentString(expected.GetAttr(key))) {
								t.Errorf("lost configured %s change for %q", key, name)
							}
						}
						if listener.GetAttr("id").IsKnown() {
							t.Errorf("changed listener %q should still recompute its ID", name)
						}
						if action == "errors" && !listener.GetAttr("custom_error_configuration").Index(cty.NumberIntVal(0)).GetAttr("custom_error_page_url").RawEquals(expected.GetAttr("custom_error_configuration").Index(cty.NumberIntVal(0)).GetAttr("custom_error_page_url")) {
							t.Error("lost custom error page edit")
						}
					}
				}
			})
		}
	}
}

func TestOtherResourceUnchanged(t *testing.T) {
	p, resource := testProvider()
	before := []cty.Value{testListener(t, resource, "private", "empty"), testListener(t, resource, "orbit", "empty")}
	after := before[:1]
	req := testRequest(t, resource, before, after)
	req.TypeName = "azurerm_other"
	raw, err := p.GRPCProvider().PlanResourceChange(t.Context(), req)
	want := testPlanned(t, resource, raw, err)
	got, err := Wrap(p.GRPCProvider(), resource).PlanResourceChange(t.Context(), req)
	if !testPlanned(t, resource, got, err).RawEquals(want) {
		t.Fatal("changed another resource's plan")
	}
}

func TestNormalizeListenersBoundaries(t *testing.T) {
	_, resource := testProvider()
	old := testListener(t, resource, "private", "empty")
	for _, name := range []string{"absent_profile", "configured_profile", "known_profile", "configured_port", "unknown_existing_id", "new_profile_id", "unknown_set", "empty_set"} {
		t.Run(name, func(t *testing.T) {
			prior := old
			attrs := old.AsValueMap()
			attrs["ssl_profile_id"] = cty.UnknownVal(cty.String)
			wantChanged := name == "absent_profile"
			switch name {
			case "configured_profile":
				priorAttrs := old.AsValueMap()
				priorAttrs["ssl_profile_name"] = cty.StringVal("profile")
				prior = cty.ObjectVal(priorAttrs)
				attrs["ssl_profile_name"] = cty.StringVal("profile")
			case "known_profile":
				prior = testListener(t, resource, "private", "populated")
				attrs = prior.AsValueMap()
				attrs["ssl_profile_id"] = cty.UnknownVal(cty.String)
			case "configured_port":
				attrs["frontend_port_name"] = cty.StringVal("other-port")
			case "unknown_existing_id":
				attrs["frontend_port_id"] = cty.UnknownVal(cty.String)
			case "new_profile_id":
				attrs["ssl_profile_id"] = cty.StringVal("/new/profile")
			}
			before := cty.SetVal([]cty.Value{prior})
			planned := cty.SetVal([]cty.Value{cty.ObjectVal(attrs)})
			switch name {
			case "unknown_set":
				planned = cty.UnknownVal(before.Type())
			case "empty_set":
				planned = cty.SetValEmpty(prior.Type())
			}
			got := normalizeCollection(before, planned, resource.Schema["http_listener"])
			changed := !got.RawEquals(planned)
			if changed != wantChanged {
				t.Fatalf("changed = %t, want %t", changed, wantChanged)
			}
			want := planned
			if wantChanged {
				want = before
			}
			if !got.RawEquals(want) {
				t.Fatalf("unexpected plan: %#v", got)
			}
		})
	}
}

func TestKnownComputedChangePreserved(t *testing.T) {
	p, resource := testProvider()
	before := []cty.Value{testListener(t, resource, "private", "empty")}
	req := testRequest(t, resource, before, before)
	delegate := &testServer{ProviderServer: p.GRPCProvider(), change: func(response *tfprotov5.PlanResourceChangeResponse) {
		planned := testPlanned(t, resource, response, nil)
		attrs := planned.AsValueMap()
		listener := before[0].AsValueMap()
		listener["frontend_port_id"] = cty.UnknownVal(cty.String)
		listener["ssl_profile_id"] = cty.UnknownVal(cty.String)
		attrs["http_listener"] = cty.SetVal([]cty.Value{cty.ObjectVal(listener)})
		response.PlannedState = testDynamic(t, cty.ObjectVal(attrs))
	}}
	response, err := Wrap(delegate, resource).PlanResourceChange(t.Context(), req)
	planned := testPlanned(t, resource, response, err).GetAttr("http_listener").AsValueSlice()[0]
	if planned.GetAttr("frontend_port_id").IsKnown() || planned.GetAttr("ssl_profile_id").IsKnown() {
		t.Fatal("normalized a listener with a real computed change")
	}
}

type testServer struct {
	tfprotov5.ProviderServer
	change func(*tfprotov5.PlanResourceChangeResponse)
}

func (s *testServer) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	response, err := s.ProviderServer.PlanResourceChange(ctx, req)
	if err == nil {
		s.change(response)
	}
	return response, err
}

func testListener(t *testing.T, resource *schema.Resource, name, profile string) cty.Value {
	t.Helper()
	fields := resource.Schema["http_listener"].Elem.(*schema.Resource).Schema
	attrs := map[string]interface{}{}
	for key := range fields {
		attrs[key] = nil
	}
	attrs["name"] = name
	attrs["frontend_ip_configuration_name"] = "private"
	attrs["frontend_port_name"] = "https"
	attrs["protocol"] = "Https"
	attrs["require_sni"] = true
	attrs["ssl_certificate_name"] = "cert"
	attrs["host_names"] = []string{"*." + name + ".example.com"}
	for _, key := range []string{"id", "frontend_ip_configuration_id", "frontend_port_id", "ssl_certificate_id"} {
		attrs[key] = "/fake/" + name + "/" + key
	}
	switch profile {
	case "empty":
		attrs["ssl_profile_id"] = ""
	case "populated":
		attrs["ssl_profile_id"] = "/fake/profile"
		attrs["ssl_profile_name"] = "profile"
	}
	attrs["custom_error_configuration"] = []interface{}{map[string]interface{}{"status_code": "HttpStatus403", "custom_error_page_url": "https://example.com/403.html", "id": attrs["ssl_profile_id"]}}
	encoded, err := json.Marshal(attrs)
	if err != nil {
		t.Fatal(err)
	}
	typ := resource.CoreConfigSchema().ImpliedType().AttributeType("http_listener").ElementType()
	value, err := ctyjson.Unmarshal(encoded, typ)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func testRequest(t *testing.T, resource *schema.Resource, before, after []cty.Value) *tfprotov5.PlanResourceChangeRequest {
	t.Helper()
	config := make([]cty.Value, len(after))
	proposed := make([]cty.Value, len(after))
	for i, listener := range after {
		attrs := listener.AsValueMap()
		for key, field := range resource.Schema["http_listener"].Elem.(*schema.Resource).Schema {
			if field.Computed {
				attrs[key] = cty.NullVal(cty.String)
			}
		}
		errorAttrs := attrs["custom_error_configuration"].Index(cty.NumberIntVal(0)).AsValueMap()
		errorAttrs["id"] = cty.NullVal(cty.String)
		attrs["custom_error_configuration"] = cty.ListVal([]cty.Value{cty.ObjectVal(errorAttrs)})
		config[i] = cty.ObjectVal(attrs)
		proposed[i] = config[i]
		for _, old := range before {
			if listener.RawEquals(old) {
				proposed[i] = listener
			}
		}
	}
	wrap := func(members []cty.Value, id cty.Value) *tfprotov5.DynamicValue {
		return testDynamic(t, cty.ObjectVal(map[string]cty.Value{"id": id, "http_listener": cty.SetVal(members)}))
	}
	return &tfprotov5.PlanResourceChangeRequest{TypeName: resourceName, PriorState: wrap(before, cty.StringVal("gateway")), Config: wrap(config, cty.NullVal(cty.String)), ProposedNewState: wrap(proposed, cty.StringVal("gateway"))}
}

func testDynamic(t *testing.T, value cty.Value) *tfprotov5.DynamicValue {
	t.Helper()
	encoded, err := msgpack.Marshal(value, value.Type())
	if err != nil {
		t.Fatal(err)
	}
	return &tfprotov5.DynamicValue{MsgPack: encoded}
}

func testPlanned(t *testing.T, resource *schema.Resource, response *tfprotov5.PlanResourceChangeResponse, err error) cty.Value {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range response.Diagnostics {
		if diagnostic.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("%s: %s", diagnostic.Summary, diagnostic.Detail)
		}
	}
	value, err := decode(response.PlannedState, resource.CoreConfigSchema().ImpliedType())
	if err != nil {
		t.Fatal(err)
	}
	return value
}
