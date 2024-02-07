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

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	StorageAccountTypeGRS  = "GRS"
	StorageAccountTypeLRS  = "LRS"
	StorageAccountTypeZRS  = "ZRS"
	StorageAccountTypeGZRS = "GZRS"
)

type MsSqlManagedInstanceModel struct {
	AdministratorLogin           string                              `tfschema:"administrator_login"`
	AdministratorLoginPassword   string                              `tfschema:"administrator_login_password"`
	Collation                    string                              `tfschema:"collation"`
	DnsZonePartnerId             string                              `tfschema:"dns_zone_partner_id"`
	DnsZone                      string                              `tfschema:"dns_zone"`
	Fqdn                         string                              `tfschema:"fqdn"`
	Identity                     []identity.SystemOrUserAssignedList `tfschema:"identity"`
	LicenseType                  string                              `tfschema:"license_type"`
	Location                     string                              `tfschema:"location"`
	MaintenanceConfigurationName string                              `tfschema:"maintenance_configuration_name"`
	MinimumTlsVersion            string                              `tfschema:"minimum_tls_version"`
	Name                         string                              `tfschema:"name"`
	ProxyOverride                string                              `tfschema:"proxy_override"`
	PublicDataEndpointEnabled    bool                                `tfschema:"public_data_endpoint_enabled"`
	ResourceGroupName            string                              `tfschema:"resource_group_name"`
	ServicePrincipalType         string                              `tfschema:"service_principal_type"`
	SkuName                      string                              `tfschema:"sku_name"`
	StorageAccountType           string                              `tfschema:"storage_account_type"`
	StorageSizeInGb              int64                               `tfschema:"storage_size_in_gb"`
	SubnetId                     string                              `tfschema:"subnet_id"`
	TimezoneId                   string                              `tfschema:"timezone_id"`
	VCores                       int64                               `tfschema:"vcores"`
	MicrosoftEntraAdministrator  []MicrosoftEntraAdministrator       `tfschema:"microsoft_entra_administrator"`
	ZoneRedundantEnabled         bool                                `tfschema:"zone_redundant_enabled"`
	Tags                         map[string]string                   `tfschema:"tags"`
}

type MicrosoftEntraAdministrator struct {
	LoginUserName                           string `tfschema:"login_username"`
	ObjectID                                string `tfschema:"object_id"`
	MicrosoftEntraAuthenticationOnlyEnabled bool   `tfschema:"microsoft_entra_authentication_only_enabled"`
	TenantID                                string `tfschema:"tenant_id"`
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
			AtLeastOneOf: []string{"administrator_login", "microsoft_entra_administrator"},
			RequiredWith: []string{"administrator_login", "administrator_login_password"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"administrator_login_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			AtLeastOneOf: []string{"administrator_login_password", "microsoft_entra_administrator"},
			RequiredWith: []string{"administrator_login", "administrator_login_password"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"microsoft_entra_administrator": {
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

					"microsoft_entra_authentication_only_enabled": {
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
			Default:  string(managedinstances.ManagedInstanceProxyOverrideDefault),
			ValidateFunc: validation.StringInSlice([]string{
				string(managedinstances.ManagedInstanceProxyOverrideDefault),
				string(managedinstances.ManagedInstanceProxyOverrideRedirect),
				string(managedinstances.ManagedInstanceProxyOverrideProxy),
			}, false),
		},

		"public_data_endpoint_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"service_principal_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(managedinstances.ServicePrincipalTypeSystemAssigned),
			}, false),
		},

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  StorageAccountTypeGRS,
			ValidateFunc: validation.StringInSlice([]string{
				StorageAccountTypeGRS,
				StorageAccountTypeLRS,
				StorageAccountTypeZRS,
				StorageAccountTypeGZRS,
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

			id := commonids.NewSqlManagedInstanceID(subscriptionId, model.ResourceGroupName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id, managedinstances.GetOperationOptions{})
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sku, err := r.expandSkuName(model.SkuName)
			if err != nil {
				return fmt.Errorf("expanding `sku_name` for SQL Managed Instance Server %q: %v", id.ID(), err)
			}

			maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(subscriptionId, model.MaintenanceConfigurationName)

			parameters := managedinstances.ManagedInstance{
				Sku:      sku,
				Identity: r.expandIdentity(model.Identity),
				Location: location.Normalize(model.Location),
				Properties: &managedinstances.ManagedInstanceProperties{
					AdministratorLogin:               pointer.To(model.AdministratorLogin),
					AdministratorLoginPassword:       pointer.To(model.AdministratorLoginPassword),
					Collation:                        pointer.To(model.Collation),
					DnsZonePartner:                   pointer.To(model.DnsZonePartnerId),
					LicenseType:                      pointer.To(managedinstances.ManagedInstanceLicenseType(model.LicenseType)),
					MaintenanceConfigurationId:       pointer.To(maintenanceConfigId.ID()),
					MinimalTlsVersion:                pointer.To(model.MinimumTlsVersion),
					ProxyOverride:                    pointer.To(managedinstances.ManagedInstanceProxyOverride(model.ProxyOverride)),
					PublicDataEndpointEnabled:        pointer.To(model.PublicDataEndpointEnabled),
					RequestedBackupStorageRedundancy: pointer.To(storageAccTypeToBackupStorageRedundancy(model.StorageAccountType)),
					StorageSizeInGB:                  pointer.To(model.StorageSizeInGb),
					SubnetId:                         pointer.To(model.SubnetId),
					TimezoneId:                       pointer.To(model.TimezoneId),
					VCores:                           pointer.To(model.VCores),
					ZoneRedundant:                    pointer.To(model.ZoneRedundantEnabled),
					Administrators:                   expandMsSqlManagedInstanceExternalAdministrators(model.MicrosoftEntraAdministrator),
				},
				Tags: pointer.To(model.Tags),
			}

			if parameters.Identity != nil && len(parameters.Identity.IdentityIds) > 0 {
				for k := range parameters.Identity.IdentityIds {
					parameters.Properties.PrimaryUserAssignedIdentityId = pointer.To(k)
					break
				}
			}

			if model.ServicePrincipalType != "" {
				parameters.Properties.ServicePrincipal = &managedinstances.ServicePrincipal{
					Type: pointer.To(managedinstances.ServicePrincipalType(model.ServicePrincipalType)),
				}
			}

			metadata.Logger.Infof("Creating %s", id)

			err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
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
			microsoftEntraAuthenticationOnlyClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAzureADOnlyAuthenticationsClient

			id, err := commonids.ParseSqlManagedInstanceID(metadata.ResourceData.Id())
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

			properties := managedinstances.ManagedInstance{
				Sku:      sku,
				Identity: r.expandIdentity(state.Identity),
				Location: location.Normalize(state.Location),
				Properties: &managedinstances.ManagedInstanceProperties{
					DnsZonePartner:                   pointer.To(state.DnsZonePartnerId),
					LicenseType:                      pointer.To(managedinstances.ManagedInstanceLicenseType(state.LicenseType)),
					MinimalTlsVersion:                pointer.To(state.MinimumTlsVersion),
					ProxyOverride:                    pointer.To(managedinstances.ManagedInstanceProxyOverride(state.ProxyOverride)),
					PublicDataEndpointEnabled:        pointer.To(state.PublicDataEndpointEnabled),
					StorageSizeInGB:                  pointer.To(state.StorageSizeInGb),
					RequestedBackupStorageRedundancy: pointer.To(storageAccTypeToBackupStorageRedundancy(state.StorageAccountType)),
					VCores:                           pointer.To(state.VCores),
					ZoneRedundant:                    pointer.To(state.ZoneRedundantEnabled),
				},
				Tags: pointer.To(state.Tags),
			}

			if properties.Identity != nil && len(properties.Identity.IdentityIds) > 0 {
				for k := range properties.Identity.IdentityIds {
					properties.Properties.PrimaryUserAssignedIdentityId = pointer.To(k)
					break
				}
			}

			if metadata.ResourceData.HasChange("maintenance_configuration_name") {
				maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(id.SubscriptionId, state.MaintenanceConfigurationName)
				properties.Properties.MaintenanceConfigurationId = pointer.To(maintenanceConfigId.ID())
			}

			if metadata.ResourceData.HasChange("administrator_login_password") {
				properties.Properties.AdministratorLoginPassword = pointer.To(state.AdministratorLoginPassword)
			}

			if metadata.ResourceData.HasChange("service_principal_type") {
				properties.Properties.ServicePrincipal = &managedinstances.ServicePrincipal{}
				if state.ServicePrincipalType == "" {
					properties.Properties.ServicePrincipal.Type = pointer.To(managedinstances.ServicePrincipalTypeNone)
				} else {
					properties.Properties.ServicePrincipal.Type = pointer.To(managedinstances.ServicePrincipalType(state.ServicePrincipalType))
				}
			}

			if metadata.ResourceData.HasChange("microsoft_entra_administrator") {
				// Need to check if Microsoft Entra authentication only is enabled or not before calling delete, else you will get the following error:
				// InvalidManagedServerAADOnlyAuthNoAADAdminPropertyName: AAD Admin is not configured,
				// Microsfot Entra Admin must be set before enabling/disabling Microsfot Entra Authentication Only.
				log.Printf("[INFO] Checking if Microsoft Entra Administrator exist")
				meAdminExists := false
				resp, err := adminClient.Get(ctx, *id)
				if err != nil {
					if !utils.ResponseWasNotFound(resp.Response) {
						return fmt.Errorf("retrieving the Administrators of %s: %+v", *id, err)
					}
				} else {
					meAdminExists = true
				}

				if meAdminExists {
					future, err := microsoftEntraAuthenticationOnlyClient.Delete(ctx, id.ResourceGroup, id.Name)
					if err != nil {
						log.Printf("[INFO] Deletion of Microsoft Entra Authentication Only failed for %s: %+v", *id, err)
						return fmt.Errorf("deleting Microsoft Entra Authentication Only for %s: %+v", *id, err)
					}

					if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
						return fmt.Errorf("waiting for the Microsoft Entra Authentication Only deletion of %s: %+v", *id, err)
					}

					resp, err := adminClient.Delete(ctx, id.ResourceGroup, id.Name)
					if err != nil {
						return fmt.Errorf("deleting the Microsoft Entra Administrator of %s: %+v", *id, err)
					}
					if err = resp.WaitForCompletionRef(ctx, client.Client); err != nil {
						return fmt.Errorf("waiting for the Microsoft Entra Administrators deletion of %s: %+v", *id, err)
					}
				}

				meAdminProps, err := expandMsSqlManagedInstanceAdministrators(state.MicrosoftEntraAdministrator)
				if err != nil {
					return err
				}
				if meAdminProps != nil {
					future, err := adminClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *meAdminProps)
					if err != nil {
						return fmt.Errorf("creating Microsoft Entra Administrator of %s: %+v", *id, err)
					}
					if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
						return fmt.Errorf("waiting for creating Microsoft Entra Administrator of %s: %+v", *id, err)
					}
				}

				if meOnlyAuthentictionsEnabled := expandMsSqlManagedInstanceMeAuthentictionOnly(state.MicrosoftEntraAdministrator); meOnlyAuthentictionsEnabled {
					meOnlyAuthentictionsProps := sql.ManagedInstanceAzureADOnlyAuthentication{
						ManagedInstanceAzureADOnlyAuthProperties: &sql.ManagedInstanceAzureADOnlyAuthProperties{
							AzureADOnlyAuthentication: pointer.To(true),
						},
					}

					future, err := microsoftEntraAuthenticationOnlyClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, meOnlyAuthentictionsProps)
					if err != nil {
						return fmt.Errorf("setting `microsoft_entra_authentication_only_enabled` for %s: %+v", *id, err)
					}
					if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
						return fmt.Errorf("waiting to set `microsoft_entra_authentication_only_enabled` for  %s: %+v", *id, err)
					}
				}

				properties.Administrators = expandMsSqlManagedInstanceExternalAdministrators(state.MicrosoftEntraAdministrator)
			}

			metadata.Logger.Infof("Updating %s", *id)

			err = client.CreateOrUpdateThenPoll(ctx, *id, properties)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MsSqlManagedInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient

			id, err := commonids.ParseSqlManagedInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedInstanceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id, managedinstances.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}
			model := MsSqlManagedInstanceModel{}

			if existing.Model != nil {
				model = MsSqlManagedInstanceModel{
					Name:              id.ManagedInstanceName,
					Location:          location.NormalizeNilable(&existing.Model.Location),
					ResourceGroupName: id.ResourceGroupName,
					Identity:          r.flattenIdentity(existing.Model.Identity),
					Tags:              pointer.From(existing.Model.Tags),

					// This value is not returned, so we'll just set whatever is in the state/config
					AdministratorLoginPassword: state.AdministratorLoginPassword,
					// This value is not returned, so we'll just set whatever is in the state/config
					DnsZonePartnerId: state.DnsZonePartnerId,
				}

				if sku := existing.Model.Sku; sku != nil {
					model.SkuName = r.normalizeSku(sku.Name)
				}

				if props := existing.Model.Properties; props != nil {
					model.LicenseType = string(pointer.From(props.LicenseType))
					model.ProxyOverride = string(pointer.From(props.ProxyOverride))
					model.StorageAccountType = backupStorageRedundancyToStorageAccType(pointer.From(props.RequestedBackupStorageRedundancy))

					model.AdministratorLogin = pointer.From(props.AdministratorLogin)

					if props.Administrators != nil {
						model.MicrosoftEntraAdministrator = flattenMsSqlManagedInstanceAdministrators(*props.Administrators)
					}
					model.Collation = pointer.From(props.Collation)
					model.DnsZone = pointer.From(props.DnsZone)
					model.Fqdn = pointer.From(props.FullyQualifiedDomainName)

					if props.MaintenanceConfigurationId != nil {
						maintenanceConfigId, err := publicmaintenanceconfigurations.ParsePublicMaintenanceConfigurationIDInsensitively(*props.MaintenanceConfigurationId)
						if err != nil {
							return err
						}
						model.MaintenanceConfigurationName = maintenanceConfigId.PublicMaintenanceConfigurationName
					}

					model.MinimumTlsVersion = pointer.From(props.MinimalTlsVersion)
					model.PublicDataEndpointEnabled = pointer.From(props.PublicDataEndpointEnabled)
					model.StorageSizeInGb = pointer.From(props.StorageSizeInGB)
					model.SubnetId = pointer.From(props.SubnetId)
					model.TimezoneId = pointer.From(props.TimezoneId)
					model.VCores = pointer.From(props.VCores)
					model.ZoneRedundantEnabled = pointer.From(props.ZoneRedundant)

					model.ServicePrincipalType = ""
					if props.ServicePrincipal != nil {
						model.ServicePrincipalType = string(pointer.From(props.ServicePrincipal.Type))
					}
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

			id, err := commonids.ParseSqlManagedInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MsSqlManagedInstanceResource) expandIdentity(input []identity.SystemOrUserAssignedList) *identity.LegacySystemAndUserAssignedMap {
	if len(input) == 0 {
		return nil
	}

	// Workaround for issue https://github.com/Azure/azure-rest-api-specs/issues/16838
	if input[0].Type == identity.TypeNone {
		return nil
	}

	var identityIds map[string]identity.UserAssignedIdentityDetails
	if len(input[0].IdentityIds) != 0 {
		identityIds = map[string]identity.UserAssignedIdentityDetails{}
		for _, id := range input[0].IdentityIds {
			identityIds[id] = identity.UserAssignedIdentityDetails{}
		}
	}

	return &identity.LegacySystemAndUserAssignedMap{
		Type:        input[0].Type,
		IdentityIds: identityIds,
	}
}

func (r MsSqlManagedInstanceResource) flattenIdentity(input *identity.LegacySystemAndUserAssignedMap) []identity.SystemOrUserAssignedList {
	if input == nil {
		return nil
	}

	identityIds := make([]string, 0)
	for k := range input.IdentityIds {
		parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(k)
		if err != nil {
			continue
		}
		identityIds = append(identityIds, parsedId.ID())
	}

	return []identity.SystemOrUserAssignedList{{
		Type:        input.Type,
		PrincipalId: input.PrincipalId,
		TenantId:    input.TenantId,
		IdentityIds: identityIds,
	}}
}

func (r MsSqlManagedInstanceResource) expandSkuName(skuName string) (*managedinstances.Sku, error) {
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

	return &managedinstances.Sku{
		Name:   skuName,
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

// the StorageAccountType property has changed to RequestedBackupStorageRedundancy with 1-1 mapping of the following values:
// GRS -> Geo
// ZRS -> Zone
// LRS -> Local
// GZRS -> GeoZone
func storageAccTypeToBackupStorageRedundancy(storageAccountType string) managedinstances.BackupStorageRedundancy {
	switch storageAccountType {
	case StorageAccountTypeZRS:
		return managedinstances.BackupStorageRedundancyZone
	case StorageAccountTypeLRS:
		return managedinstances.BackupStorageRedundancyLocal
	case StorageAccountTypeGZRS:
		return managedinstances.BackupStorageRedundancyGeoZone
	}
	return managedinstances.BackupStorageRedundancyGeo
}

func backupStorageRedundancyToStorageAccType(backupStorageRedundancy managedinstances.BackupStorageRedundancy) string {
	switch backupStorageRedundancy {
	case managedinstances.BackupStorageRedundancyZone:
		return StorageAccountTypeZRS
	case managedinstances.BackupStorageRedundancyLocal:
		return StorageAccountTypeLRS
	case managedinstances.BackupStorageRedundancyGeoZone:
		return StorageAccountTypeGZRS
	}
	return StorageAccountTypeGRS
}

func expandMsSqlManagedInstanceMeAuthentictionOnly(input []MicrosoftEntraAdministrator) bool {
	if len(input) == 0 {
		return false
	}

	if ok := input[0].MicrosoftEntraAuthenticationOnlyEnabled; ok {
		return input[0].MicrosoftEntraAuthenticationOnlyEnabled
	}

	return false
}

func expandMsSqlManagedInstanceExternalAdministrators(input []MicrosoftEntraAdministrator) *managedinstances.ManagedInstanceExternalAdministrator {
	if len(input) == 0 {
		return nil
	}

	admin := input[0]
	adminParams := managedinstances.ManagedInstanceExternalAdministrator{
		AdministratorType: pointer.To(managedinstances.AdministratorTypeActiveDirectory),
		Login:             pointer.To(admin.LoginUserName),
		Sid:               pointer.To(admin.ObjectID),
	}

	if admin.TenantID != "" {
		adminParams.TenantId = pointer.To(admin.TenantID)
	}

	adminParams.AzureADOnlyAuthentication = pointer.To(admin.MicrosoftEntraAuthenticationOnlyEnabled)

	return &adminParams
}

func expandMsSqlManagedInstanceAdministrators(input []MicrosoftEntraAdministrator) (*sql.ManagedInstanceAdministrator, error) {
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

func flattenMsSqlManagedInstanceAdministrators(admin sql.ManagedInstanceExternalAdministrator) []MicrosoftEntraAdministrator {
	results := make([]MicrosoftEntraAdministrator, 0)
	return append(results, MicrosoftEntraAdministrator{
		LoginUserName:                           pointer.From(admin.Login),
		ObjectID:                                pointer.From(admin.Sid).String(),
		TenantID:                                pointer.From(admin.TenantID).String(),
		MicrosoftEntraAuthenticationOnlyEnabled: pointer.From(admin.AzureADOnlyAuthentication),
	})
}
