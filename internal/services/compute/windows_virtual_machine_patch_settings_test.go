// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	networksdk "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	computeclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	networkclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// This transport exercises the real SDK serializer, request builder and poller
// without opening a socket or requiring Azure credentials.
type windowsPatchSettingsTransport struct {
	t                *testing.T
	puts             []virtualmachines.VirtualMachine
	patches          []map[string]any
	patchSettings    map[string]any
	provisionVMAgent *bool
	patchError       bool
}

func (m *windowsPatchSettingsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host != "offline.invalid" {
		return nil, fmt.Errorf("unexpected host %q", req.URL.Host)
	}
	if req.URL.Query().Get("api-version") != "2024-03-01" {
		m.t.Errorf("unexpected API version: %s", req.URL)
	}
	status := http.StatusOK
	var body any
	switch req.Method {
	case http.MethodPut:
		var payload virtualmachines.VirtualMachine
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, err
		}
		m.puts = append(m.puts, payload)
		if payload.Properties.OsProfile != nil {
			encoded, _ := json.Marshal(payload.Properties.OsProfile.WindowsConfiguration.PatchSettings)
			if err := json.Unmarshal(encoded, &m.patchSettings); err != nil {
				return nil, err
			}
		}
	case http.MethodPatch:
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, err
		}
		m.patches = append(m.patches, payload)
		if m.patchError {
			status = http.StatusBadRequest
			body = map[string]any{"error": map[string]any{"code": "PatchRejected", "message": "offline rejection"}}
		} else {
			properties := payload["properties"].(map[string]any)
			profile := properties["osProfile"].(map[string]any)
			config := profile["windowsConfiguration"].(map[string]any)
			settings := config["patchSettings"].(map[string]any)
			for key, value := range settings {
				m.patchSettings[key] = value
			}
		}
	case http.MethodGet:
		if strings.HasSuffix(req.URL.Path, "/instanceView") {
			body = map[string]any{"statuses": []map[string]any{{"code": "PowerState/deallocated"}}}
		}
	default:
		return nil, fmt.Errorf("unexpected method/path: %s %s", req.Method, req.URL.Path)
	}
	if body == nil {
		config := map[string]any{"patchSettings": m.patchSettings}
		if m.provisionVMAgent != nil {
			config["provisionVMAgent"] = *m.provisionVMAgent
		}
		body = map[string]any{
			"location": "westus2",
			"properties": map[string]any{
				"provisioningState": "Succeeded",
				"osProfile":         map[string]any{"windowsConfiguration": config},
			},
		}
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set(sdkclient.SkipPollingDelayHeader, "true")
	return &http.Response{StatusCode: status, Header: headers, Body: io.NopCloser(strings.NewReader(string(encoded))), Request: req}, nil
}

func windowsPatchSettingsClient(t *testing.T, transport *windowsPatchSettingsTransport) *clients.Client {
	t.Helper()
	client, err := virtualmachines.NewVirtualMachinesClientWithBaseURI(environments.ResourceManagerAPI("https://offline.invalid"))
	if err != nil {
		t.Fatal(err)
	}
	client.Client.AuthorizeRequest = nil
	client.Client.DisableRetries = true
	client.Client.SetTransport(transport)
	return &clients.Client{
		Account:     &clients.ResourceManagerAccount{SubscriptionId: "00000000-0000-0000-0000-000000000000"},
		StopContext: context.Background(),
		Features:    features.UserFeatures{SkipImportCheckOnCreateAndAllowOverwritingExistingResources: true},
		Compute:     &computeclient.Client{VirtualMachinesClient: client},
		Network:     &networkclient.Client{Client: &networksdk.Client{}},
	}
}

func windowsPatchSettingsConfig() map[string]any {
	return map[string]any{
		"name": "testvm", "resource_group_name": "test", "location": "westus2", "size": "Standard_F2",
		"network_interface_ids": []any{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/networkInterfaces/test"},
		"os_managed_disk_id":    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Compute/disks/test",
		"os_disk":               []any{map[string]any{"caching": "ReadWrite"}},
	}
}

func windowsPatchSettingsApply(t *testing.T, resource *pluginsdk.Resource, old *terraform.InstanceState, config map[string]any, meta *clients.Client) (*terraform.InstanceState, error) {
	t.Helper()
	// Include every schema attribute with typed nulls, as Terraform's SDK bridge
	// does. NewResourceConfigRaw alone does not retain a protocol RawConfig.
	attributes := map[string]cty.Value{}
	for key, typ := range resource.CoreConfigSchema().ImpliedType().AttributeTypes() {
		attributes[key] = cty.NullVal(typ)
	}
	for key, value := range config {
		attributes[key] = windowsPatchSettingsConfigValue(value)
	}
	value, err := resource.CoreConfigSchema().CoerceValue(cty.ObjectVal(attributes))
	if err != nil {
		return old, err
	}
	raw := terraform.NewResourceConfigShimmed(value, resource.CoreConfigSchema())
	raw.CtyValue = value
	diff, err := resource.Diff(t.Context(), old, raw, meta)
	if err != nil {
		return old, err
	}
	if diff == nil {
		return old, nil
	}
	diff.RawConfig = value
	state, diagnostics := resource.Apply(t.Context(), old, diff, meta)
	if diagnostics.HasError() {
		return state, fmt.Errorf("%v", diagnostics)
	}
	return state, nil
}

func windowsPatchSettingsConfigValue(value any) cty.Value {
	switch value := value.(type) {
	case nil:
		return cty.NullVal(cty.DynamicPseudoType)
	case string:
		return cty.StringVal(value)
	case bool:
		return cty.BoolVal(value)
	case int:
		return cty.NumberIntVal(int64(value))
	case []any:
		values := make([]cty.Value, len(value))
		for index, item := range value {
			values[index] = windowsPatchSettingsConfigValue(item)
		}
		return cty.TupleVal(values)
	case map[string]any:
		values := map[string]cty.Value{}
		for key, item := range value {
			values[key] = windowsPatchSettingsConfigValue(item)
		}
		return cty.ObjectVal(values)
	default:
		panic(fmt.Sprintf("unsupported test config type %T", value))
	}
}

func TestWindowsVirtualMachineImportedDiskPatchSettingsCreate(t *testing.T) {
	cases := []struct {
		name   string
		fields map[string]any
		want   map[string]any
	}{
		{name: "unset"},
		{name: "explicit null", fields: map[string]any{"patch_mode": nil, "patch_assessment_mode": nil}},
		{name: "mode only", fields: map[string]any{"patch_mode": "AutomaticByOS"}, want: map[string]any{"patchMode": "AutomaticByOS"}},
		{name: "assessment only", fields: map[string]any{"patch_assessment_mode": "AutomaticByPlatform"}, want: map[string]any{"assessmentMode": "AutomaticByPlatform"}},
		{name: "both", fields: map[string]any{"patch_mode": "AutomaticByPlatform", "patch_assessment_mode": "AutomaticByPlatform"}, want: map[string]any{"patchMode": "AutomaticByPlatform", "assessmentMode": "AutomaticByPlatform"}},
		{name: "manual and image default", fields: map[string]any{"patch_mode": "Manual", "patch_assessment_mode": "ImageDefault"}, want: map[string]any{"patchMode": "Manual", "assessmentMode": "ImageDefault"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			transport := &windowsPatchSettingsTransport{t: t, patchSettings: map[string]any{}}
			config := windowsPatchSettingsConfig()
			for key, value := range tc.fields {
				config[key] = value
			}
			state, err := windowsPatchSettingsApply(t, resourceWindowsVirtualMachine(), nil, config, windowsPatchSettingsClient(t, transport))
			if err != nil {
				t.Fatal(err)
			}
			if len(transport.puts) != 1 {
				t.Fatalf("expected one PUT, got %d", len(transport.puts))
			}
			if transport.puts[0].Properties.OsProfile != nil {
				t.Fatal("attached disk PUT must omit osProfile")
			}
			if transport.puts[0].Properties.StorageProfile.OsDisk.CreateOption != virtualmachines.DiskCreateOptionTypesAttach {
				t.Fatal("expected Attach create option")
			}
			if tc.want == nil {
				if len(transport.patches) != 0 {
					t.Fatalf("unset patch fields caused PATCH: %v", transport.patches)
				}
			} else {
				want := map[string]any{"properties": map[string]any{"osProfile": map[string]any{"windowsConfiguration": map[string]any{"patchSettings": tc.want}}}}
				if len(transport.patches) != 1 || !reflect.DeepEqual(transport.patches[0], want) {
					t.Fatalf("expected PATCH %v, got %v", want, transport.patches)
				}
				for key, wantValue := range tc.fields {
					if got := state.Attributes[key]; got != wantValue {
						t.Errorf("read %s: got %v, want %v", key, got, wantValue)
					}
				}
			}
			if state.ID == "" || state.Identity["name"] != "testvm" {
				t.Fatalf("missing resource ID/identity: %#v", state)
			}
		})
	}
}

func TestWindowsVirtualMachineImportedDiskPatchSettingsFailurePreservesIdentity(t *testing.T) {
	transport := &windowsPatchSettingsTransport{t: t, patchSettings: map[string]any{}, patchError: true}
	config := windowsPatchSettingsConfig()
	config["patch_mode"] = "AutomaticByPlatform"
	state, err := windowsPatchSettingsApply(t, resourceWindowsVirtualMachine(), nil, config, windowsPatchSettingsClient(t, transport))
	if err == nil || !strings.Contains(err.Error(), "updating patch settings for Windows") {
		t.Fatalf("expected patch failure, got %v", err)
	}
	if state == nil || state.ID == "" || state.Identity["name"] != "testvm" {
		t.Fatalf("created VM identity lost on PATCH failure: %#v", state)
	}
}

func TestWindowsVirtualMachineImportedDiskPatchSettingsUpdate(t *testing.T) {
	for _, agent := range []*bool{nil, pointer.To(true), pointer.To(false)} {
		name := "agent omitted"
		if agent != nil {
			name = fmt.Sprintf("agent %t", *agent)
		}
		t.Run(name, func(t *testing.T) {
			transport := &windowsPatchSettingsTransport{t: t, patchSettings: map[string]any{}, provisionVMAgent: agent}
			meta := windowsPatchSettingsClient(t, transport)
			resource := resourceWindowsVirtualMachine()
			config := windowsPatchSettingsConfig()
			config["patch_mode"], config["patch_assessment_mode"] = "AutomaticByOS", "ImageDefault"
			old, err := windowsPatchSettingsApply(t, resource, nil, config, meta)
			if err != nil {
				t.Fatal(err)
			}
			wantAgent := "true"
			if agent != nil && !*agent {
				wantAgent = "false"
			}
			if old.Attributes["provision_vm_agent"] != wantAgent {
				t.Fatalf("agent state: got %q, want %q", old.Attributes["provision_vm_agent"], wantAgent)
			}
			transport.patches = nil
			config["patch_mode"], config["patch_assessment_mode"] = "AutomaticByPlatform", "AutomaticByPlatform"
			state, err := windowsPatchSettingsApply(t, resource, old, config, meta)
			if agent != nil && !*agent {
				if err == nil || !strings.Contains(err.Error(), "provision_vm_agent") {
					t.Fatalf("known-disabled agent should be rejected, got %v", err)
				}
				if len(transport.patches) != 0 {
					t.Fatal("known-disabled agent caused PATCH")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(transport.patches) != 1 {
				t.Fatalf("expected one update PATCH, got %v", transport.patches)
			}
			if state.Attributes["patch_mode"] != "AutomaticByPlatform" || state.Attributes["patch_assessment_mode"] != "AutomaticByPlatform" {
				t.Fatalf("patch fields did not roundtrip: %v", state.Attributes)
			}
			transport.patches = nil
			if _, err := windowsPatchSettingsApply(t, resource, state, config, meta); err != nil {
				t.Fatal(err)
			}
			if len(transport.patches) != 0 {
				t.Fatal("unchanged patch fields caused PATCH")
			}
		})
	}
}

func TestWindowsVirtualMachineImportedDiskInheritsPatchSettings(t *testing.T) {
	transport := &windowsPatchSettingsTransport{
		t: t, provisionVMAgent: pointer.To(true),
		patchSettings: map[string]any{"patchMode": "AutomaticByOS", "assessmentMode": "ImageDefault"},
	}
	resource := resourceWindowsVirtualMachine()
	meta := windowsPatchSettingsClient(t, transport)
	config := windowsPatchSettingsConfig()
	state, err := windowsPatchSettingsApply(t, resource, nil, config, meta)
	if err != nil {
		t.Fatal(err)
	}
	if len(transport.patches) != 0 {
		t.Fatalf("omitted fields overwrote inherited patch settings: %v", transport.patches)
	}
	if state.Attributes["patch_mode"] != "AutomaticByOS" || state.Attributes["patch_assessment_mode"] != "ImageDefault" {
		t.Fatalf("inherited settings not read: %v", state.Attributes)
	}
	if _, err := windowsPatchSettingsApply(t, resource, state, config, meta); err != nil {
		t.Fatal(err)
	}
	if len(transport.patches) != 0 {
		t.Fatal("computed patch fields caused an update")
	}
}

func TestWindowsVirtualMachineImagePatchSettingsCreate(t *testing.T) {
	transport := &windowsPatchSettingsTransport{t: t, patchSettings: map[string]any{}}
	config := windowsPatchSettingsConfig()
	delete(config, "os_managed_disk_id")
	config["admin_username"] = "adminuser"
	config["admin_password"] = "P@ssword12345!"
	config["os_disk"] = []any{map[string]any{"caching": "ReadWrite", "storage_account_type": "Standard_LRS"}}
	config["source_image_reference"] = []any{map[string]any{"publisher": "MicrosoftWindowsServer", "offer": "WindowsServer", "sku": "2016-Datacenter", "version": "latest"}}
	config["patch_mode"], config["patch_assessment_mode"] = "AutomaticByPlatform", "AutomaticByPlatform"
	state, err := windowsPatchSettingsApply(t, resourceWindowsVirtualMachine(), nil, config, windowsPatchSettingsClient(t, transport))
	if err != nil {
		t.Fatal(err)
	}
	if len(transport.puts) != 1 || transport.puts[0].Properties.OsProfile == nil {
		t.Fatal("image VM must include osProfile in its PUT")
	}
	if len(transport.patches) != 0 {
		t.Fatalf("image VM got an extra PATCH: %v", transport.patches)
	}
	if state.Attributes["patch_mode"] != "AutomaticByPlatform" || state.Attributes["patch_assessment_mode"] != "AutomaticByPlatform" {
		t.Fatalf("image patch settings did not roundtrip: %v", state.Attributes)
	}
}
