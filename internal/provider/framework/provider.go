// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	providerfunction "github.com/hashicorp/terraform-provider-azurerm/internal/provider/function"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk/frameworkhelpers"

	pluginsdkprovider "github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

type azureRmFrameworkProvider struct {
	V2Provider interface{ Meta() interface{} }
	ProviderConfig
}

var _ provider.Provider = &azureRmFrameworkProvider{}

var _ provider.ProviderWithFunctions = &azureRmFrameworkProvider{}

var _ provider.ProviderWithEphemeralResources = &azureRmFrameworkProvider{}

func (p *azureRmFrameworkProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		providerfunction.NewNormaliseResourceIDFunction,
		providerfunction.NewParseResourceIDFunction,
	}
}

func NewFrameworkProvider(primary interface{ Meta() interface{} }) provider.Provider {
	return &azureRmFrameworkProvider{
		V2Provider: primary,
	}
}

func NewFrameworkV5Provider() provider.Provider {
	return &azureRmFrameworkProvider{}
}

func (p *azureRmFrameworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "azurerm"
}

func (p *azureRmFrameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"subscription_id": schema.StringAttribute{
				// Note: There is no equivalent of `DefaultFunc` in the provider schema package. This property is Required, but can be
				// set via env var instead of provider config, so needs to be toggled in schema based on the presence of that env var.
				Optional:    true,
				Description: "The Subscription ID which should be used.",
			},

			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Client ID which should be used.",
			},

			"client_id_file_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to a file containing the Client ID which should be used.",
			},

			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Tenant ID which should be used.",
			},

			"auxiliary_tenant_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(3),
				},
			},

			"environment": schema.StringAttribute{
				Optional:    true,
				Description: "The Cloud Environment which should be used. Possible values are public, usgovernment, and china. Defaults to public. Not used and should not be specified when `metadata_host` is specified.",
			},

			"metadata_host": schema.StringAttribute{
				Optional:    true,
				Description: "The Hostname which should be used for the Azure Metadata Service.",
			},

			// Client Certificate specific fields
			"client_certificate": schema.StringAttribute{
				Optional:    true,
				Description: "Base64 encoded PKCS#12 certificate bundle to use when authenticating as a Service Principal using a Client Certificate",
			},

			"client_certificate_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate.",
			},

			"client_certificate_password": schema.StringAttribute{
				Optional:    true,
				Description: "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate",
			},

			// Client Secret specific fields
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			"client_secret_file_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to a file containing the Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			"ado_pipeline_service_connection_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Azure DevOps Pipeline Service Connection ID.",
			},

			// OIDC specific fields
			"oidc_request_token": schema.StringAttribute{
				Optional:    true,
				Description: "The bearer token for the request to the OIDC provider. For use when authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_request_url": schema.StringAttribute{
				Optional:    true,
				Description: "The URL for the OIDC provider from which to request an ID token. For use when authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token": schema.StringAttribute{
				Optional:    true,
				Description: "The OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token_file_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to a file containing an OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},

			"use_oidc": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow OpenID Connect to be used for authentication",
			},

			// Managed Service Identity specific fields
			"use_msi": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow Managed Service Identity to be used for Authentication.",
			},
			"msi_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "The path to a custom endpoint for Managed Service Identity - in most circumstances this should be detected automatically. ",
			},

			// Azure CLI specific fields
			"use_cli": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow Azure CLI to be used for Authentication.",
			},

			// Azure AKS Workload Identity fields
			"use_aks_workload_identity": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow Azure AKS Workload Identity to be used for Authentication.",
			},

			// Managed Tracking GUID for User-agent
			"partner_id": schema.StringAttribute{
				Optional:    true,
				Description: "A GUID/UUID that is registered with Microsoft to facilitate partner resource usage attribution.",
			},

			"disable_correlation_request_id": schema.BoolAttribute{
				Optional:    true,
				Description: "This will disable the x-ms-correlation-request-id header.",
			},

			"disable_terraform_partner_id": schema.BoolAttribute{
				Optional:    true,
				Description: "This will disable the Terraform Partner ID which is used if a custom `partner_id` isn't specified.",
			},

			// Advanced feature flags
			"skip_provider_registration": schema.BoolAttribute{
				Optional:           true,
				Description:        "Should the AzureRM Provider skip registering all of the Resource Providers that it supports, if they're not already registered?",
				DeprecationMessage: "This property is deprecated and will be removed in v5.0 of the AzureRM provider. Please use the `resource_provider_registrations` property instead.",
			},

			"storage_use_azuread": schema.BoolAttribute{
				Optional:    true,
				Description: "Should the AzureRM Provider use Azure AD Authentication when accessing the Storage Data Plane APIs?",
			},

			"resource_provider_registrations": schema.StringAttribute{
				Optional:    true,
				Description: "The set of Resource Providers which should be automatically registered for the subscription.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						resourceproviders.ProviderRegistrationsNone,
						resourceproviders.ProviderRegistrationsLegacy,
						resourceproviders.ProviderRegistrationsCore,
						resourceproviders.ProviderRegistrationsExtended,
						resourceproviders.ProviderRegistrationsAll,
					),
				},
			},

			"resource_providers_to_register": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "A list of Resource Providers to explicitly register for the subscription, in addition to those specified by the `resource_provider_registrations` property.",
				Validators: []validator.List{
					frameworkhelpers.WrappedListValidator{
						Func:         resourceproviders.EnhancedValidate,
						Desc:         "EnhancedValidate returns a validation function which attempts to validate the Resource Provider against the list of Resource Provider supported by this Azure Environment.",
						MarkdownDesc: "EnhancedValidate returns a validation function which attempts to validate the Resource Provider against the list of Resource Provider supported by this Azure Environment.",
					},
				},
			},
		},

		Blocks: map[string]schema.Block{
			"features": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 1),
				},
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"api_management": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"purge_soft_delete_on_destroy": schema.BoolAttribute{
										Optional: true,
									},

									"recover_soft_deleted": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"app_configuration": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"purge_soft_delete_on_destroy": schema.BoolAttribute{
										Optional: true,
									},

									"recover_soft_deleted": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"application_insights": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"disable_generated_rule": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"cognitive_account": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"purge_soft_delete_on_destroy": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"key_vault": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"purge_soft_delete_on_destroy": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault` resources will be permanently deleted (e.g purged), when destroyed",
										Optional:    true,
									},

									"purge_soft_deleted_certificates_on_destroy": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault_certificate` resources will be permanently deleted (e.g purged), when destroyed",
										Optional:    true,
									},

									"purge_soft_deleted_keys_on_destroy": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault_key` resources will be permanently deleted (e.g purged), when destroyed",
										Optional:    true,
									},

									"purge_soft_deleted_secrets_on_destroy": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault_secret` resources will be permanently deleted (e.g purged), when destroyed",
										Optional:    true,
									},

									"purge_soft_deleted_hardware_security_modules_on_destroy": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault_managed_hardware_security_module` resources will be permanently deleted (e.g purged), when destroyed",
										Optional:    true,
									},

									"purge_soft_deleted_hardware_security_module_keys_on_destroy": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault_managed_hardware_security_module_key` resources will be permanently deleted (e.g purged), when destroyed",
										Optional:    true,
									},

									"recover_soft_deleted_certificates": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault_certificate` resources will be restored, instead of creating new ones",
										Optional:    true,
									},

									"recover_soft_deleted_key_vaults": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault` resources will be restored, instead of creating new ones",
										Optional:    true,
									},

									"recover_soft_deleted_keys": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault_key` resources will be restored, instead of creating new ones",
										Optional:    true,
									},

									"recover_soft_deleted_secrets": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault_secret` resources will be restored, instead of creating new ones",
										Optional:    true,
									},

									"recover_soft_deleted_hardware_security_module_keys": schema.BoolAttribute{
										Description: "When enabled soft-deleted `azurerm_key_vault_managed_hardware_security_module_key` resources will be restored, instead of creating new ones",
										Optional:    true,
									},
								},
							},
						},
						"log_analytics_workspace": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"permanently_delete_on_destroy": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"template_deployment": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"delete_nested_items_during_deletion": schema.BoolAttribute{
										Required: true,
									},
								},
							},
						},
						"virtual_machine": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"delete_os_disk_on_deletion": schema.BoolAttribute{
										Optional: true,
									},
									"skip_shutdown_and_force_delete": schema.BoolAttribute{
										Optional: true,
									},
									"detach_implicit_data_disk_on_deletion": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"virtual_machine_scale_set": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"force_delete": schema.BoolAttribute{
										Optional: true,
									},
									"reimage_on_manual_upgrade": schema.BoolAttribute{
										Optional: true,
									},
									"roll_instances_when_required": schema.BoolAttribute{
										Optional: true,
									},
									"scale_to_zero_before_deletion": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"resource_group": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"prevent_deletion_if_contains_resources": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"managed_disk": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"expand_without_downtime": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"storage": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"data_plane_available": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"subscription": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"prevent_cancellation_on_destroy": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"postgresql_flexible_server": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"restart_server_on_configuration_value_change": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"machine_learning": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"purge_soft_deleted_workspace_on_destroy": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"recovery_service": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"vm_backup_stop_protection_and_retain_data_on_destroy": schema.BoolAttribute{
										Optional: true,
									},
									"vm_backup_suspend_protection_and_retain_data_on_destroy": schema.BoolAttribute{
										Optional: true,
									},
									"purge_protected_items_from_vault_on_destroy": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"recovery_services_vaults": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"recover_soft_deleted_backup_protected_vm": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"netapp": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"delete_backups_on_backup_vault_destroy": schema.BoolAttribute{
										Optional:    true,
										Description: "When enabled, backups will be deleted when the `azurerm_netapp_backup_vault` resource is destroyed",
									},
									"prevent_volume_destruction": schema.BoolAttribute{
										Description: "When enabled, the volume will not be destroyed, safeguarding from severe data loss",
										Optional:    true,
									},
								},
							},
						},
						"databricks_workspace": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"force_delete": schema.BoolAttribute{
										Optional:    true,
										Description: "When enabled, the managed resource group that contains the Unity Catalog data will be forcibly deleted when the workspace is destroyed, regardless of contents.",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if !features.FivePointOh() {
		response.Schema.Blocks["features"].(schema.ListNestedBlock).NestedObject.Blocks["virtual_machine"].(schema.ListNestedBlock).NestedObject.Attributes["graceful_shutdown"] = schema.BoolAttribute{
			Optional:           true,
			DeprecationMessage: "'graceful_shutdown' has been deprecated and will be removed from v5.0 of the AzureRM provider.",
		}
	}
}

func (p *azureRmFrameworkProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var data ProviderModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	if p.V2Provider != nil {
		v := p.V2Provider.Meta()

		response.ResourceData = v
		response.DataSourceData = v
		response.EphemeralResourceData = v
	} else {
		p.Load(ctx, &data, request.TerraformVersion, &response.Diagnostics)

		response.DataSourceData = &p.ProviderConfig
		response.ResourceData = &p.ProviderConfig
	}
}

func (p *azureRmFrameworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	var output []func() datasource.DataSource

	for _, service := range pluginsdkprovider.SupportedFrameworkServices() {
		output = append(output, service.FrameworkDataSources()...)
	}

	return output
}

func (p *azureRmFrameworkProvider) Resources(_ context.Context) []func() resource.Resource {
	var output []func() resource.Resource

	for _, service := range pluginsdkprovider.SupportedFrameworkServices() {
		output = append(output, service.FrameworkResources()...)
	}

	return output
}

func (p *azureRmFrameworkProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	var output []func() ephemeral.EphemeralResource

	for _, service := range pluginsdkprovider.SupportedFrameworkServices() {
		output = append(output, service.EphemeralResources()...)
	}

	return output
}
