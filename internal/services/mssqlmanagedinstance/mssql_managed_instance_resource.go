// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedInstanceModel struct {
	AdministratorLogin                string                              `tfschema:"administrator_login"`
	AdministratorLoginPassword        string                              `tfschema:"administrator_login_password"`
	Collation                         string                              `tfschema:"collation"`
	DnsZonePartnerId                  string                              `tfschema:"dns_zone_partner_id"`
	DnsZone                           string                              `tfschema:"dns_zone"`
	Fqdn                              string                              `tfschema:"fqdn"`
	Identity                          []identity.SystemOrUserAssignedList `tfschema:"identity"`
	LicenseType                       string                              `tfschema:"license_type"`
	Location                          string                              `tfschema:"location"`
	MaintenanceConfigurationName      string                              `tfschema:"maintenance_configuration_name"`
	MinimumTlsVersion                 string                              `tfschema:"minimum_tls_version"`
	Name                              string                              `tfschema:"name"`
	ProxyOverride                     string                              `tfschema:"proxy_override"`
	PublicDataEndpointEnabled         bool                                `tfschema:"public_data_endpoint_enabled"`
	ResourceGroupName                 string                              `tfschema:"resource_group_name"`
	SkuName                           string                              `tfschema:"sku_name"`
	StorageAccountType                string                              `tfschema:"storage_account_type"`
	StorageSizeInGb                   int64                               `tfschema:"storage_size_in_gb"`
	SubnetId                          string                              `tfschema:"subnet_id"`
	TimezoneId                        string                              `tfschema:"timezone_id"`
	VCores                            int64                               `tfschema:"vcores"`
	AzureActiveDirectoryAdministrator []AzureActiveDirectoryAdministrator `tfschema:"azure_active_directory_administrator"`
	ZoneRedundantEnabled              bool                                `tfschema:"zone_redundant_enabled"`
	Tags                              map[string]string                   `tfschema:"tags"`
}

type AzureActiveDirectoryAdministrator struct {
	LoginUserName                    string `tfschema:"login_username"`
	ObjectID                         string `tfschema:"object_id"`
	AzureADAuthenticationOnlyEnabled bool   `tfschema:"azuread_authentication_only_enabled"`
	TenantID                         string `tfschema:"tenant_id"`
}

var _ sdk.Resource = MsSqlManagedInstanceResource{}
var _ sdk.ResourceWithUpdate = MsSqlManagedInstanceResource{}
var _ sdk.ResourceWithCustomizeDiff = MsSqlManagedInstanceResource{}

type MsSqlManagedInstanceResource struct{}

func (r MsSqlManagedInstanceResource) ResourceType() string {
	return "azurerm_mssql_managed_instance"
}

func (r MsSqlManagedInstanceResource) ModelObject() interface{} {
	return &MsSqlManagedInstanceModel{}
}

func (r MsSqlManagedInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedInstanceID
}

func (r MsSqlManagedInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlManagedInstanceServerName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"BC_Gen4",
				"BC_Gen5",
				"BC_Gen8IH",
				"BC_Gen8IM",
				"GP_Gen4",
				"GP_Gen5",
				"GP_Gen8IH",
				"GP_Gen8IM",
			}, false),
		},

		"license_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"LicenseIncluded",
				"BasePrice",
			}, true),
		},

		"storage_size_in_gb": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(32, 16384),
		},

		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"vcores": {
			Type:     schema.TypeInt,
			Required: true,
			ValidateFunc: validation.IntInSlice([]int{
				4,
				6,
				8,
				10,
				12,
				16,
				20,
				24,
				32,
				40,
				48,
				56,
				64,
				80,
				96,
				128,
			}),
		},

		"administrator_login": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"administrator_login", "azure_active_directory_administrator"},
			RequiredWith: []string{"administrator_login", "administrator_login_password"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"administrator_login_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			AtLeastOneOf: []string{"administrator_login_password", "azure_active_directory_administrator"},
			RequiredWith: []string{"administrator_login", "administrator_login_password"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"azure_active_directory_administrator": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"login_username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"object_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"azuread_authentication_only_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"collation": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "SQL_Latin1_General_CP1_CI_AS",
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"dns_zone_partner_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validate.ManagedInstanceID,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"maintenance_configuration_name": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "SQL_Default",
			ValidateFunc: validation.Any(
				validation.StringInSlice([]string{
					"SQL_Default",
				}, false),
				validation.StringMatch(regexp.MustCompile(`^SQL_[A-Za-z0-9]+_MI_\d+$`), "expected a name in the format `SQL_{Location}_MI_{Number}` or `SQL_Default`"),
			),
		},

		"minimum_tls_version": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "1.2",
			ValidateFunc: validation.StringInSlice([]string{
				"1.0",
				"1.1",
				"1.2",
			}, false),
		},

		"proxy_override": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  string(sql.ManagedInstanceProxyOverrideDefault),
			ValidateFunc: validation.StringInSlice([]string{
				string(sql.ManagedInstanceProxyOverrideDefault),
				string(sql.ManagedInstanceProxyOverrideRedirect),
				string(sql.ManagedInstanceProxyOverrideProxy),
			}, false),
		},

		"public_data_endpoint_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(sql.StorageAccountTypeGRS),
			ValidateFunc: validation.StringInSlice([]string{
				string(sql.StorageAccountTypeGRS),
				string(sql.StorageAccountTypeLRS),
				string(sql.StorageAccountTypeZRS),
			}, false),
		},

		"timezone_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "UTC",
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"zone_redundant_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": tags.Schema(),
	}
}

func (r MsSqlManagedInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dns_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"fqdn": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (r MsSqlManagedInstanceResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			// dns_zone_partner_id can only be set on init
			if oldVal, newVal := rd.GetChange("dns_zone_partner_id"); oldVal.(string) == "" && newVal.(string) != "" {
				if err := rd.ForceNew("dns_zone_partner_id"); err != nil {
					return err
				}
			}

			// system-assigned identity can't be removed due to https://github.com/Azure/azure-rest-api-specs/issues/16838
			if oldVal, newVal := rd.GetChange("identity.#"); oldVal.(int) == 1 && newVal.(int) == 0 {
				if err := rd.ForceNew("identity"); err != nil {
					return err
				}
			}

			_, aadAdminOk := rd.GetOk("azure_active_directory_administrator")
			authOnlyEnabled := rd.Get("azure_active_directory_administrator.0.azuread_authentication_only_enabled").(bool)
			_, loginOk := rd.GetOk("administrator_login")
			_, pwsOk := rd.GetOk("administrator_login_password")
			if aadAdminOk && !authOnlyEnabled && (!loginOk || !pwsOk) {
				return fmt.Errorf("`administrator_login` and `administrator_login_password` are required when `azuread_authentication_only_enabled` is false")
			}

			return nil
		},
	}
}

func (r MsSqlManagedInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 24 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MsSqlManagedInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := parse.NewManagedInstanceID(subscriptionId, model.ResourceGroupName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sku, err := r.expandSkuName(model.SkuName)
			if err != nil {
				return fmt.Errorf("expanding `sku_name` for SQL Managed Instance Server %q: %v", id.ID(), err)
			}

			maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(subscriptionId, model.MaintenanceConfigurationName)

			parameters := sql.ManagedInstance{
				Sku:      sku,
				Identity: r.expandIdentity(model.Identity),
				Location: pointer.To(location.Normalize(model.Location)),
				ManagedInstanceProperties: &sql.ManagedInstanceProperties{
					AdministratorLogin:         pointer.To(model.AdministratorLogin),
					AdministratorLoginPassword: pointer.To(model.AdministratorLoginPassword),
					Collation:                  pointer.To(model.Collation),
					DNSZonePartner:             pointer.To(model.DnsZonePartnerId),
					LicenseType:                sql.ManagedInstanceLicenseType(model.LicenseType),
					MaintenanceConfigurationID: pointer.To(maintenanceConfigId.ID()),
					MinimalTLSVersion:          pointer.To(model.MinimumTlsVersion),
					ProxyOverride:              sql.ManagedInstanceProxyOverride(model.ProxyOverride),
					PublicDataEndpointEnabled:  pointer.To(model.PublicDataEndpointEnabled),
					StorageAccountType:         sql.StorageAccountType(model.StorageAccountType),
					StorageSizeInGB:            pointer.To(int32(model.StorageSizeInGb)),
					SubnetID:                   pointer.To(model.SubnetId),
					TimezoneID:                 pointer.To(model.TimezoneId),
					VCores:                     pointer.To(int32(model.VCores)),
					ZoneRedundant:              pointer.To(model.ZoneRedundantEnabled),
				},
				Tags: tags.FromTypedObject(model.Tags),
			}

			administrators, err := expandMsSqlManagedInstanceExternalAdministrators(model.AzureActiveDirectoryAdministrator)
			if err != nil {
				return fmt.Errorf("expanding `azure_active_directory_administrator` for SQL Managed Instance Server %q: %v", id.ID(), err)
			}
			parameters.ManagedInstanceProperties.Administrators = administrators

			if parameters.Identity != nil && len(parameters.Identity.UserAssignedIdentities) > 0 {
				for k := range parameters.Identity.UserAssignedIdentities {
					parameters.ManagedInstanceProperties.PrimaryUserAssignedIdentityID = pointer.To(k)
					break
				}
			}

			metadata.Logger.Infof("Creating %s", id)

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				if response.WasConflict(future.Response()) {
					return fmt.Errorf("sql managed instance names need to be globally unique and %q is already in use", id.Name)
				}

				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MsSqlManagedInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 24 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient
			adminClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAdministratorsClient
			azureADAuthenticationOnlyClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAzureADOnlyAuthenticationsClient

			id, err := parse.ManagedInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedInstanceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("Updating %s", id)

			sku, err := r.expandSkuName(state.SkuName)
			if err != nil {
				return fmt.Errorf("expanding `sku_name` for SQL Managed Instance Server %q: %v", id.ID(), err)
			}

			properties := sql.ManagedInstance{
				Sku:      sku,
				Identity: r.expandIdentity(state.Identity),
				Location: pointer.To(location.Normalize(state.Location)),
				ManagedInstanceProperties: &sql.ManagedInstanceProperties{
					DNSZonePartner:            pointer.To(state.DnsZonePartnerId),
					LicenseType:               sql.ManagedInstanceLicenseType(state.LicenseType),
					MinimalTLSVersion:         pointer.To(state.MinimumTlsVersion),
					ProxyOverride:             sql.ManagedInstanceProxyOverride(state.ProxyOverride),
					PublicDataEndpointEnabled: pointer.To(state.PublicDataEndpointEnabled),
					StorageSizeInGB:           pointer.To(int32(state.StorageSizeInGb)),
					VCores:                    pointer.To(int32(state.VCores)),
					ZoneRedundant:             pointer.To(state.ZoneRedundantEnabled),
				},
				Tags: tags.FromTypedObject(state.Tags),
			}

			if properties.Identity != nil && len(properties.Identity.UserAssignedIdentities) > 0 {
				for k := range properties.Identity.UserAssignedIdentities {
					properties.ManagedInstanceProperties.PrimaryUserAssignedIdentityID = pointer.To(k)
					break
				}
			}

			if metadata.ResourceData.HasChange("maintenance_configuration_name") {
				maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(id.SubscriptionId, state.MaintenanceConfigurationName)
				properties.MaintenanceConfigurationID = pointer.To(maintenanceConfigId.ID())
			}

			if metadata.ResourceData.HasChange("administrator_login_password") {
				properties.AdministratorLoginPassword = pointer.To(state.AdministratorLoginPassword)
			}

			if metadata.ResourceData.HasChange("azure_active_directory_administrator") {
				// Need to check if Microsoft AAD authentication only is enabled or not before calling delete, else you will get the following error:
				// InvalidManagedServerAADOnlyAuthNoAADAdminPropertyName: AAD Admin is not configured,
				// AAD Admin must be set before enabling/disabling AAD Authentication Only.
				log.Printf("[INFO] Checking if AAD Administrator exists")
				aadAdminExists := false
				resp, err := adminClient.Get(ctx, id.ResourceGroup, id.Name)
				if err != nil {
					if !utils.ResponseWasNotFound(resp.Response) {
						return fmt.Errorf("retrieving the Administrators of %s: %+v", *id, err)
					}
				} else {
					aadAdminExists = true
				}

				aadAdminRemoved := checkAADAdminRemoved(metadata.ResourceData)
				if aadAdminRemoved && aadAdminExists {
					log.Printf("[INFO] Disabling AAD Authentication Only for %s: %+v", *id, err)
					future, err := azureADAuthenticationOnlyClient.Delete(ctx, id.ResourceGroup, id.Name)
					if err != nil {
						return fmt.Errorf("disabling AAD Authentication Only for %s: %+v", *id, err)
					}

					if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
						return fmt.Errorf("waiting to disable AAD Authentication Only for %s: %+v", *id, err)
					}

					log.Printf("[INFO] Removing AAD Administrator for %s: %+v", *id, err)
					resp, err := adminClient.Delete(ctx, id.ResourceGroup, id.Name)
					if err != nil {
						return fmt.Errorf("removing the AAD Administrator for %s: %+v", *id, err)
					}
					if err = resp.WaitForCompletionRef(ctx, client.Client); err != nil {
						return fmt.Errorf("waiting for removal of AAD Administrator for %s: %+v", *id, err)
					}
				}

				aadAdminProps, err := expandMsSqlManagedInstanceAdministrators(state.AzureActiveDirectoryAdministrator)

				if err != nil {
					return err
				}
				if aadAdminProps != nil {
					aadAdminChanged := checkAADAdminChanged(metadata.ResourceData)
					if aadAdminChanged {
						// if enabling AAD Authentication only, AAD admin must be set first.
						log.Printf("[INFO] Creating/updating AAD Administrator for %s: %+v", *id, err)
						future, err := adminClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *aadAdminProps)
						if err != nil {
							return fmt.Errorf("creating AAD Administrator for %s: %+v", *id, err)
						}
						if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
							return fmt.Errorf("waiting for creation of AAD Administrator for %s: %+v", *id, err)
						}

					}

					if metadata.ResourceData.HasChange("azure_active_directory_administrator.0.azuread_authentication_only_enabled") {
						aadOnlyAuthenticationsProps := sql.ManagedInstanceAzureADOnlyAuthentication{
							ManagedInstanceAzureADOnlyAuthProperties: &sql.ManagedInstanceAzureADOnlyAuthProperties{
								AzureADOnlyAuthentication: pointer.To(expandMsSqlManagedInstanceAadAuthenticationOnly(state.AzureActiveDirectoryAdministrator)),
							},
						}
						resp, err := azureADAuthenticationOnlyClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, aadOnlyAuthenticationsProps)
						if err != nil {
							return fmt.Errorf("setting `azuread_authentication_only_enabled` for %s: %+v", *id, err)
						}
						if err = resp.WaitForCompletionRef(ctx, client.Client); err != nil {
							return fmt.Errorf("waiting to set `azuread_authentication_only_enabled` for  %s: %+v", *id, err)
						}
					}
				}
			}

			metadata.Logger.Infof("Updating %s", *id)

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, properties)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func checkAADAdminRemoved(d *schema.ResourceData) bool {
	old, new := d.GetChange("azure_active_directory_administrator")
	if len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
		return true
	}
	return false
}

func checkAADAdminChanged(d *schema.ResourceData) bool {
	oldLogin, newLogin := d.GetChange("azure_active_directory_administrator.0.login_username")
	oldObjId, newObjId := d.GetChange("azure_active_directory_administrator.0.object_id")
	oldTenantID, newTenantID := d.GetChange("azure_active_directory_administrator.0.tenant_id")
	return oldLogin.(string) != newLogin.(string) || oldObjId.(string) != newObjId.(string) || oldTenantID.(string) != newTenantID.(string)
}

func (r MsSqlManagedInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient

			id, err := parse.ManagedInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedInstanceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := MsSqlManagedInstanceModel{
				Name:              id.Name,
				Location:          location.NormalizeNilable(existing.Location),
				ResourceGroupName: id.ResourceGroup,
				Identity:          r.flattenIdentity(existing.Identity),
				Tags:              tags.ToTypedObject(existing.Tags),

				// This value is not returned, so we'll just set whatever is in the state/config
				AdministratorLoginPassword: state.AdministratorLoginPassword,
				// This value is not returned, so we'll just set whatever is in the state/config
				DnsZonePartnerId: state.DnsZonePartnerId,
			}

			if sku := existing.Sku; sku != nil && sku.Name != nil {
				model.SkuName = r.normalizeSku(*sku.Name)
			}

			if props := existing.ManagedInstanceProperties; props != nil {
				model.LicenseType = string(props.LicenseType)
				model.ProxyOverride = string(props.ProxyOverride)
				model.StorageAccountType = string(props.StorageAccountType)

				if props.AdministratorLogin != nil {
					model.AdministratorLogin = *props.AdministratorLogin
				}
				if props.Administrators != nil {
					model.AzureActiveDirectoryAdministrator = flattenMsSqlManagedInstanceAdministrators(*props.Administrators)
				}
				if props.Collation != nil {
					model.Collation = *props.Collation
				}
				if props.DNSZone != nil {
					model.DnsZone = *props.DNSZone
				}
				if props.FullyQualifiedDomainName != nil {
					model.Fqdn = *props.FullyQualifiedDomainName
				}
				if props.MaintenanceConfigurationID != nil {
					maintenanceConfigId, err := publicmaintenanceconfigurations.ParsePublicMaintenanceConfigurationIDInsensitively(*props.MaintenanceConfigurationID)
					if err != nil {
						return err
					}
					model.MaintenanceConfigurationName = maintenanceConfigId.PublicMaintenanceConfigurationName
				}
				if props.MinimalTLSVersion != nil {
					model.MinimumTlsVersion = *props.MinimalTLSVersion
				}
				if props.PublicDataEndpointEnabled != nil {
					model.PublicDataEndpointEnabled = *props.PublicDataEndpointEnabled
				}
				if props.StorageSizeInGB != nil {
					model.StorageSizeInGb = int64(*props.StorageSizeInGB)
				}
				if props.SubnetID != nil {
					model.SubnetId = *props.SubnetID
				}
				if props.TimezoneID != nil {
					model.TimezoneId = *props.TimezoneID
				}
				if props.VCores != nil {
					model.VCores = int64(*props.VCores)
				}

				if props.ZoneRedundant != nil {
					model.ZoneRedundantEnabled = *props.ZoneRedundant
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MsSqlManagedInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 24 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient

			id, err := parse.ManagedInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlManagedInstanceResource) expandIdentity(input []identity.SystemOrUserAssignedList) *sql.ResourceIdentity {
	if len(input) == 0 {
		return nil
	}

	// Workaround for issue https://github.com/Azure/azure-rest-api-specs/issues/16838
	if input[0].Type == identity.TypeNone {
		return nil
	}

	var identityIds map[string]*sql.UserIdentity
	if len(input[0].IdentityIds) != 0 {
		identityIds = map[string]*sql.UserIdentity{}
		for _, id := range input[0].IdentityIds {
			identityIds[id] = &sql.UserIdentity{}
		}
	}

	return &sql.ResourceIdentity{
		Type:                   sql.IdentityType(input[0].Type),
		UserAssignedIdentities: identityIds,
	}
}

func (r MsSqlManagedInstanceResource) flattenIdentity(input *sql.ResourceIdentity) []identity.SystemOrUserAssignedList {
	if input == nil {
		return nil
	}

	principalId := ""
	if input.PrincipalID != nil {
		principalId = input.PrincipalID.String()
	}

	tenantId := ""
	if input.TenantID != nil {
		tenantId = input.TenantID.String()
	}

	var identityIds = make([]string, 0)
	for k := range input.UserAssignedIdentities {
		parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(k)
		if err != nil {
			continue
		}
		identityIds = append(identityIds, parsedId.ID())
	}

	return []identity.SystemOrUserAssignedList{{
		Type:        identity.Type(input.Type),
		PrincipalId: principalId,
		TenantId:    tenantId,
		IdentityIds: identityIds,
	}}
}

func (r MsSqlManagedInstanceResource) expandSkuName(skuName string) (*sql.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier string
	switch parts[0] {
	case "GP":
		tier = "GeneralPurpose"
	case "BC":
		tier = "BusinessCritical"
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	return &sql.Sku{
		Name:   pointer.To(skuName),
		Tier:   pointer.To(tier),
		Family: pointer.To(parts[1]),
	}, nil
}

func (r MsSqlManagedInstanceResource) normalizeSku(sku string) string {
	switch sku {
	case "MIBC64G8IH", "BC_G8IH":
		return "BC_Gen8IH"
	case "MIBC64G8IM", "BC_G8IM":
		return "BC_Gen8IM"
	case "MIGP4G8IH", "GP_G8IH":
		return "GP_Gen8IH"
	case "MIGP4G8IM", "GP_G8IM":
		return "GP_Gen8IM"
	}

	return sku
}

func expandMsSqlManagedInstanceAadAuthenticationOnly(input []AzureActiveDirectoryAdministrator) bool {
	if len(input) == 0 {
		return false
	}

	if ok := input[0].AzureADAuthenticationOnlyEnabled; ok {
		return input[0].AzureADAuthenticationOnlyEnabled
	}

	return false
}

func expandMsSqlManagedInstanceExternalAdministrators(input []AzureActiveDirectoryAdministrator) (*sql.ManagedInstanceExternalAdministrator, error) {
	if len(input) == 0 {
		return nil, nil
	}

	admin := input[0]
	sid, err := uuid.FromString(admin.ObjectID)
	if err != nil {
		return nil, err
	}

	adminParams := sql.ManagedInstanceExternalAdministrator{
		AdministratorType: sql.AdministratorTypeActiveDirectory,
		Login:             pointer.To(admin.LoginUserName),
		Sid:               pointer.To(sid),
	}

	if admin.TenantID != "" {
		tenantId, err := uuid.FromString(admin.TenantID)
		if err != nil {
			return nil, err
		}
		adminParams.TenantID = pointer.To(tenantId)
	}

	adminParams.AzureADOnlyAuthentication = pointer.To(admin.AzureADAuthenticationOnlyEnabled)

	return &adminParams, nil
}

func expandMsSqlManagedInstanceAdministrators(input []AzureActiveDirectoryAdministrator) (*sql.ManagedInstanceAdministrator, error) {
	if len(input) == 0 {
		return nil, nil
	}

	admin := input[0]
	sid, err := uuid.FromString(admin.ObjectID)
	if err != nil {
		return nil, err
	}

	adminProps := sql.ManagedInstanceAdministrator{
		ManagedInstanceAdministratorProperties: &sql.ManagedInstanceAdministratorProperties{
			AdministratorType: pointer.To(string(sql.AdministratorTypeActiveDirectory)),
			Login:             pointer.To(admin.LoginUserName),
			Sid:               pointer.To(sid),
		},
	}

	if admin.TenantID != "" {
		tenantId, err := uuid.FromString(admin.TenantID)
		if err != nil {
			return nil, err
		}
		adminProps.ManagedInstanceAdministratorProperties.TenantID = pointer.To(tenantId)
	}

	return pointer.To(adminProps), nil
}

func flattenMsSqlManagedInstanceAdministrators(admin sql.ManagedInstanceExternalAdministrator) []AzureActiveDirectoryAdministrator {
	results := make([]AzureActiveDirectoryAdministrator, 0)
	return append(results, AzureActiveDirectoryAdministrator{
		LoginUserName:                    pointer.From(admin.Login),
		ObjectID:                         pointer.From(admin.Sid).String(),
		TenantID:                         pointer.From(admin.TenantID).String(),
		AzureADAuthenticationOnlyEnabled: pointer.From(admin.AzureADOnlyAuthentication),
	})
}
