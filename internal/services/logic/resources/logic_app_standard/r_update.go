// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_standard

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := commonids.ParseLogicAppId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			data := LogicAppResourceModel{}

			if err := metadata.Decode(&data); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", *id, err)
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s for update", id)
			}

			siteEnvelope := *existing.Model.Properties

			sc, err := client.GetConfiguration(ctx, *id)
			if err != nil || sc.Model == nil {
				return fmt.Errorf("reading site_config for %s: %v", *id, err)
			}

			existingSiteConfig := sc.Model.Properties
			siteEnvelope.SiteConfig = existingSiteConfig

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil || appSettingsResp.Model == nil {
				return fmt.Errorf("reading App Settings for Linux %s: %+v", id, err)
			}

			currentAppSettings := make([]webapps.NameValuePair, 0)
			if appSettingsResp.Model.Properties != nil {
				currentAppSettings = expandAppSettings(*appSettingsResp.Model.Properties)
			}
			existingSiteConfig.AppSettings = pointer.To(currentAppSettings)

			// lintignore:R019 // deliberate subset: only the fields that require rebuilding site_config; the remaining attributes are applied elsewhere in the update
			if metadata.ResourceData.HasChanges(
				"site_config",
				"app_settings",
				"version",
				"storage_account_name",
				"storage_account_access_key",
				"storage_key_vault_secret_id",
			) {
				existingSiteConfig, err = expandLogicAppStandardSiteConfigForUpdate(data.SiteConfig, metadata, existingSiteConfig)
				if err != nil {
					return fmt.Errorf("expanding site_config update for %s: %v", *id, err)
				}

				siteEnvelope.SiteConfig = existingSiteConfig
			}

			if metadata.ResourceData.HasChange("site_config.0.linux_fx_version") {
				kind := "functionapp,workflowapp"
				if metadata.ResourceData.Get("site_config.0.linux_fx_version").(string) != "" {
					kind = "functionapp,linux,container,workflowapp"
				}
				existing.Model.Kind = pointer.To(kind)
			}

			if metadata.ResourceData.HasChange("app_service_plan_id") {
				siteEnvelope.ServerFarmId = pointer.To(data.AppServicePlanId)
			}

			if metadata.ResourceData.HasChange("enabled") {
				siteEnvelope.Enabled = pointer.To(data.Enabled)
			}

			if metadata.ResourceData.HasChange("client_affinity_enabled") {
				siteEnvelope.ClientAffinityEnabled = pointer.To(data.ClientAffinityEnabled)
			}

			if metadata.ResourceData.HasChange("client_certificate_mode") {
				siteEnvelope.ClientCertMode = pointer.ToEnum[webapps.ClientCertMode](data.ClientCertificateMode)
				siteEnvelope.ClientCertEnabled = pointer.To(data.ClientCertificateMode != "")
			}

			if metadata.ResourceData.HasChange("https_only") {
				siteEnvelope.HTTPSOnly = pointer.To(data.HTTPSOnly)
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				siteEnvelope.PublicNetworkAccess = pointer.To(data.PublicNetworkAccess)
				// Both fields must agree so the existing site configuration does not override the update.
				siteEnvelope.SiteConfig.PublicNetworkAccess = siteEnvelope.PublicNetworkAccess
			}

			if metadata.ResourceData.HasChange("vnet_content_share_enabled") {
				siteEnvelope.VnetContentShareEnabled = pointer.To(data.VNETContentShareEnabled)
			}

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := data.VirtualNetworkSubnetId
				if subnetId == "" {
					if _, err := client.DeleteSwiftVirtualNetwork(ctx, *id); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					siteEnvelope.VirtualNetworkSubnetId = empty
				} else {
					siteEnvelope.VirtualNetworkSubnetId = pointer.To(subnetId)
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(data.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				existing.Model.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("key_vault_reference_identity_id") && data.KeyvaultReferenceIdentityId != "" {
				siteEnvelope.KeyVaultReferenceIdentity = pointer.To(data.KeyvaultReferenceIdentityId)
			}

			existing.Model.Properties = pointer.To(siteEnvelope)

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = pointer.To(data.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStrings := helpers.ExpandConnectionStrings(data.ConnectionStrings)
				if connectionStrings.Properties != nil {
					if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStrings); err != nil {
						return fmt.Errorf("setting Connection Strings for Linux %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("ftp_publish_basic_authentication_enabled") {
				policy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: data.FtpPublishBasicAuthEnabled,
					},
				}

				if _, err := client.UpdateFtpAllowed(ctx, *id, policy); err != nil {
					return fmt.Errorf("updating FTP publish basic authentication policy for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("scm_publish_basic_authentication_enabled") {
				policy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: data.SCMPublishBasicAuthEnabled,
					},
				}

				if _, err := client.UpdateScmAllowed(ctx, *id, policy); err != nil {
					return fmt.Errorf("updating SCM publish basic authentication policy for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
