// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-01-01/agents"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sreclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sreagent/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type sreAgentTransport func(*http.Request) (*http.Response, error)

func (f sreAgentTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func sreAgentTestMetadata(t *testing.T, transport sreAgentTransport) (sdk.ResourceMetaData, agents.AgentId, context.Context) {
	t.Helper()
	client, err := agents.NewAgentsClientWithBaseURI(environments.AzurePublic().ResourceManager)
	if err != nil {
		t.Fatal(err)
	}
	client.Client.SetTransport(transport)
	client.Client.DisableRetries = true
	client.Client.AuthorizeRequest = nil
	id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
	md := sdk.NewResourceMetaData(&clients.Client{
		Account:  &clients.ResourceManagerAccount{SubscriptionId: id.SubscriptionId},
		Features: features.UserFeatures{PersistIDOnCreateBeforePollingForCompletion: true},
		SreAgent: &sreclient.Client{Agents: client},
	}, SreAgentResource{})
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)
	return md, id, ctx
}

func sreAgentHTTPResponse(request *http.Request, status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header: http.Header{"Content-Type": []string{"application/json"}, sdkclient.SkipPollingDelayHeader: []string{"true"}},
		Body:   io.NopCloser(strings.NewReader(body)), Request: request, ContentLength: int64(len(body)),
	}
}

func TestSreAgentReadHTTP(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNotFound, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodGet || req.URL.Query().Get("api-version") != "2026-01-01" {
					t.Fatalf("unexpected request %s %s", req.Method, req.URL)
				}
				body := `{"location":"eastus","properties":{"actionConfiguration":{"mode":"ReadOnly","accessLevel":"Low"},"knowledgeGraphConfiguration":{"managedResources":[]},"defaultModel":{"name":"server","provider":"provider"}}}`
				if status != http.StatusOK {
					body = `{"error":{"code":"Failure","message":"test error"}}`
				}
				return sreAgentHTTPResponse(req, status, body), nil
			})
			md.SetID(id)
			err := (SreAgentResource{}).Read().Func(ctx, md)
			switch status {
			case http.StatusNotFound:
				if err != nil || md.ResourceData.Id() != "" {
					t.Fatalf("404 should clear ID: %v", err)
				}
			case http.StatusForbidden:
				if err == nil || md.ResourceData.Id() != id.ID() {
					t.Fatalf("403 must surface error and retain ID: %v", err)
				}
			default:
				if err != nil {
					t.Fatal(err)
				}
				var model SreAgentModel
				if err := md.Decode(&model); err != nil {
					t.Fatal(err)
				}
				if model.ActionConfiguration[0].Mode != "ReadOnly" || len(model.ResourcesConfiguration[0].ResourceIDs) != 0 {
					t.Fatal("read must preserve actual server values independently of configuration restrictions")
				}
			}
		})
	}
}

func TestSreAgentUpdateUsesSelectivePATCH(t *testing.T) {
	calls := 0
	md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		calls++
		if req.Method == http.MethodGet {
			return sreAgentHTTPResponse(req, http.StatusOK, `{"location":"eastus","identity":{"type":"SystemAssigned"},"tags":{"purpose":"changed"},"properties":{"provisioningState":"Succeeded","actionConfiguration":{"mode":"ReadOnly"},"knowledgeGraphConfiguration":{"managedResources":[]}}}`), nil
		}
		if req.Method != http.MethodPatch {
			t.Fatalf("update must only use PATCH, got %s", req.Method)
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		assertJSONEqual(t, body, `{"tags":{"purpose":"changed"}}`)
		return sreAgentHTTPResponse(req, http.StatusOK, `{"properties":{"provisioningState":"Succeeded"}}`), nil
	})
	md.SetID(id)
	r := SreAgentResource{}
	if err := r.flatten(md, &id, &agents.Agent{
		Location: "eastus",
		Identity: &identity.LegacySystemAndUserAssignedMap{Type: identity.TypeSystemAssigned},
		Properties: &agents.AgentProperties{
			ActionConfiguration:         &agents.ActionConfiguration{Mode: pointer.To(agents.AgentModeReadOnly)},
			KnowledgeGraphConfiguration: &agents.KnowledgeGraphConfiguration{ManagedResources: pointer.To([]string{})},
		},
	}); err != nil {
		t.Fatal(err)
	}
	wrapped := sdk.WrappedResource(r)
	state := md.ResourceData.State()
	diff, err := wrapped.Diff(ctx, state, terraform.NewResourceConfigRaw(map[string]interface{}{
		"name": "agent", "resource_group_name": "example", "location": "eastus", "tags": map[string]interface{}{"purpose": "changed"},
		"identity": sreAgentSystemIdentityConfig(),
	}), md.Client)
	if err != nil {
		t.Fatal(err)
	}
	updated, diags := wrapped.Apply(ctx, state, diff, md.Client)
	if diags.HasError() {
		t.Fatalf("update failed: %v", diags)
	}
	if updated.Attributes["action_configuration.0.mode"] != "ReadOnly" || updated.Attributes["resources_configuration.0.resource_ids.#"] != "0" {
		t.Fatalf("unrelated PATCH must preserve unsupported server configuration: %#v", updated.Attributes)
	}
	if calls != 3 {
		t.Fatalf("expected PATCH, polling GET and refresh GET, got %d requests", calls)
	}
}

func TestSreAgentResourceScopeRoundTrip(t *testing.T) {
	prefix := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/"
	a, b, c, d := prefix+"one", prefix+"two", prefix+"three", prefix+"four"
	desired := []string{a, b, c}
	serverModel := agents.Agent{
		Location: "eastus",
		Identity: &identity.LegacySystemAndUserAssignedMap{Type: identity.TypeSystemAssigned},
		Properties: &agents.AgentProperties{
			ProvisioningState: pointer.To(agents.AgentProvisioningStateSucceeded),
			KnowledgeGraphConfiguration: &agents.KnowledgeGraphConfiguration{
				Identity: pointer.To(testIdentityID), ManagedResources: pointer.To([]string{c, a, b}),
			},
		},
	}
	patches := 0
	md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		switch req.Method {
		case http.MethodPatch:
			patches++
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return nil, err
			}
			props, ok := body["properties"].(map[string]interface{})
			if !ok || len(body) != 1 || len(props) != 1 {
				return nil, fmt.Errorf("scope update included unrelated settings: %#v", body)
			}
			scope, ok := props["knowledgeGraphConfiguration"].(map[string]interface{})
			if !ok || len(scope) != 2 || scope["identity"] != testIdentityID {
				return nil, fmt.Errorf("incomplete scope update: %#v", props)
			}
			rawIDs, ok := scope["managedResources"].([]interface{})
			if !ok {
				return nil, fmt.Errorf("managedResources is not an array: %#v", scope)
			}
			got := make([]string, len(rawIDs))
			for i, raw := range rawIDs {
				var ok bool
				got[i], ok = raw.(string)
				if !ok {
					return nil, fmt.Errorf("resource ID is not a string: %#v", raw)
				}
			}
			want := slices.Clone(desired)
			slices.Sort(got)
			slices.Sort(want)
			if !slices.Equal(got, want) {
				return nil, fmt.Errorf("scope membership = %v, want %v", got, want)
			}
		case http.MethodGet:
		default:
			return nil, fmt.Errorf("unexpected request %s", req.Method)
		}
		body, err := json.Marshal(serverModel)
		if err != nil {
			return nil, err
		}
		return sreAgentHTTPResponse(req, http.StatusOK, string(body)), nil
	})
	md.SetID(id)
	r := SreAgentResource{}
	if err := r.flatten(md, &id, &serverModel); err != nil {
		t.Fatal(err)
	}
	wrapped := sdk.WrappedResource(r)
	state := md.ResourceData.State()
	for _, step := range []struct {
		name    string
		ids     []string
		changed bool
	}{
		{"reorder three", []string{b, c, a}, false},
		{"replace one of three", []string{b, d, a}, true},
		{"contract to two", []string{d, a}, true},
		{"reorder two", []string{a, d}, false},
	} {
		t.Run(step.name, func(t *testing.T) {
			desired = step.ids
			readback := slices.Clone(desired)
			slices.Reverse(readback)
			serverModel.Properties.KnowledgeGraphConfiguration.ManagedResources = &readback
			configIDs := make([]interface{}, len(desired))
			for i, id := range desired {
				configIDs[i] = id
			}
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name": "agent", "resource_group_name": "example", "location": "eastus", "identity": sreAgentSystemIdentityConfig(),
				"resources_configuration": []interface{}{map[string]interface{}{"identity_id": testIdentityID, "resource_ids": configIDs}},
			})
			diff, err := wrapped.Diff(ctx, state, config, md.Client)
			if err != nil {
				t.Fatal(err)
			}
			changed := diff != nil && len(diff.Attributes) != 0
			if changed != step.changed || (diff != nil && (diff.RequiresNew() || diff.Destroy)) {
				t.Fatalf("unexpected scope diff: %#v", diff)
			}
			before := patches
			if step.changed {
				updated, diags := wrapped.Apply(ctx, state, diff, md.Client)
				if diags.HasError() {
					t.Fatal(diags)
				}
				state = updated
				if patches != before+1 {
					t.Fatalf("expected one scope PATCH, got %d", patches-before)
				}
			}
			noChange, err := wrapped.Diff(ctx, state, config, md.Client)
			if err != nil {
				t.Fatal(err)
			}
			if noChange != nil && (len(noChange.Attributes) != 0 || noChange.RequiresNew() || noChange.Destroy ||
				!reflect.DeepEqual(noChange.Identity, state.Identity)) {
				t.Fatalf("scope readback produced drift: %#v", noChange)
			}
			if state.Attributes["resources_configuration.0.resource_ids.#"] != fmt.Sprint(len(desired)) {
				t.Fatalf("unexpected scope size: %#v", state.Attributes)
			}
		})
	}
	if patches != 2 {
		t.Fatalf("order-only changes must not issue PATCH: got %d", patches)
	}
}

func TestSreAgentNestedConfigurationUpdates(t *testing.T) {
	for _, mode := range []agents.AgentMode{agents.AgentModeAutonomous, agents.AgentModeReview} {
		t.Run(string(mode), func(t *testing.T) {
			actionID, resourcesID := testIdentityID+"-actions", testIdentityID+"-resources"
			scopeID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/managed"
			previousMode := agents.AgentModeReview
			if mode == previousMode {
				previousMode = agents.AgentModeAutonomous
			}
			patches := 0
			md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				switch req.Method {
				case http.MethodPatch:
					patches++
					body, err := io.ReadAll(req.Body)
					if err != nil {
						return nil, err
					}
					assertJSONEqual(t, body, fmt.Sprintf(`{"properties":{"actionConfiguration":{"identity":%q,"mode":%q,"accessLevel":"Low"},"knowledgeGraphConfiguration":{"identity":%q,"managedResources":[%q]}}}`, actionID, mode, resourcesID, scopeID))
				case http.MethodGet:
				default:
					return nil, fmt.Errorf("unexpected request %s", req.Method)
				}
				return sreAgentHTTPResponse(req, http.StatusOK, fmt.Sprintf(`{"location":"eastus","identity":{"type":"SystemAssigned"},"properties":{"provisioningState":"Succeeded","actionConfiguration":{"identity":%q,"mode":%q,"accessLevel":"Low"},"knowledgeGraphConfiguration":{"identity":%q,"managedResources":[%q]}}}`, actionID, mode, resourcesID, scopeID)), nil
			})
			md.SetID(id)
			r := SreAgentResource{}
			if err := r.flatten(md, &id, &agents.Agent{
				Location: "eastus", Identity: &identity.LegacySystemAndUserAssignedMap{Type: identity.TypeSystemAssigned},
				Properties: &agents.AgentProperties{
					ActionConfiguration: &agents.ActionConfiguration{
						Identity: pointer.To(testIdentityID), Mode: &previousMode, AccessLevel: pointer.To(agents.AgentAccessLevelLow),
					},
					KnowledgeGraphConfiguration: &agents.KnowledgeGraphConfiguration{
						Identity: pointer.To(testIdentityID), ManagedResources: pointer.To([]string{scopeID}),
					},
				},
			}); err != nil {
				t.Fatal(err)
			}
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name": "agent", "resource_group_name": "example", "location": "eastus", "identity": sreAgentSystemIdentityConfig(),
				"action_configuration":    []interface{}{map[string]interface{}{"identity_id": actionID, "mode": string(mode), "access_level": "Low"}},
				"resources_configuration": []interface{}{map[string]interface{}{"identity_id": resourcesID, "resource_ids": []interface{}{scopeID}}},
			})
			wrapped := sdk.WrappedResource(r)
			diff, err := wrapped.Diff(ctx, md.ResourceData.State(), config, md.Client)
			if err != nil {
				t.Fatal(err)
			}
			state, diags := wrapped.Apply(ctx, md.ResourceData.State(), diff, md.Client)
			if diags.HasError() {
				t.Fatal(diags)
			}
			if patches != 1 || state.Attributes["action_configuration.0.identity_id"] != actionID ||
				state.Attributes["resources_configuration.0.identity_id"] != resourcesID || state.Attributes["action_configuration.0.mode"] != string(mode) {
				t.Fatalf("nested configuration update lost desired values: patches=%d state=%#v", patches, state.Attributes)
			}
		})
	}
}

func TestSreAgentCreateRequiresImport(t *testing.T) {
	md, _, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Fatalf("existing resource must not be overwritten: %s", req.Method)
		}
		return sreAgentHTTPResponse(req, http.StatusOK, `{"location":"eastus"}`), nil
	})
	for key, value := range map[string]string{"name": "agent", "resource_group_name": "example", "location": "eastus"} {
		if err := md.ResourceData.Set(key, value); err != nil {
			t.Fatal(err)
		}
	}
	if err := (SreAgentResource{}).Create().Func(ctx, md); err == nil || !strings.Contains(err.Error(), "import") {
		t.Fatalf("expected import error, got %v", err)
	}
}

func TestSreAgentCreateCapturesIdentityBeforePolling(t *testing.T) {
	for _, failed := range []bool{false, true} {
		t.Run(fmt.Sprintf("failed=%t", failed), func(t *testing.T) {
			var md sdk.ResourceMetaData
			var id agents.AgentId
			var ctx context.Context
			polled := false
			md, id, ctx = sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.Path, "/operations/") {
					polled = true
					if md.ResourceData.Id() != id.ID() {
						t.Error("ID must be persisted before polling")
						return nil, fmt.Errorf("ID missing before polling")
					}
					identityData, err := md.ResourceData.Identity()
					if err != nil {
						t.Error(err)
						return nil, err
					}
					if identityData.Get("subscription_id") != id.SubscriptionId || identityData.Get("name") != id.AgentName {
						t.Errorf("resource identity missing before polling: %#v", identityData)
						return nil, fmt.Errorf("resource identity missing")
					}
					if failed {
						return sreAgentHTTPResponse(req, http.StatusOK, `{"status":"Failed","error":{"code":"TestFailure","message":"synthetic failure"}}`), nil
					}
					return sreAgentHTTPResponse(req, http.StatusOK, `{"status":"Succeeded"}`), nil
				}
				switch req.Method {
				case http.MethodGet:
					return sreAgentHTTPResponse(req, http.StatusNotFound, `{"error":{"code":"ResourceNotFound","message":"absent"}}`), nil
				case http.MethodPut:
					result := sreAgentHTTPResponse(req, http.StatusCreated, `{"properties":{"provisioningState":"Accepted"}}`)
					result.Header.Set("Azure-AsyncOperation", "https://management.azure.com/operations/example?api-version=2026-01-01")
					return result, nil
				default:
					t.Errorf("unexpected %s request", req.Method)
					return nil, fmt.Errorf("unexpected request")
				}
			})
			for key, value := range map[string]string{"name": id.AgentName, "resource_group_name": id.ResourceGroupName, "location": "eastus"} {
				if err := md.ResourceData.Set(key, value); err != nil {
					t.Fatal(err)
				}
				if err := md.ResourceData.Set("identity", sreAgentSystemIdentityConfig()); err != nil {
					t.Fatal(err)
				}
			}
			err := (SreAgentResource{}).Create().Func(ctx, md)
			if (err != nil) != failed {
				t.Fatalf("unexpected create result: %v", err)
			}
			if !polled || md.ResourceData.Id() != id.ID() {
				t.Fatal("polling must retain the ID on both success and failure")
			}
		})
	}
}

func TestSreAgentDeleteAbsent(t *testing.T) {
	md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodGet {
			return sreAgentHTTPResponse(req, http.StatusNotFound, `{"error":{"code":"ResourceNotFound","message":"absent"}}`), nil
		}
		if req.Method != http.MethodDelete {
			return nil, fmt.Errorf("unexpected %s request", req.Method)
		}
		return sreAgentHTTPResponse(req, http.StatusNoContent, ""), nil
	})
	md.SetID(id)
	if err := (SreAgentResource{}).Delete().Func(ctx, md); err != nil {
		t.Fatal(err)
	}
}

func TestSreAgentDeleteInitialFailure(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodDelete {
					return nil, fmt.Errorf("unexpected %s request", req.Method)
				}
				return sreAgentHTTPResponse(req, status, `{"error":{"code":"Failure","message":"synthetic failure"}}`), nil
			})
			md.SetID(id)
			err := (SreAgentResource{}).Delete().Func(ctx, md)
			if status == http.StatusNotFound && err != nil {
				t.Fatalf("an already absent agent must delete successfully: %v", err)
			}
			if status == http.StatusForbidden && err == nil {
				t.Fatal("403 must not be treated as successful deletion")
			}
		})
	}
}

func TestSreAgentServerDefaultsDoNotDrift(t *testing.T) {
	r := SreAgentResource{}
	id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
	md := sdk.NewResourceMetaData(nil, r)
	md.SetID(id)
	model := agents.Agent{
		Location: "eastus",
		Identity: &identity.LegacySystemAndUserAssignedMap{Type: identity.TypeSystemAssigned},
		Properties: &agents.AgentProperties{
			ActionConfiguration: &agents.ActionConfiguration{
				Mode: pointer.To(agents.AgentModeReadOnly), AccessLevel: pointer.To(agents.AgentAccessLevelLow), Identity: pointer.To(testIdentityID),
			},
			KnowledgeGraphConfiguration: &agents.KnowledgeGraphConfiguration{Identity: pointer.To(testIdentityID), ManagedResources: pointer.To([]string{})},
			DefaultModel:                &agents.DefaultModel{Name: pointer.To("server"), Provider: pointer.To("provider")},
		},
	}
	if err := r.flatten(md, &id, &model); err != nil {
		t.Fatal(err)
	}
	base := map[string]interface{}{"name": "agent", "resource_group_name": "example", "location": "eastus", "identity": sreAgentSystemIdentityConfig()}
	wrapped := sdk.WrappedResource(r)
	diff, err := wrapped.Diff(context.Background(), md.ResourceData.State(), terraform.NewResourceConfigRaw(base), &clients.Client{})
	if err != nil {
		t.Fatal(err)
	}
	if diff != nil && (len(diff.Attributes) != 0 || diff.RequiresNew() || diff.Destroy || !reflect.DeepEqual(diff.Identity, md.ResourceData.State().Identity)) {
		t.Fatalf("omitted server defaults produced drift: attributes=%#v identity=%#v", diff.Attributes, diff.Identity)
	}
	for _, block := range []string{"action_configuration", "resources_configuration"} {
		t.Run(block+" reset", func(t *testing.T) {
			config := map[string]interface{}{"name": "agent", "resource_group_name": "example", "location": "eastus", "identity": sreAgentSystemIdentityConfig(), block: []interface{}{}}
			_, err := wrapped.Diff(context.Background(), md.ResourceData.State(), terraform.NewResourceConfigRaw(config), &clients.Client{})
			if err == nil || !strings.Contains(err.Error(), block) {
				t.Fatalf("explicit reset must fail at plan time, got %v", err)
			}
		})
		t.Run(block+" unknown", func(t *testing.T) {
			if _, err := wrapped.Diff(context.Background(), md.ResourceData.State(), sreAgentUnknownConfig(wrapped, block, false), &clients.Client{}); err != nil {
				t.Fatalf("unknown block is not a reset request: %v", err)
			}
		})
		for _, empty := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s protocol empty=%t", block, empty), func(t *testing.T) {
				values := protocolValues(wrapped)
				blockType := wrapped.CoreConfigSchema().ImpliedType().AttributeType(block)
				if empty {
					values[block] = cty.ListValEmpty(blockType.ElementType())
				}
				protocol := terraform.NewResourceConfigShimmed(cty.ObjectVal(values), wrapped.CoreConfigSchema())
				diff, err := wrapped.Diff(context.Background(), md.ResourceData.State(), protocol, &clients.Client{})
				if err != nil {
					t.Fatal(err)
				}
				if diff != nil && (len(diff.Attributes) != 0 || diff.RequiresNew() || diff.Destroy || !reflect.DeepEqual(diff.Identity, md.ResourceData.State().Identity)) {
					t.Fatalf("an omitted or zero-length dynamic block must preserve remote settings: %#v", diff)
				}
			})
		}
	}
}

func TestSreAgentLiteralEmptyBlockIsNotConfigurationSyntax(t *testing.T) {
	core := sdk.WrappedResource(SreAgentResource{}).CoreConfigSchema()
	bodySchema := &hcl.BodySchema{}
	for name := range core.Attributes {
		bodySchema.Attributes = append(bodySchema.Attributes, hcl.AttributeSchema{Name: name})
	}
	for name := range core.BlockTypes {
		bodySchema.Blocks = append(bodySchema.Blocks, hcl.BlockHeaderSchema{Type: name})
	}
	for _, name := range []string{"action_configuration", "resources_configuration"} {
		t.Run(name, func(t *testing.T) {
			file, diags := hclsyntax.ParseConfig([]byte(name+" = []"), "test.tf", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatal(diags)
			}
			if _, diags := file.Body.Content(bodySchema); !diags.HasErrors() {
				t.Fatalf("%s must use nested-block syntax; assignment is not a reset interface", name)
			}
		})
	}
}

func TestSreAgentUnknownNestedValuesRemainPlannable(t *testing.T) {
	wrapped := sdk.WrappedResource(SreAgentResource{})
	for _, block := range []string{"action_configuration", "resources_configuration"} {
		t.Run(block, func(t *testing.T) {
			if _, err := wrapped.Diff(context.Background(), nil, sreAgentUnknownConfig(wrapped, block, true), &clients.Client{}); err != nil {
				t.Fatalf("unknown references must remain plannable: %v", err)
			}
		})
	}
}

func sreAgentUnknownConfig(resource *pluginsdk.Resource, block string, nested bool) *terraform.ResourceConfig {
	schema := resource.CoreConfigSchema()
	values := protocolValues(resource)
	values[block] = cty.UnknownVal(schema.ImpliedType().AttributeType(block))
	if nested {
		fields := map[string]cty.Value{"identity_id": cty.UnknownVal(cty.String)}
		if block == "action_configuration" {
			fields["mode"] = cty.StringVal("Review")
			fields["access_level"] = cty.StringVal("Low")
		} else {
			fields["resource_ids"] = cty.SetVal([]cty.Value{cty.UnknownVal(cty.String)})
		}
		values[block] = cty.ListVal([]cty.Value{cty.ObjectVal(fields)})
	}
	return terraform.NewResourceConfigShimmed(cty.ObjectVal(values), schema)
}

func protocolValues(resource *pluginsdk.Resource) map[string]cty.Value {
	values := make(map[string]cty.Value)
	for name, attributeType := range resource.CoreConfigSchema().ImpliedType().AttributeTypes() {
		values[name] = cty.NullVal(attributeType)
	}
	values["name"] = cty.StringVal("agent")
	values["resource_group_name"] = cty.StringVal("example")
	values["location"] = cty.StringVal("eastus")
	values["identity"] = cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
		"type":         cty.StringVal("SystemAssigned"),
		"identity_ids": cty.SetValEmpty(cty.String),
		"principal_id": cty.NullVal(cty.String),
		"tenant_id":    cty.NullVal(cty.String),
	})})
	return values
}

func sreAgentSystemIdentityConfig() []interface{} {
	return []interface{}{map[string]interface{}{"type": "SystemAssigned"}}
}
