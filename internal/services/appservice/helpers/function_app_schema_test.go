// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers_test

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func TestMergeUserAppSettings(t *testing.T) {
	cases := []struct {
		service  []webapps.NameValuePair
		user     map[string]string
		expected []webapps.NameValuePair // Note: The function doesn't preserve the order of the list
	}{
		{
			service: []webapps.NameValuePair{
				{
					Name:  pointer.To("test"),
					Value: pointer.To("ServiceValue"),
				},
			},
			user: map[string]string{
				"test": "UserValue",
			},
			expected: []webapps.NameValuePair{
				{
					Name:  pointer.To("test"),
					Value: pointer.To("UserValue"),
				},
			},
		},
		{
			service: []webapps.NameValuePair{
				{
					Name:  pointer.To("test"),
					Value: pointer.To("ServiceValue"),
				},
				{
					Name:  pointer.To("test2"),
					Value: pointer.To("ServiceValue2"),
				},
				{
					Name:  pointer.To("test3"),
					Value: pointer.To("ServiceValue3"),
				},
				{
					Name:  pointer.To("test4"),
					Value: pointer.To("ServiceValue4"),
				},
			},
			user: map[string]string{
				"test":  "UserValue",
				"test4": "UserValue4",
			},
			expected: []webapps.NameValuePair{
				{
					Name:  pointer.To("test"),
					Value: pointer.To("UserValue"),
				},
				{
					Name:  pointer.To("test2"),
					Value: pointer.To("ServiceValue2"),
				},
				{
					Name:  pointer.To("test3"),
					Value: pointer.To("ServiceValue3"),
				},
				{
					Name:  pointer.To("test4"),
					Value: pointer.To("UserValue4"),
				},
			},
		},
	}

	for _, v := range cases {
		actualRaw := helpers.MergeUserAppSettings(&v.service, v.user)
		actual := *actualRaw
		sort.Slice(actual, func(i, j int) bool {
			return *actual[i].Name < *actual[j].Name
		})
		if !reflect.DeepEqual(actual, v.expected) {
			t.Fatalf("expected %+v, got %+v", v.expected, actual)
		}
	}
}

func windowsFunctionAppSiteConfigRaw(stack map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"site_config": []interface{}{
			map[string]interface{}{
				"application_stack": []interface{}{stack},
			},
		},
	}
}

func appSettingValue(settings *[]webapps.NameValuePair, name string) string {
	if settings == nil {
		return ""
	}
	for _, v := range *settings {
		if pointer.From(v.Name) == name {
			return pointer.From(v.Value)
		}
	}
	return ""
}

func TestExpandSiteConfigWindowsFunctionApp_PowerShellStackSettings(t *testing.T) {
	cases := []struct {
		name                  string
		stack                 helpers.ApplicationStackWindowsFunctionApp
		existing              *webapps.SiteConfig
		expectedNetFramework  string
		expectedPowerShell    string
		expectedWorkerRuntime string
		expected32BitWorker   bool
	}{
		{
			name:                  "create powershell 7.6 uses the .NET 10 host",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v4.0", PowerShellCoreVersion: "7.6"},
			expectedNetFramework:  "v10.0",
			expectedPowerShell:    "7.6",
			expectedWorkerRuntime: "powershell",
		},
		{
			name:                  "existing powershell 7.6 app on a v4.0 host is corrected",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v4.0", PowerShellCoreVersion: "7.6"},
			existing:              &webapps.SiteConfig{NetFrameworkVersion: pointer.To("v4.0"), PowerShellVersion: pointer.To("7.6"), Use32BitWorkerProcess: pointer.To(true)},
			expectedNetFramework:  "v10.0",
			expectedPowerShell:    "7.6",
			expectedWorkerRuntime: "powershell",
		},
		{
			name:                  "existing powershell 7.6 app on the correct host keeps it",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v4.0", PowerShellCoreVersion: "7.6"},
			existing:              &webapps.SiteConfig{NetFrameworkVersion: pointer.To("v10.0"), PowerShellVersion: pointer.To("7.6")},
			expectedNetFramework:  "v10.0",
			expectedPowerShell:    "7.6",
			expectedWorkerRuntime: "powershell",
		},
		{
			name:                  "upgrade from powershell 7.4 to 7.6",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v4.0", PowerShellCoreVersion: "7.6"},
			existing:              &webapps.SiteConfig{NetFrameworkVersion: pointer.To("v4.0"), PowerShellVersion: pointer.To("7.4")},
			expectedNetFramework:  "v10.0",
			expectedPowerShell:    "7.6",
			expectedWorkerRuntime: "powershell",
		},
		{
			name:                  "powershell 7.4 is unchanged",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v4.0", PowerShellCoreVersion: "7.4"},
			expectedNetFramework:  "v4.0",
			expectedPowerShell:    "7.4",
			expectedWorkerRuntime: "powershell",
			expected32BitWorker:   true,
		},
		{
			name:                  "dotnet is unchanged",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v8.0"},
			expectedNetFramework:  "v8.0",
			expectedWorkerRuntime: "dotnet",
			expected32BitWorker:   true,
		},
		{
			name:                  "node is unchanged",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v4.0", NodeVersion: "~22"},
			expectedNetFramework:  "v4.0",
			expectedWorkerRuntime: "node",
			expected32BitWorker:   true,
		},
		{
			name:                  "powershell 7 is unchanged",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v4.0", PowerShellCoreVersion: "7"},
			expectedNetFramework:  "v4.0",
			expectedPowerShell:    "7",
			expectedWorkerRuntime: "powershell",
			expected32BitWorker:   true,
		},
		{
			name:                  "powershell 7.2 is unchanged",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v4.0", PowerShellCoreVersion: "7.2"},
			expectedNetFramework:  "v4.0",
			expectedPowerShell:    "7.2",
			expectedWorkerRuntime: "powershell",
			expected32BitWorker:   true,
		},
		{
			name:                  "java is unchanged",
			stack:                 helpers.ApplicationStackWindowsFunctionApp{DotNetVersion: "v4.0", JavaVersion: "17"},
			expectedNetFramework:  "v4.0",
			expectedWorkerRuntime: "java",
			expected32BitWorker:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stackRaw := map[string]interface{}{
				"dotnet_version":          tc.stack.DotNetVersion,
				"node_version":            tc.stack.NodeVersion,
				"powershell_core_version": tc.stack.PowerShellCoreVersion,
			}

			appData := schema.TestResourceDataRaw(t, map[string]*schema.Schema{"site_config": helpers.SiteConfigSchemaWindowsFunctionApp()}, windowsFunctionAppSiteConfigRaw(stackRaw))
			appExisting := copySiteConfig(tc.existing)
			app, err := helpers.ExpandSiteConfigWindowsFunctionApp([]helpers.SiteConfigWindowsFunctionApp{{ApplicationStack: []helpers.ApplicationStackWindowsFunctionApp{tc.stack}, Use32BitWorker: true}}, appExisting, sdk.ResourceMetaData{ResourceData: appData}, "~4", "storage", false)
			if err != nil {
				t.Fatalf("expanding app site config: %+v", err)
			}

			slotData := schema.TestResourceDataRaw(t, map[string]*schema.Schema{"site_config": helpers.SiteConfigSchemaWindowsFunctionAppSlot()}, windowsFunctionAppSiteConfigRaw(stackRaw))
			slotExisting := copySiteConfig(tc.existing)
			slot, err := helpers.ExpandSiteConfigWindowsFunctionAppSlot([]helpers.SiteConfigWindowsFunctionAppSlot{{ApplicationStack: []helpers.ApplicationStackWindowsFunctionApp{tc.stack}, Use32BitWorker: true}}, slotExisting, sdk.ResourceMetaData{ResourceData: slotData}, "~4", "storage", false)
			if err != nil {
				t.Fatalf("expanding slot site config: %+v", err)
			}

			for kind, actual := range map[string]*webapps.SiteConfig{"app": app, "slot": slot} {
				if actual.Use32BitWorkerProcess == nil || *actual.Use32BitWorkerProcess != tc.expected32BitWorker {
					t.Errorf("%s: expected use32BitWorkerProcess %t, got %+v", kind, tc.expected32BitWorker, actual.Use32BitWorkerProcess)
				}
				if tc.existing != nil && ((kind == "app" && actual != appExisting) || (kind == "slot" && actual != slotExisting)) {
					t.Errorf("%s: expected the existing site config to be updated in place", kind)
				}
				if v := pointer.From(actual.NetFrameworkVersion); v != tc.expectedNetFramework {
					t.Errorf("%s: expected netFrameworkVersion %q, got %q", kind, tc.expectedNetFramework, v)
				}
				if v := pointer.From(actual.PowerShellVersion); v != tc.expectedPowerShell {
					t.Errorf("%s: expected powerShellVersion %q, got %q", kind, tc.expectedPowerShell, v)
				}
				if v := appSettingValue(actual.AppSettings, "FUNCTIONS_WORKER_RUNTIME"); v != tc.expectedWorkerRuntime {
					t.Errorf("%s: expected FUNCTIONS_WORKER_RUNTIME %q, got %q", kind, tc.expectedWorkerRuntime, v)
				}
				if v := appSettingValue(actual.AppSettings, "FUNCTIONS_EXTENSION_VERSION"); v != "~4" {
					t.Errorf("%s: expected FUNCTIONS_EXTENSION_VERSION ~4, got %q", kind, v)
				}
				for _, setting := range *actual.AppSettings {
					if pointer.From(setting.Name) == "FUNCTIONS_WORKER_RUNTIME_VERSION" {
						t.Errorf("%s: Windows worker version must use siteConfig.powerShellVersion, not FUNCTIONS_WORKER_RUNTIME_VERSION", kind)
					}
				}
			}
		})
	}
}

func copySiteConfig(input *webapps.SiteConfig) *webapps.SiteConfig {
	if input == nil {
		return nil
	}
	output := *input
	return &output
}

func TestFlattenSiteConfigWindowsFunctionApp_PowerShellStackSettings(t *testing.T) {
	cases := []struct {
		name           string
		powershell     string
		host           string
		worker32Bit    bool
		expectedDotNet string
	}{
		{name: "powershell 7.6", powershell: "7.6", host: "v10.0", expectedDotNet: "v4.0"},
		{name: "legacy powershell 7.6", powershell: "7.6", host: "v4.0", worker32Bit: true, expectedDotNet: "v4.0"},
		{name: "powershell 7.6 host drift", powershell: "7.6", host: "v8.0", expectedDotNet: "v4.0"},
		{name: "powershell 7.4 unchanged", powershell: "7.4", host: "v8.0", worker32Bit: true, expectedDotNet: "v8.0"},
		{name: "dotnet unchanged", host: "v10.0", worker32Bit: true, expectedDotNet: "v10.0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := &webapps.SiteConfig{
				NetFrameworkVersion:   pointer.To(tc.host),
				PowerShellVersion:     pointer.To(tc.powershell),
				Use32BitWorkerProcess: pointer.To(tc.worker32Bit),
			}
			app, err := helpers.FlattenSiteConfigWindowsFunctionApp(input)
			if err != nil {
				t.Fatalf("flattening app site config: %+v", err)
			}
			slot, err := helpers.FlattenSiteConfigWindowsFunctionAppSlot(input)
			if err != nil {
				t.Fatalf("flattening slot site config: %+v", err)
			}
			for kind, stack := range map[string]helpers.ApplicationStackWindowsFunctionApp{"app": app.ApplicationStack[0], "slot": slot.ApplicationStack[0]} {
				if stack.DotNetVersion != tc.expectedDotNet || stack.PowerShellCoreVersion != tc.powershell {
					t.Errorf("%s: unexpected application stack: %+v", kind, stack)
				}
			}
			if app.Use32BitWorker != tc.worker32Bit || slot.Use32BitWorker != tc.worker32Bit {
				t.Fatal("flatten must report the actual worker bitness")
			}
		})
	}
}

func TestWindowsFunctionAppStack_PowerShell76Diff(t *testing.T) {
	const dotNetVersionKey = "site_config.0.application_stack.0.dotnet_version"
	const workerKey = "site_config.0.use_32_bit_worker"

	cases := []struct {
		name         string
		stateDotNet  string
		statePS      string
		stateWorker  string
		configStack  map[string]interface{}
		configWorker interface{}
		dotNetDiff   bool
		workerDiff   bool
	}{
		{
			name:        "refresh/import of a powershell 7.6 app is stable",
			stateDotNet: "v4.0",
			statePS:     "7.6",
			stateWorker: "false",
			configStack: map[string]interface{}{"powershell_core_version": "7.6"},
		},
		{
			name:        "legacy powershell 7.6 app on v4.0 and 32-bit has no diff",
			stateDotNet: "v4.0",
			statePS:     "7.6",
			stateWorker: "true",
			configStack: map[string]interface{}{"powershell_core_version": "7.6"},
		},
		{
			name:         "explicit 64-bit is stable",
			stateDotNet:  "v4.0",
			statePS:      "7.6",
			stateWorker:  "false",
			configStack:  map[string]interface{}{"powershell_core_version": "7.6"},
			configWorker: false,
		},
		{
			name:         "explicit 64-bit corrects a legacy 32-bit worker",
			stateDotNet:  "v4.0",
			statePS:      "7.6",
			stateWorker:  "true",
			configStack:  map[string]interface{}{"powershell_core_version": "7.6"},
			configWorker: false,
			workerDiff:   true,
		},
		{
			name:        "downgrade to powershell 7.4 restores the default worker",
			stateDotNet: "v4.0",
			statePS:     "7.6",
			stateWorker: "false",
			configStack: map[string]interface{}{"powershell_core_version": "7.4"},
			workerDiff:  true,
		},
		{
			name:        "powershell 7.4 diffs are not suppressed",
			stateDotNet: "v8.0",
			statePS:     "7.4",
			stateWorker: "false",
			configStack: map[string]interface{}{"powershell_core_version": "7.4"},
			dotNetDiff:  true,
			workerDiff:  true,
		},
		{
			name:        "dotnet diffs are not suppressed",
			stateDotNet: "v10.0",
			stateWorker: "false",
			configStack: map[string]interface{}{"dotnet_version": "v4.0"},
			dotNetDiff:  true,
			workerDiff:  true,
		},
	}

	resources := map[string]*schema.Resource{
		"app":  {Schema: map[string]*schema.Schema{"site_config": helpers.SiteConfigSchemaWindowsFunctionApp()}},
		"slot": {Schema: map[string]*schema.Schema{"site_config": helpers.SiteConfigSchemaWindowsFunctionAppSlot()}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for kind, resource := range resources {
				stateRaw := windowsFunctionAppSiteConfigRaw(map[string]interface{}{
					"dotnet_version":          tc.stateDotNet,
					"powershell_core_version": tc.statePS,
				})
				stateRaw["site_config"].([]interface{})[0].(map[string]interface{})["use_32_bit_worker"] = tc.stateWorker == "true"
				stateData := schema.TestResourceDataRaw(t, resource.Schema, stateRaw)
				stateData.SetId("test")
				state := stateData.State()
				raw := windowsFunctionAppSiteConfigRaw(tc.configStack)
				if tc.configWorker != nil {
					raw["site_config"].([]interface{})[0].(map[string]interface{})["use_32_bit_worker"] = tc.configWorker
				}
				config := terraform.NewResourceConfigRaw(raw)

				diff, err := resource.Diff(context.Background(), state, config, nil)
				if err != nil {
					t.Fatalf("%s: diffing: %+v", kind, err)
				}

				if tc.statePS == "7.6" && tc.configStack["powershell_core_version"] == "7.6" && !tc.workerDiff && !tc.dotNetDiff && diff != nil && !diff.Empty() {
					t.Errorf("%s: expected a clean follow-up plan, got %+v", kind, diff)
				}
				for key, expected := range map[string]bool{dotNetVersionKey: tc.dotNetDiff, workerKey: tc.workerDiff} {
					var attrDiff *terraform.ResourceAttrDiff
					if diff != nil {
						attrDiff = diff.Attributes[key]
					}
					if hasDiff := attrDiff != nil; hasDiff != expected {
						t.Errorf("%s: expected %s diff %t, got %t (%+v)", kind, key, expected, hasDiff, attrDiff)
					}
				}
			}
		})
	}
}
