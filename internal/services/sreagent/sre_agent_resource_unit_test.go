// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-10-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const testIdentityID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.ManagedIdentity/userAssignedIdentities/actions"

func TestSreAgentSchema(t *testing.T) {
	r := SreAgentResource{}
	wrapped := sdk.WrappedResource(r)
	if err := wrapped.InternalValidate(nil, true); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"vnet_configuration", "subnet_id", "agent_space_id"} {
		if _, ok := wrapped.Schema[name]; ok {
			t.Fatalf("unexpected unsupported field %q", name)
		}
	}
	if !r.Attributes()["default_model"].Computed {
		t.Fatal("server-selected default model must be computed")
	}
	for _, name := range []string{"action_configuration", "resources_configuration", "identity", "tags"} {
		if wrapped.Schema[name].ForceNew {
			t.Fatalf("%s must not invent replacement behavior", name)
		}
	}
	validate := wrapped.Schema["name"].ValidateFunc
	for _, value := range []string{"ab", "Agent-1", "a" + strings.Repeat("b", 31)} {
		if _, errs := validate(value, "name"); len(errs) != 0 {
			t.Fatalf("valid name %q rejected: %v", value, errs)
		}
	}
	for _, value := range []string{"a", "1agent", "agent-", "agent_name", "a" + strings.Repeat("b", 32)} {
		if _, errs := validate(value, "name"); len(errs) == 0 {
			t.Fatalf("invalid name %q accepted", value)
		}
	}
}

func TestSreAgentPatchOwnedFields(t *testing.T) {
	config := SreAgentModel{
		ActionConfiguration: []SreAgentActionConfiguration{{
			IdentityID: testIdentityID, Mode: "Review", AccessLevel: "Low",
		}},
		ResourcesConfiguration: []SreAgentResourcesConfiguration{{
			IdentityID:  testIdentityID + "-resources",
			ResourceIDs: []string{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/scoped"},
		}},
	}
	cases := []struct {
		name    string
		changed string
		want    string
	}{
		{"tags removed", "tags", `{"tags":{}}`},
		{"unrelated fields unchanged", "", `{}`},
		{"complete action", "action_configuration", `{"properties":{"actionConfiguration":{"accessLevel":"Low","identity":"` + testIdentityID + `","mode":"Review"}}}`},
		{"complete resource scope", "resources_configuration", `{"properties":{"knowledgeGraphConfiguration":{"identity":"` + testIdentityID + `-resources","managedResources":["/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/scoped"]}}}`},
		{"simultaneous nested edits", "both", `{"properties":{"actionConfiguration":{"accessLevel":"Low","identity":"` + testIdentityID + `","mode":"Review"},"knowledgeGraphConfiguration":{"identity":"` + testIdentityID + `-resources","managedResources":["/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/scoped"]}}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			patch, err := SreAgentResource{}.expandPatch(config, func(key string) bool {
				return key == tc.changed || tc.changed == "both" && (key == "action_configuration" || key == "resources_configuration")
			})
			if err != nil {
				t.Fatal(err)
			}
			got, err := json.Marshal(patch)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, got, tc.want)
		})
	}
}

func TestSreAgentRejectsEmptyResourceScope(t *testing.T) {
	config := SreAgentModel{
		Identity:               []identity.ModelSystemAssignedUserAssigned{{Type: identity.TypeSystemAssigned}},
		ResourcesConfiguration: []SreAgentResourcesConfiguration{{IdentityID: testIdentityID}},
	}
	if _, err := (SreAgentResource{}).expandCreate(config); err == nil {
		t.Fatal("create must reject an empty resource scope")
	}
	if _, err := (SreAgentResource{}).expandPatch(config, func(key string) bool { return key == "resources_configuration" }); err == nil {
		t.Fatal("update must not translate scope removal into an empty array")
	}
}

func TestSreAgentPatchRejectsUnsupportedBlockRemoval(t *testing.T) {
	for _, block := range []string{"action_configuration", "resources_configuration"} {
		t.Run(block, func(t *testing.T) {
			_, err := SreAgentResource{}.expandPatch(SreAgentModel{}, func(key string) bool { return key == block })
			if err == nil || !strings.Contains(err.Error(), block) {
				t.Fatalf("expected explicit removal error for %s, got %v", block, err)
			}
		})
	}
}

func TestSreAgentCreateIdentitySeparation(t *testing.T) {
	config := SreAgentModel{
		Location: "East US",
		Identity: []identity.ModelSystemAssignedUserAssigned{{
			Type: identity.TypeUserAssigned, IdentityIds: []string{testIdentityID, testIdentityID + "-resources"},
		}},
		ActionConfiguration: []SreAgentActionConfiguration{{
			IdentityID: testIdentityID, Mode: "Review", AccessLevel: "High",
		}},
		ResourcesConfiguration: []SreAgentResourcesConfiguration{{
			IdentityID:  testIdentityID + "-resources",
			ResourceIDs: []string{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/three", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/one", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/two"},
		}},
	}
	model, err := SreAgentResource{}.expandCreate(config)
	if err != nil {
		t.Fatal(err)
	}
	_, hasRootIdentity := model.Identity.IdentityIds[testIdentityID]
	if model.Location != "eastus" || !hasRootIdentity {
		t.Fatal("incorrect create envelope")
	}
	if pointer.From(model.Properties.ActionConfiguration.Identity) != testIdentityID ||
		pointer.From(model.Properties.KnowledgeGraphConfiguration.Identity) != testIdentityID+"-resources" {
		t.Fatal("identity references were collapsed")
	}
	if !reflect.DeepEqual(pointer.From(model.Properties.KnowledgeGraphConfiguration.ManagedResources), config.ResourcesConfiguration[0].ResourceIDs) {
		t.Fatal("resource IDs lost during mapping")
	}
	if model.Properties.DefaultModel != nil {
		t.Fatal("create must leave the server-selected model unspecified")
	}
}

func TestSreAgentMinimalCreateOmitsOptionalSettings(t *testing.T) {
	model, err := (SreAgentResource{}).expandCreate(SreAgentModel{
		Location: "eastus", Identity: []identity.ModelSystemAssignedUserAssigned{{Type: identity.TypeSystemAssigned}},
	})
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(model)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, body, `{"location":"eastus","identity":{"type":"SystemAssigned","userAssignedIdentities":null}}`)
}

func TestSreAgentRequiresManagedIdentity(t *testing.T) {
	for _, values := range [][]identity.ModelSystemAssignedUserAssigned{
		nil,
		{{Type: identity.TypeNone}},
		{{Type: identity.TypeUserAssigned}},
		{{Type: identity.TypeSystemAssignedUserAssigned}},
	} {
		config := SreAgentModel{Location: "eastus", Identity: values}
		if _, err := (SreAgentResource{}).expandCreate(config); err == nil {
			t.Fatalf("create accepted missing identity: %#v", values)
		}
		if _, err := (SreAgentResource{}).expandPatch(config, func(key string) bool { return key == "identity" }); err == nil {
			t.Fatalf("update accepted missing identity: %#v", values)
		}
	}
	if !(SreAgentResource{}).Arguments()["identity"].Required {
		t.Fatal("identity must be required in configuration")
	}
}

func TestSreAgentFlatten(t *testing.T) {
	r := SreAgentResource{}
	id := agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent")
	md := sdk.NewResourceMetaData(nil, r)
	md.SetID(id)
	model := agents.Agent{
		Location: "East US",
		Properties: &agents.AgentProperties{
			AgentEndpoint: pointer.To("https://example.invalid"),
			DefaultModel: &agents.DefaultModel{
				Name: pointer.To("server-default"), Provider: pointer.To("server-provider"),
			},
		},
	}
	if err := r.flatten(md, &id, &model); err != nil {
		t.Fatal(err)
	}
	var state SreAgentModel
	if err := md.Decode(&state); err != nil {
		t.Fatal(err)
	}
	if state.Endpoint != "https://example.invalid" || len(state.DefaultModel) != 1 || state.DefaultModel[0].Name != "server-default" {
		t.Fatalf("missing computed fields: %#v", state)
	}
	model.Properties = nil
	if err := r.flatten(md, &id, &model); err != nil {
		t.Fatal(err)
	}
	if err := md.Decode(&state); err != nil {
		t.Fatal(err)
	}
	if state.Endpoint != "" || len(state.DefaultModel) != 0 || len(state.ActionConfiguration) != 0 {
		t.Fatalf("stale optional state: %#v", state)
	}
	if err := r.flatten(md, &id, nil); err == nil {
		t.Fatal("nil successful model must be rejected")
	}
}

func TestSreAgentResourceIDSet(t *testing.T) {
	schema := SreAgentResource{}.Arguments()["resources_configuration"].Elem.(*pluginsdk.Resource).Schema["resource_ids"]
	if schema.Type != pluginsdk.TypeSet {
		t.Fatal("resource IDs must have order-independent Terraform state")
	}
}

func assertJSONEqual(t *testing.T, got []byte, want string) {
	t.Helper()
	var actual, expected interface{}
	if err := json.Unmarshal(got, &actual); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(want), &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unexpected JSON\n got: %s\nwant: %s", got, want)
	}
}
