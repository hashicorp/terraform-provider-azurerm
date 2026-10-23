// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func TestValidateFunctionAppContainerName(t *testing.T) {
	for _, name := range []string{"a", "app-123", "acctest-ca-260916163258940770", strings.Repeat("a", 32)} {
		if err := helpers.ValidateFunctionAppContainerName(name); err != nil {
			t.Errorf("valid name %q rejected: %v", name, err)
		}
	}
	for _, name := range []string{"", "App", "1app", "app-", "app--name", "app.name", strings.Repeat("a", 33)} {
		if err := helpers.ValidateFunctionAppContainerName(name); err == nil {
			t.Errorf("invalid name %q accepted", name)
		}
	}
}

func TestFunctionAppContainerSiteConfig(t *testing.T) {
	input := &webapps.SiteConfig{
		LinuxFxVersion:              pointer.To("DOCKER|mcr.microsoft.com/azure-functions/dotnet:4-dotnet8"),
		AppSettings:                 pointer.To([]webapps.NameValuePair{{Name: pointer.To("example"), Value: pointer.To("value")}}),
		AcrUseManagedIdentityCreds:  pointer.To(true),
		AcrUserManagedIdentityID:    pointer.To("identity-client-id"),
		MinimumElasticInstanceCount: pointer.To(int64(1)),
		FunctionAppScaleLimit:       pointer.To(int64(10)),
		PreWarmedInstanceCount:      pointer.To(int64(0)),
		ManagedPipelineMode:         pointer.To(webapps.ManagedPipelineModeIntegrated),
		LoadBalancing:               pointer.To(webapps.SiteLoadBalancingLeastRequests),
		PublicNetworkAccess:         pointer.To("Enabled"),
		Use32BitWorkerProcess:       pointer.To(false),
	}
	actual := helpers.FunctionAppContainerSiteConfig(input)
	encoded, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 6 {
		t.Fatalf("expected only the six supported fields, got %s", encoded)
	}
	for _, name := range []string{"linuxFxVersion", "appSettings", "acrUseManagedIdentityCreds", "acrUserManagedIdentityID", "minimumElasticInstanceCount", "functionAppScaleLimit"} {
		if _, ok := fields[name]; !ok {
			t.Errorf("missing %s", name)
		}
	}
	if !reflect.DeepEqual(input.AppSettings, actual.AppSettings) || !reflect.DeepEqual(input.LinuxFxVersion, actual.LinuxFxVersion) {
		t.Fatal("image or app settings were changed")
	}
	if input.PreWarmedInstanceCount == nil {
		t.Fatal("input was mutated")
	}
	if helpers.FunctionAppContainerSiteConfig(nil) != nil {
		t.Fatal("expected nil for a nil config")
	}
}

func TestFunctionAppContainerSiteProperties(t *testing.T) {
	input := &webapps.SiteProperties{
		ManagedEnvironmentId: pointer.To("/subscriptions/example/managedEnvironments/test"),
		Enabled:              pointer.To(true),
		SiteConfig: &webapps.SiteConfig{
			LinuxFxVersion: pointer.To("DOCKER|example/image:latest"),
		},
		PublicNetworkAccess:       pointer.To("Enabled"),
		VnetImagePullEnabled:      pointer.To(false),
		VnetBackupRestoreEnabled:  pointer.To(false),
		DailyMemoryTimeQuota:      pointer.To(int64(0)),
		ClientCertMode:            pointer.To(webapps.ClientCertModeOptional),
		ClientCertEnabled:         pointer.To(false),
		HTTPSOnly:                 pointer.To(false),
		DefaultHostName:           pointer.To("read-only.example.com"),
		HostNames:                 pointer.To([]string{"read-only.example.com"}),
		State:                     pointer.To("Running"),
		KeyVaultReferenceIdentity: pointer.To("SystemAssigned"),
		DaprConfig:                &webapps.DaprConfig{Enabled: pointer.To(true), AppId: pointer.To("example")},
		ResourceConfig:            &webapps.ResourceConfig{Cpu: pointer.To(0.5), Memory: pointer.To("1Gi")},
		WorkloadProfileName:       pointer.To("Consumption"),
	}
	actual := helpers.FunctionAppContainerSiteProperties(input)
	if actual.ManagedEnvironmentId != input.ManagedEnvironmentId || actual.SiteConfig.LinuxFxVersion != input.SiteConfig.LinuxFxVersion {
		t.Fatal("supported site properties were changed")
	}
	if !reflect.DeepEqual(actual.DaprConfig, input.DaprConfig) || !reflect.DeepEqual(actual.ResourceConfig, input.ResourceConfig) ||
		!reflect.DeepEqual(actual.WorkloadProfileName, input.WorkloadProfileName) || !reflect.DeepEqual(actual.KeyVaultReferenceIdentity, input.KeyVaultReferenceIdentity) {
		t.Fatal("existing Container Apps configuration was changed")
	}
	encoded, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 6 {
		t.Fatalf("unexpected properties in payload: %s", encoded)
	}
	if input.DailyMemoryTimeQuota == nil || input.SiteConfig == actual.SiteConfig {
		t.Fatal("input was mutated")
	}
	if helpers.FunctionAppContainerSiteProperties(nil) != nil {
		t.Fatal("expected nil for nil properties")
	}
	for _, profile := range []*string{nil, pointer.To("")} {
		input.WorkloadProfileName = profile
		actual := helpers.FunctionAppContainerSiteProperties(input)
		if actual.ResourceConfig != nil || actual.WorkloadProfileName != nil {
			t.Fatal("resource configuration must not be sent without a workload profile")
		}
		if input.ResourceConfig == nil {
			t.Fatal("input resource configuration was mutated")
		}
	}
}

func TestValidateFunctionAppContainerSiteConfig(t *testing.T) {
	for _, test := range []struct {
		name    string
		value   cty.Value
		wantErr bool
	}{
		{"null", cty.NullVal(cty.DynamicPseudoType), false},
		{"unknown", cty.DynamicVal, false},
		{"omitted", cty.ObjectVal(map[string]cty.Value{"managed_pipeline_mode": cty.NullVal(cty.String)}), false},
		{"empty blocks", cty.ObjectVal(map[string]cty.Value{"app_service_logs": cty.ListValEmpty(cty.EmptyObject)}), false},
		{"supported unknown", cty.ObjectVal(map[string]cty.Value{"application_stack": cty.DynamicVal}), false},
		{"replicas", cty.ObjectVal(map[string]cty.Value{"elastic_instance_minimum": cty.NumberIntVal(0), "app_scale_limit": cty.NumberIntVal(10)}), false},
		{"explicit zero", cty.ObjectVal(map[string]cty.Value{"pre_warmed_instance_count": cty.NumberIntVal(0)}), true},
		{"explicit default", cty.ObjectVal(map[string]cty.Value{"managed_pipeline_mode": cty.StringVal("Integrated")}), true},
		{"explicit false", cty.ObjectVal(map[string]cty.Value{"use_32_bit_worker": cty.False}), true},
		{"unsupported unknown", cty.ObjectVal(map[string]cty.Value{"load_balancing_mode": cty.UnknownVal(cty.String)}), true},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := helpers.ValidateFunctionAppContainerSiteConfig(test.value)
			if (err != nil) != test.wantErr {
				t.Fatalf("want error %t, got %v", test.wantErr, err)
			}
		})
	}
}

func TestFunctionAppContainerSiteConfigDefaults(t *testing.T) {
	config := helpers.SiteConfigLinuxFunctionApp{AppScaleLimit: 10, ElasticInstanceMinimum: 1}
	helpers.SetFunctionAppContainerSiteConfigDefaults(&config)
	if config.ManagedPipelineMode != "Integrated" || config.LoadBalancing != "LeastRequests" ||
		config.FtpsState != "Disabled" || config.MinTlsVersion != "1.2" || config.ScmMinTlsVersion != "1.2" ||
		config.IpRestrictionDefaultAction != "Allow" || config.ScmIpRestrictionDefaultAction != "Allow" {
		t.Fatalf("unexpected placeholder defaults: %+v", config)
	}
	if config.AppScaleLimit != 10 || config.ElasticInstanceMinimum != 1 {
		t.Fatal("supported replica settings were changed")
	}
}
