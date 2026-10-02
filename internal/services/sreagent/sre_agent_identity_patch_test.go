// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-10-01/agents"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestSreAgentIdentityPatchWire(t *testing.T) {
	a, b, c := testIdentityID, testIdentityID+"-b", testIdentityID+"-c"
	cases := []struct {
		name    string
		kind    identity.Type
		old     []string
		desired []string
		want    map[string]interface{}
	}{
		{"remove one", identity.TypeUserAssigned, []string{a, b}, []string{b}, map[string]interface{}{a: nil, b: map[string]interface{}{}}},
		{"add and remove", identity.TypeUserAssigned, []string{a, b}, []string{b, c}, map[string]interface{}{a: nil, b: map[string]interface{}{}, c: map[string]interface{}{}}},
		{"case only", identity.TypeUserAssigned, []string{a}, []string{strings.ToUpper(a)}, map[string]interface{}{strings.ToUpper(a): map[string]interface{}{}}},
		{"retain both", identity.TypeUserAssigned, []string{a, b}, []string{a, b}, map[string]interface{}{a: map[string]interface{}{}, b: map[string]interface{}{}}},
		{"combined remove one", identity.TypeSystemAssignedUserAssigned, []string{a, b}, []string{b}, map[string]interface{}{a: nil, b: map[string]interface{}{}}},
		{"system only", identity.TypeSystemAssigned, []string{a, b}, nil, nil},
		{"add to system", identity.TypeSystemAssignedUserAssigned, nil, []string{c}, map[string]interface{}{c: map[string]interface{}{}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			desired := testSreAgentIdentity(tc.kind, tc.desired)
			previous := testSreAgentIdentity(identity.TypeSystemAssignedUserAssigned, tc.old)
			patches := 0
			expected := map[string]interface{}{
				"identity": map[string]interface{}{
					"type": strings.ReplaceAll(string(tc.kind), ", ", ","), "userAssignedIdentities": tc.want,
				},
				"tags": map[string]string{"stage": "updated"},
			}
			input := agents.AgentPatch{
				Identity: desired, Tags: pointer.To(map[string]string{"stage": "updated"}),
			}
			if len(tc.desired) > 0 {
				actionID := tc.desired[0]
				expected["properties"] = map[string]interface{}{"actionConfiguration": map[string]interface{}{
					"identity": actionID, "mode": "Review", "accessLevel": "Low",
				}}
				input.Properties = &agents.AgentPatchProperties{ActionConfiguration: &agents.ActionConfiguration{
					Identity: pointer.To(actionID), Mode: pointer.To(agents.AgentModeReview), AccessLevel: pointer.To(agents.AgentAccessLevelLow),
				}}
			}
			want, err := json.Marshal(expected)
			if err != nil {
				t.Fatal(err)
			}
			md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				if req.URL.Query().Get("api-version") != "2026-10-01" {
					return nil, fmt.Errorf("unexpected API version: %s", req.URL)
				}
				switch req.Method {
				case http.MethodPatch:
					patches++
					body, err := io.ReadAll(req.Body)
					if err != nil {
						return nil, err
					}
					assertJSONEqual(t, body, string(want))
				case http.MethodGet:
				default:
					return nil, fmt.Errorf("unexpected method: %s", req.Method)
				}
				return sreAgentHTTPResponse(req, http.StatusOK, `{"properties":{"provisioningState":"Succeeded"}}`), nil
			})
			if err := updateSreAgent(ctx, md.Client.SreAgent.Agents, id, input, previous); err != nil {
				t.Fatal(err)
			}
			if patches != 1 {
				t.Fatalf("expected exactly one PATCH, got %d", patches)
			}
		})
	}
}

func TestSreAgentIdentityPatchPollingFailure(t *testing.T) {
	patches, polls := 0, 0
	md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		switch {
		case req.Method == http.MethodPatch:
			patches++
			result := sreAgentHTTPResponse(req, http.StatusAccepted, `{"properties":{"provisioningState":"Accepted"}}`)
			result.Header.Set("Azure-AsyncOperation", "https://management.azure.com/operations/identity-update?api-version=2026-10-01")
			return result, nil
		case req.Method == http.MethodGet && req.URL.Path == "/operations/identity-update":
			polls++
			return sreAgentHTTPResponse(req, http.StatusOK, `{"status":"Failed","error":{"code":"IdentityUpdateFailed","message":"synthetic identity failure"}}`), nil
		default:
			return nil, fmt.Errorf("unexpected polling request %s %s", req.Method, req.URL)
		}
	})
	err := updateSreAgent(ctx, md.Client.SreAgent.Agents, id, agents.AgentPatch{
		Identity: testSreAgentIdentity(identity.TypeSystemAssignedUserAssigned, []string{testIdentityID}),
	}, testSreAgentIdentity(identity.TypeSystemAssigned, nil))
	if err == nil || !strings.Contains(err.Error(), "IdentityUpdateFailed") || patches != 1 || polls != 1 {
		t.Fatalf("failed identity LRO must surface the service error: patches=%d polls=%d err=%v", patches, polls, err)
	}
}

func TestSreAgentIdentityRemovalUsesPriorState(t *testing.T) {
	a, b := testIdentityID, testIdentityID+"-b"
	patches := 0
	md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodPatch {
			patches++
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			assertJSONEqual(t, body, fmt.Sprintf(`{"identity":{"type":"UserAssigned","userAssignedIdentities":{%q:null,%q:{}}}}`, a, b))
		} else if req.Method != http.MethodGet {
			return nil, fmt.Errorf("unexpected method: %s", req.Method)
		}
		return sreAgentHTTPResponse(req, http.StatusOK, fmt.Sprintf(`{"location":"eastus","identity":{"type":"UserAssigned","userAssignedIdentities":{%q:{}}},"properties":{"provisioningState":"Succeeded"}}`, b)), nil
	})
	md.SetID(id)
	resource := SreAgentResource{}
	if err := resource.flatten(md, &id, &agents.Agent{Location: "eastus", Identity: testSreAgentIdentity(identity.TypeUserAssigned, []string{a, b})}); err != nil {
		t.Fatal(err)
	}
	wrapped := sdk.WrappedResource(resource)
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name": "agent", "resource_group_name": "example", "location": "eastus",
		"identity": []interface{}{map[string]interface{}{"type": "UserAssigned", "identity_ids": []interface{}{b}}},
	})
	diff, err := wrapped.Diff(ctx, md.ResourceData.State(), config, md.Client)
	if err != nil {
		t.Fatal(err)
	}
	state, diags := wrapped.Apply(ctx, md.ResourceData.State(), diff, md.Client)
	if diags.HasError() {
		t.Fatal(diags)
	}
	if patches != 1 || state.Attributes["identity.0.identity_ids.#"] != "1" {
		t.Fatalf("partial removal was not applied: patches=%d state=%#v", patches, state.Attributes)
	}
}

func TestSreAgentLastIdentityRemovalRejectedDuringPlan(t *testing.T) {
	md, id, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
		t.Fatalf("invalid identity plan must not issue %s", req.Method)
		return nil, nil
	})
	md.SetID(id)
	r := SreAgentResource{}
	if err := r.flatten(md, &id, &agents.Agent{Location: "eastus", Identity: testSreAgentIdentity(identity.TypeSystemAssigned, nil)}); err != nil {
		t.Fatal(err)
	}
	wrapped := sdk.WrappedResource(r)
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name": "agent", "resource_group_name": "example", "location": "eastus", "identity": []interface{}{},
	})
	if _, err := wrapped.Diff(ctx, md.ResourceData.State(), config, md.Client); err == nil {
		t.Fatal("removal of the last managed identity must fail during planning")
	}
}

func testSreAgentIdentity(kind identity.Type, ids []string) *identity.LegacySystemAndUserAssignedMap {
	value := &identity.LegacySystemAndUserAssignedMap{Type: kind, IdentityIds: map[string]identity.UserAssignedIdentityDetails{}}
	for _, id := range ids {
		value.IdentityIds[id] = identity.UserAssignedIdentityDetails{}
	}
	return value
}
