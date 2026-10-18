// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2024-11-01/extensions"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	containersclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
)

func TestKubernetesClusterExtensionUpdate(t *testing.T) {
	cases := []struct {
		name                string
		previous            map[string]interface{}
		settings            map[string]string
		protected           map[string]string
		previousProtected   map[string]interface{}
		advertisedProtected map[string]string
		omitProtectedMap    bool
		method              string
		errorContains       string
	}{
		{name: "remove one key", previous: map[string]interface{}{"keep": "old", "remove": "true"}, settings: map[string]string{"keep": "new"}, method: http.MethodPut},
		{name: "remove last key", previous: map[string]interface{}{"remove": "true"}, method: http.MethodPut},
		{name: "change value uses patch", previous: map[string]interface{}{"keep": "old"}, settings: map[string]string{"keep": "new"}, method: http.MethodPatch},
		{name: "add key uses patch", previous: map[string]interface{}{}, settings: map[string]string{"add": "true"}, method: http.MethodPatch},
		{name: "preserve protected value from state", previous: map[string]interface{}{"remove": "true"}, protected: map[string]string{"protected": "fixture-value"}, previousProtected: map[string]interface{}{"protected": "fixture-value"}, advertisedProtected: map[string]string{"protected": ""}, method: http.MethodPut},
		{name: "rotate protected value with key removal", previous: map[string]interface{}{"remove": "true"}, protected: map[string]string{"protected": "new-fixture-value"}, previousProtected: map[string]interface{}{"protected": "old-fixture-value"}, advertisedProtected: map[string]string{"protected": ""}, method: http.MethodPut},
		{name: "supply previously unmanaged protected value", previous: map[string]interface{}{"remove": "true"}, protected: map[string]string{"protected": "fixture-value"}, advertisedProtected: map[string]string{"protected": ""}, method: http.MethodPut},
		{name: "reject missing protected map", previous: map[string]interface{}{"remove": "true"}, omitProtectedMap: true, errorContains: "cannot be verified"},
		{name: "reject unknown protected value", previous: map[string]interface{}{"remove": "true"}, advertisedProtected: map[string]string{"unmanaged": ""}, errorContains: "not available in Terraform state"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clusterID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.ContainerService/managedClusters/test"
			id := extensions.NewScopedExtensionID(clusterID, "test")
			current := map[string]string{}
			for k, v := range tc.previous {
				current[k] = v.(string)
			}
			advertised := tc.advertisedProtected
			if advertised == nil {
				advertised = map[string]string{}
			}
			writes := []string{}
			payload := map[string]interface{}{}
			var mutex sync.Mutex
			response := func() map[string]interface{} {
				result := map[string]interface{}{
					"id": id.ID(), "name": "test", "type": "Microsoft.KubernetesConfiguration/extensions",
					"systemData": map[string]interface{}{"createdBy": "not-a-writable-field"},
					"identity":   map[string]interface{}{"type": "SystemAssigned", "principalId": "generated", "tenantId": "generated"},
					"plan":       map[string]interface{}{"name": "plan", "product": "product", "publisher": "publisher"},
					"properties": map[string]interface{}{
						"extensionType": "Microsoft.ArgoCD", "autoUpgradeMinorVersion": false, "releaseTrain": "Stable", "version": "1.0.19", "currentVersion": "1.0.19", "scope": map[string]interface{}{"cluster": map[string]interface{}{"releaseNamespace": "argocd"}},
						"configurationSettings": current, "configurationProtectedSettings": advertised,
						"provisioningState": "Succeeded", "packageUri": "not-a-writable-field", "customLocationSettings": map[string]string{"generated": "true"}, "isSystemExtension": false,
						"aksAssignedIdentity": map[string]interface{}{"type": "SystemAssigned", "principalId": "generated", "tenantId": "generated"},
					},
				}
				if tc.omitProtectedMap {
					delete(result["properties"].(map[string]interface{}), "configurationProtectedSettings")
				}
				return result
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				mutex.Lock()
				defer mutex.Unlock()
				if req.URL.Path != id.ID() || req.URL.Query().Get("api-version") != "2024-11-01" {
					t.Errorf("unexpected request target")
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set(sdkclient.SkipPollingDelayHeader, "true")
				switch req.Method {
				case http.MethodGet:
				case http.MethodPut, http.MethodPatch:
					writes = append(writes, req.Method)
					if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
						t.Errorf("decoding request: %v", err)
					}
					props := payload["properties"].(map[string]interface{})
					if req.Method == http.MethodPut {
						current = map[string]string{}
					}
					if m, ok := props["configurationSettings"].(map[string]interface{}); ok {
						for k, v := range m {
							current[k] = v.(string)
						}
					}
				default:
					t.Errorf("unexpected HTTP method %s", req.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				if err := json.NewEncoder(w).Encode(response()); err != nil {
					t.Errorf("encoding response: %v", err)
				}
			}))
			defer server.Close()
			client, err := extensions.NewExtensionsClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			resource := KubernetesClusterExtensionResource{}
			metadata := sdk.NewResourceMetaData(&clients.Client{Containers: &containersclient.Client{KubernetesExtensionsClient: client}}, resource)
			protected := tc.protected
			if protected == nil {
				protected = map[string]string{}
			}
			settings := tc.settings
			if settings == nil {
				settings = map[string]string{}
			}
			previousSettings := make(map[string]string, len(tc.previous))
			for k, v := range tc.previous {
				previousSettings[k] = v.(string)
			}
			previousProtected := make(map[string]string, len(tc.previousProtected))
			for k, v := range tc.previousProtected {
				previousProtected[k] = v.(string)
			}
			if err := metadata.Encode(&KubernetesClusterExtensionModel{Name: "test", ClusterID: clusterID, ExtensionType: "Microsoft.ArgoCD", Version: "1.0.19", ReleaseNamespace: "argocd", ConfigurationSettings: previousSettings, ConfigurationProtectedSettings: previousProtected}); err != nil {
				t.Fatal(err)
			}
			metadata.ResourceData.SetId(id.ID())
			previousState := metadata.ResourceData.State()
			desired := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name": "test", "cluster_id": clusterID, "extension_type": "Microsoft.ArgoCD", "version": "1.0.19", "release_namespace": "argocd",
				"configuration_settings": stringMapForExtensionTest(settings), "configuration_protected_settings": stringMapForExtensionTest(protected),
			})
			schemaMap := schema.InternalMap(sdk.WrappedResource(resource).SchemaMap())
			diff, err := schemaMap.Diff(t.Context(), previousState, desired, nil, nil, false)
			if err != nil {
				t.Fatal(err)
			}
			metadata.ResourceData, err = schemaMap.Data(previousState, diff)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
			defer cancel()
			err = resource.Update().Func(ctx, metadata)
			mutex.Lock()
			defer mutex.Unlock()
			if tc.errorContains != "" {
				if err == nil || !strings.Contains(err.Error(), tc.errorContains) {
					t.Fatalf("expected safety error containing %q, got %v", tc.errorContains, err)
				}
				if len(writes) != 0 {
					t.Fatalf("unsafe write occurred after protected-settings check failed")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(writes, []string{tc.method}) {
				t.Fatalf("expected exactly one %s, got %v", tc.method, writes)
			}
			if !reflect.DeepEqual(current, settings) {
				t.Error("Azure configuration did not converge to the requested map")
			}
			if metadata.ResourceData.Id() != id.ID() {
				t.Error("resource ID changed")
			}
			props := payload["properties"].(map[string]interface{})
			if tc.method == http.MethodPut {
				if !reflect.DeepEqual(props["configurationProtectedSettings"], stringMapForExtensionTest(protected)) {
					t.Error("protected values were not preserved from configuration/state")
				}
				if props["extensionType"] != "Microsoft.ArgoCD" || props["version"] != "1.0.19" || props["releaseTrain"] != "Stable" || props["autoUpgradeMinorVersion"] != false {
					t.Error("extension version or upgrade policy changed")
				}
				if !reflect.DeepEqual(props["scope"], map[string]interface{}{"cluster": map[string]interface{}{"releaseNamespace": "argocd"}}) {
					t.Error("extension scope changed")
				}
				if !reflect.DeepEqual(payload["identity"], map[string]interface{}{"type": "SystemAssigned"}) || !reflect.DeepEqual(props["aksAssignedIdentity"], map[string]interface{}{"type": "SystemAssigned"}) {
					t.Error("identity type changed or generated IDs were replayed")
				}
				if !reflect.DeepEqual(payload["plan"], map[string]interface{}{"name": "plan", "product": "product", "publisher": "publisher"}) {
					t.Error("marketplace plan changed")
				}
				for _, key := range []string{"id", "name", "type", "systemData"} {
					if _, ok := payload[key]; ok {
						t.Errorf("read-only field %s was replayed", key)
					}
				}
				for _, key := range []string{"currentVersion", "provisioningState", "packageUri", "customLocationSettings", "isSystemExtension", "statuses", "errorInfo"} {
					if _, ok := props[key]; ok {
						t.Errorf("read-only property %s was replayed", key)
					}
				}
			} else if _, ok := props["configurationProtectedSettings"]; ok {
				t.Error("unchanged protected settings should not be patched")
			}
		})
	}
}

func stringMapForExtensionTest(input map[string]string) map[string]interface{} {
	result := make(map[string]interface{}, len(input))
	for k, v := range input {
		result[k] = v
	}
	return result
}

func TestKubernetesClusterExtensionFullUpdateRejectsIncompleteResponse(t *testing.T) {
	for _, existing := range []extensions.Extension{{}, {Properties: &extensions.ExtensionProperties{}}} {
		if _, err := expandKubernetesClusterExtensionFullUpdate(existing, nil, nil, nil); err == nil {
			t.Fatal("expected incomplete GET to prevent PUT")
		}
	}
}

func TestKubernetesClusterExtensionFullUpdatePreservesNamespaceScope(t *testing.T) {
	extensionType := "Microsoft.ArgoCD"
	upgrade := true
	train := "Preview"
	namespace := "target-namespace"
	generated := "generated-principal"
	empty := map[string]string{}
	existing := extensions.Extension{
		Properties: &extensions.ExtensionProperties{
			ExtensionType:                  &extensionType,
			AutoUpgradeMinorVersion:        &upgrade,
			ReleaseTrain:                   &train,
			Scope:                          &extensions.Scope{Namespace: &extensions.ScopeNamespace{TargetNamespace: &namespace}},
			ConfigurationProtectedSettings: &empty,
			AksAssignedIdentity:            &extensions.ExtensionPropertiesAksAssignedIdentity{PrincipalId: &generated},
		},
	}
	payload, err := expandKubernetesClusterExtensionFullUpdate(existing, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(payload.Properties.Scope, existing.Properties.Scope) || payload.Properties.AutoUpgradeMinorVersion != existing.Properties.AutoUpgradeMinorVersion || payload.Properties.ReleaseTrain != existing.Properties.ReleaseTrain || payload.Properties.Version != nil {
		t.Error("namespace scope or automatic upgrade policy changed")
	}
	if payload.Properties.AksAssignedIdentity != nil {
		t.Error("generated identity-only object must not be replayed")
	}
	if payload.Properties.ConfigurationSettings == nil || *payload.Properties.ConfigurationSettings == nil || payload.Properties.ConfigurationProtectedSettings == nil || *payload.Properties.ConfigurationProtectedSettings == nil {
		t.Error("empty settings must be explicit objects, not null")
	}
}
