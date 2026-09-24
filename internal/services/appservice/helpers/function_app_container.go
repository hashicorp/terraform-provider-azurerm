// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-cty/cty"
)

func ValidateFunctionAppContainerName(name string) error {
	if len(name) > 32 || !regexp.MustCompile(`^[a-z]([a-z0-9-]*[a-z0-9])?$`).MatchString(name) || strings.Contains(name, "--") {
		return fmt.Errorf("`name` must be at most 32 characters, start with a lowercase letter, end with a lowercase letter or digit, and contain only lowercase letters, digits and single hyphens when `container_app_environment_id` is configured")
	}
	return nil
}

// Container Apps uses a subset of SiteConfig. In particular, App Service defaults
// must not be sent on either create or a read-modify-write update.
func FunctionAppContainerSiteConfig(input *webapps.SiteConfig) *webapps.SiteConfig {
	if input == nil {
		return nil
	}
	return &webapps.SiteConfig{
		LinuxFxVersion:              input.LinuxFxVersion,
		AppSettings:                 input.AppSettings,
		AcrUseManagedIdentityCreds:  input.AcrUseManagedIdentityCreds,
		AcrUserManagedIdentityID:    input.AcrUserManagedIdentityID,
		MinimumElasticInstanceCount: input.MinimumElasticInstanceCount,
		FunctionAppScaleLimit:       input.FunctionAppScaleLimit,
	}
}

func FunctionAppContainerSiteProperties(input *webapps.SiteProperties) *webapps.SiteProperties {
	if input == nil {
		return nil
	}
	properties := &webapps.SiteProperties{
		ManagedEnvironmentId:      input.ManagedEnvironmentId,
		SiteConfig:                FunctionAppContainerSiteConfig(input.SiteConfig),
		KeyVaultReferenceIdentity: input.KeyVaultReferenceIdentity,
		DaprConfig:                input.DaprConfig,
	}
	// Consumption-only environments return resource defaults but reject them on write.
	if input.WorkloadProfileName != nil && *input.WorkloadProfileName != "" {
		properties.WorkloadProfileName = input.WorkloadProfileName
		properties.ResourceConfig = input.ResourceConfig
	}
	return properties
}

func ValidateFunctionAppContainerSiteConfig(raw cty.Value) error {
	if raw.IsNull() || !raw.IsKnown() {
		return nil
	}
	fields := raw.AsValueMap()
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		value := fields[name]
		if value.IsNull() {
			continue
		}
		if value.IsKnown() && (value.Type().IsListType() || value.Type().IsSetType() || value.Type().IsTupleType()) && value.Length().IsKnown() && value.LengthInt() == 0 {
			continue
		}
		switch name {
		case "application_stack", "application_insights_key", "application_insights_connection_string",
			"container_registry_use_managed_identity", "container_registry_managed_identity_client_id",
			"elastic_instance_minimum", "app_scale_limit":
			continue
		}
		return fmt.Errorf("`site_config.%s` is not supported by this resource when `container_app_environment_id` is configured", name)
	}
	return nil
}

// These schema defaults are placeholders for settings that do not apply to
// Container Apps. Keep them stable in resource state without sending them to Azure.
func SetFunctionAppContainerSiteConfigDefaults(config *SiteConfigLinuxFunctionApp) {
	config.ManagedPipelineMode = string(webapps.ManagedPipelineModeIntegrated)
	config.LoadBalancing = string(webapps.SiteLoadBalancingLeastRequests)
	config.FtpsState = string(webapps.FtpsStateDisabled)
	config.MinTlsVersion = string(webapps.SupportedTlsVersionsOnePointTwo)
	config.ScmMinTlsVersion = string(webapps.SupportedTlsVersionsOnePointTwo)
	config.IpRestrictionDefaultAction = string(webapps.DefaultActionAllow)
	config.ScmIpRestrictionDefaultAction = string(webapps.DefaultActionAllow)
}
