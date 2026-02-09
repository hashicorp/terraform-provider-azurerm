// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	keyvault "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DatabricksServerlessWorkspaceResource struct{}

var _ sdk.ResourceWithCustomizeDiff = DatabricksServerlessWorkspaceResource{}

var _ sdk.ResourceWithIdentity = DatabricksServerlessWorkspaceResource{}

type DatabricksServerlessWorkspaceModel struct {
	Name                            string                            `tfschema:"name"`
	ResourceGroupName               string                            `tfschema:"resource_group_name"`
	Location                        string                            `tfschema:"location"`
	EnhancedSecurityCompliance      []EnhancedSecurityComplianceModel `tfschema:"enhanced_security_compliance"`
	ManagedServicesCmkKeyVaultId    string                            `tfschema:"managed_services_cmk_key_vault_id"`
	ManagedServicesCmkKeyVaultKeyId string                            `tfschema:"managed_services_cmk_key_vault_key_id"`
	PublicNetworkAccessEnabled      bool                              `tfschema:"public_network_access_enabled"`
	WorkspaceId                     string                            `tfschema:"workspace_id"`
	WorkspaceUrl                    string                            `tfschema:"workspace_url"`
	Tags                            map[string]string                 `tfschema:"tags"`
}

type EnhancedSecurityComplianceModel struct {
	AutomaticClusterUpdateEnabled      bool     `tfschema:"automatic_cluster_update_enabled"`
	ComplianceSecurityProfileEnabled   bool     `tfschema:"compliance_security_profile_enabled"`
	ComplianceSecurityProfileStandards []string `tfschema:"compliance_security_profile_standards"`
	EnhancedSecurityMonitoringEnabled  bool     `tfschema:"enhanced_security_monitoring_enabled"`
}

func (r DatabricksServerlessWorkspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"enhanced_security_compliance": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"automatic_cluster_update_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      false,
						AtLeastOneOf: r.databricksServerlessWorkspaceEnhancedSecurityComplianceConstraint(),
					},
					"compliance_security_profile_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      false,
						AtLeastOneOf: r.databricksServerlessWorkspaceEnhancedSecurityComplianceConstraint(),
					},
					"compliance_security_profile_standards": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"HIPAA",
							}, false),
						},
					},
					"enhanced_security_monitoring_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      false,
						AtLeastOneOf: r.databricksServerlessWorkspaceEnhancedSecurityComplianceConstraint(),
					},
				},
			},
		},

		"managed_services_cmk_key_vault_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
		},

		"managed_services_cmk_key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyVaultValidate.KeyVaultChildID,
			RequiredWith: []string{
				"managed_service_cmk_key_vault_id",
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r DatabricksServerlessWorkspaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"workspace_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (DatabricksServerlessWorkspaceResource) ModelObject() interface{} {
	return &DatabricksServerlessWorkspaceResource{}
}

func (DatabricksServerlessWorkspaceResource) ResourceType() string {
	return "azurerm_databricks_serverless_workspace"
}

func (r DatabricksServerlessWorkspaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			ctx, cancel := context.WithTimeout(ctx, metadata.ResourceData.Timeout(schema.TimeoutCreate))
			defer cancel()

			client := metadata.Client.DataBricks.WorkspacesClient
			keyVaultClient := metadata.Client.KeyVault
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config DatabricksServerlessWorkspaceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := workspaces.NewWorkspaceID(subscriptionId, config.ResourceGroupName, config.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			location := location.Normalize(config.Location)
			encryption, err := r.expandDatabricksServerlessWorkspaceEncryption(config, id, keyVaultClient, ctx)
			if err != nil {
				return err
			}

			publicNetworkAccess := workspaces.PublicNetworkAccessEnabled
			if !config.PublicNetworkAccessEnabled {
				publicNetworkAccess = workspaces.PublicNetworkAccessDisabled
			}

			workspace := workspaces.Workspace{
				Location: location,
				Sku: &workspaces.Sku{
					Name: "premium",
				},
				Properties: workspaces.WorkspaceProperties{
					ComputeMode:                workspaces.ComputeModeServerless,
					Encryption:                 encryption,
					EnhancedSecurityCompliance: r.expandDatabricksServerlessWorkspaceEnhancedSecurityComplianceDefinition(config.EnhancedSecurityCompliance),
					PublicNetworkAccess:        &publicNetworkAccess,
				},
				Tags: pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, workspace); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return r.Read().Func(ctx, metadata)
		},
	}
}

func (r DatabricksServerlessWorkspaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			ctx, cancel := context.WithTimeout(ctx, metadata.ResourceData.Timeout(schema.TimeoutUpdate))
			defer cancel()

			client := metadata.Client.DataBricks.WorkspacesClient
			keyVaultClient := metadata.Client.KeyVault

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model is nil", id)
			}

			model := *existing.Model
			props := model.Properties
			var config DatabricksServerlessWorkspaceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("sku") {
				if model.Sku == nil {
					model.Sku = &workspaces.Sku{}
				}
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				publicNetworkAccess := workspaces.PublicNetworkAccessEnabled
				if !config.PublicNetworkAccessEnabled {
					publicNetworkAccess = workspaces.PublicNetworkAccessDisabled
				}

				props.PublicNetworkAccess = &publicNetworkAccess
			}

			if metadata.ResourceData.HasChanges("managed_services_cmk_key_vault_id", "managed_services_cmk_key_vault_key_id") {
				encryption, err := r.expandDatabricksServerlessWorkspaceEncryption(config, *id, keyVaultClient, ctx)
				if err != nil {
					return err
				}

				props.Encryption = encryption
			}

			if metadata.ResourceData.HasChange("enhanced_security_compliance") {
				props.EnhancedSecurityCompliance = r.expandDatabricksServerlessWorkspaceEnhancedSecurityComplianceDefinition(config.EnhancedSecurityCompliance)
			}

			model.Properties = props

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return r.Read().Func(ctx, metadata)
		},
	}
}

func (r DatabricksServerlessWorkspaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			ctx, cancel := context.WithTimeout(ctx, metadata.ResourceData.Timeout(schema.TimeoutRead))
			defer cancel()

			client := metadata.Client.DataBricks.WorkspacesClient
			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var config DatabricksServerlessWorkspaceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := DatabricksServerlessWorkspaceModel{
				Name:              id.WorkspaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.EnhancedSecurityCompliance = r.flattenDatabricksServerlessWorkspaceEnhancedSecurityComplianceDefinition(model.Properties.EnhancedSecurityCompliance)

				if encryption := model.Properties.Encryption; encryption != nil {
					if managedServices := encryption.Entities.ManagedServices; managedServices != nil {
						if keyVaultProperties := managedServices.KeyVaultProperties; keyVaultProperties != nil {
							if keyVaultProperties.KeyVaultUri != "" {
								key, err := keyVaultParse.NewNestedItemID(keyVaultProperties.KeyVaultUri, keyVaultParse.NestedItemTypeKey, keyVaultProperties.KeyName, keyVaultProperties.KeyVersion)
								if err == nil {
									state.ManagedServicesCmkKeyVaultKeyId = key.ID()
								}
							}
						}
					}
				}

				if model.Properties.PublicNetworkAccess != nil {
					state.PublicNetworkAccessEnabled = *model.Properties.PublicNetworkAccess == workspaces.PublicNetworkAccessEnabled
				}

				if model.Properties.WorkspaceId != nil {
					state.WorkspaceId = *model.Properties.WorkspaceId
				}

				if model.Properties.WorkspaceURL != nil {
					state.WorkspaceUrl = *model.Properties.WorkspaceURL
				}

				state.Tags = pointer.From(model.Tags)
			}

			// Always set `ManagedServicesCmkKeyVaultId` to keep the state file consistent with the configuration file
			state.ManagedServicesCmkKeyVaultId = config.ManagedServicesCmkKeyVaultId

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func (DatabricksServerlessWorkspaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			ctx, cancel := context.WithTimeout(ctx, metadata.ResourceData.Timeout(schema.TimeoutDelete))
			defer cancel()

			client := metadata.Client.DataBricks.WorkspacesClient
			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			deleteOperationOptions := workspaces.DefaultDeleteOperationOptions()
			if metadata.Client.Features.DatabricksWorkspace.ForceDelete {
				deleteOperationOptions.ForceDeletion = pointer.To(true)
			}

			if err = client.DeleteThenPoll(ctx, *id, deleteOperationOptions); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (DatabricksServerlessWorkspaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspaces.ValidateWorkspaceID
}

func (DatabricksServerlessWorkspaceResource) Identity() resourceids.ResourceId {
	return &workspaces.WorkspaceId{}
}

func (DatabricksServerlessWorkspaceResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			oldComplianceSecurityProfileEnabled, newComplianceSecurityProfileEnabled := metadata.ResourceDiff.GetChange("enhanced_security_compliance.0.compliance_security_profile_enabled")
			if metadata.ResourceDiff.HasChange("enhanced_security_compliance.0.compliance_security_profile_enabled") {
				if oldComplianceSecurityProfileEnabled.(bool) && !newComplianceSecurityProfileEnabled.(bool) {
					metadata.ResourceDiff.ForceNew("enhanced_security_compliance.0.compliance_security_profile_enabled")
				}
			}

			oldComplianceSecurityProfileStandards, newComplianceSecurityProfileStandards := metadata.ResourceDiff.GetChange("enhanced_security_compliance.0.compliance_security_profile_standards")
			if metadata.ResourceDiff.HasChange("enhanced_security_compliance.0.compliance_security_profile_standards") {
				removedStandards := oldComplianceSecurityProfileStandards.(*pluginsdk.Set).Difference(newComplianceSecurityProfileStandards.(*pluginsdk.Set))
				if removedStandards.Len() > 0 {
					metadata.ResourceDiff.ForceNew("enhanced_security_compliance.0.compliance_security_profile_standards")
				}
			}

			automaticClusterUpdateEnabled := metadata.ResourceDiff.Get("enhanced_security_compliance.0.automatic_cluster_update_enabled").(bool)
			enhancedSecurityMonitoringEnabled := metadata.ResourceDiff.Get("enhanced_security_compliance.0.enhanced_security_monitoring_enabled").(bool)
			if newComplianceSecurityProfileEnabled.(bool) && (!automaticClusterUpdateEnabled || !enhancedSecurityMonitoringEnabled) {
				return errors.New("`automatic_cluster_update_enabled` and `enhanced_security_monitoring_enabled` must be set to `true` when `compliance_security_profile_enabled` is set to `true`")
			}

			if !newComplianceSecurityProfileEnabled.(bool) && newComplianceSecurityProfileStandards.(*pluginsdk.Set).Len() > 0 {
				return errors.New("`compliance_security_profile_standards` cannot be set when `compliance_security_profile_enabled` is `false`")
			}

			return nil
		},
	}
}

func (DatabricksServerlessWorkspaceResource) expandDatabricksServerlessWorkspaceEncryption(input DatabricksServerlessWorkspaceModel, id workspaces.WorkspaceId, keyVaultClient *keyvault.Client, ctx context.Context) (*workspaces.WorkspacePropertiesEncryption, error) {
	if input.ManagedServicesCmkKeyVaultKeyId == "" {
		return nil, nil
	}

	// set default subscription as current subscription for key vault look-up...
	servicesResourceSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	if input.ManagedServicesCmkKeyVaultId != "" {
		// If they passed the 'managed_services_cmk_key_vault_id' parse the Key Vault ID
		// to extract the correct key vault subscription for the exists call...
		v, err := commonids.ParseKeyVaultID(input.ManagedServicesCmkKeyVaultId)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a Key Vault ID: %+v", input.ManagedServicesCmkKeyVaultId, err)
		}

		servicesResourceSubscriptionId = commonids.NewSubscriptionID(v.SubscriptionId)
	}

	key, err := keyVaultParse.ParseNestedItemID(input.ManagedServicesCmkKeyVaultKeyId)
	if err != nil {
		return nil, err
	}

	// make sure the key vault exists
	_, err = keyVaultClient.KeyVaultIDFromBaseUrl(ctx, servicesResourceSubscriptionId, key.KeyVaultBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed services Key Vault in subscription %q at URL %q: %+v", servicesResourceSubscriptionId, key.KeyVaultBaseUrl, err)
	}

	encryption := &workspaces.WorkspacePropertiesEncryption{
		Entities: workspaces.EncryptionEntitiesDefinition{
			ManagedServices: &workspaces.EncryptionV2{
				KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
				KeyVaultProperties: &workspaces.EncryptionV2KeyVaultProperties{
					KeyName:     key.Name,
					KeyVersion:  key.Version,
					KeyVaultUri: key.KeyVaultBaseUrl,
				},
			},
		},
	}

	return encryption, nil
}

func (DatabricksServerlessWorkspaceResource) expandDatabricksServerlessWorkspaceEnhancedSecurityComplianceDefinition(input []EnhancedSecurityComplianceModel) *workspaces.EnhancedSecurityComplianceDefinition {
	if len(input) == 0 {
		return nil
	}

	automaticClusterUpdateEnabled := workspaces.AutomaticClusterUpdateValueDisabled
	if input[0].AutomaticClusterUpdateEnabled {
		automaticClusterUpdateEnabled = workspaces.AutomaticClusterUpdateValueEnabled
	}

	complianceSecurityProfileEnabled := workspaces.ComplianceSecurityProfileValueDisabled
	if input[0].ComplianceSecurityProfileEnabled {
		complianceSecurityProfileEnabled = workspaces.ComplianceSecurityProfileValueEnabled
	}

	complianceSecurityProfileStandards := input[0].ComplianceSecurityProfileStandards
	if complianceSecurityProfileEnabled == workspaces.ComplianceSecurityProfileValueEnabled && len(complianceSecurityProfileStandards) == 0 {
		complianceSecurityProfileStandards = append(complianceSecurityProfileStandards, string(validate.ComplianceStandardNONE))
	}

	enhancedSecurityMonitoringEnabled := workspaces.EnhancedSecurityMonitoringValueDisabled
	if input[0].EnhancedSecurityMonitoringEnabled {
		enhancedSecurityMonitoringEnabled = workspaces.EnhancedSecurityMonitoringValueEnabled
	}

	return &workspaces.EnhancedSecurityComplianceDefinition{
		AutomaticClusterUpdate: &workspaces.AutomaticClusterUpdateDefinition{
			Value: &automaticClusterUpdateEnabled,
		},
		EnhancedSecurityMonitoring: &workspaces.EnhancedSecurityMonitoringDefinition{
			Value: &enhancedSecurityMonitoringEnabled,
		},
		ComplianceSecurityProfile: &workspaces.ComplianceSecurityProfileDefinition{
			Value:               &complianceSecurityProfileEnabled,
			ComplianceStandards: &complianceSecurityProfileStandards,
		},
	}
}

func (DatabricksServerlessWorkspaceResource) flattenDatabricksServerlessWorkspaceEnhancedSecurityComplianceDefinition(input *workspaces.EnhancedSecurityComplianceDefinition) []EnhancedSecurityComplianceModel {
	if input == nil {
		return []EnhancedSecurityComplianceModel{}
	}

	enhancedSecurityCompliance := make([]EnhancedSecurityComplianceModel, 1)

	if automaticClusterUpdate := input.AutomaticClusterUpdate; automaticClusterUpdate != nil {
		if value := automaticClusterUpdate.Value; value != nil {
			enhancedSecurityCompliance[0].AutomaticClusterUpdateEnabled = pointer.From(value) == workspaces.AutomaticClusterUpdateValueEnabled
		}
	}

	if enhancedSecurityMonitoring := input.EnhancedSecurityMonitoring; enhancedSecurityMonitoring != nil {
		if value := input.EnhancedSecurityMonitoring.Value; value != nil {
			enhancedSecurityCompliance[0].EnhancedSecurityMonitoringEnabled = pointer.From(value) == workspaces.EnhancedSecurityMonitoringValueEnabled
		}
	}

	if complianceSecurityProfile := input.ComplianceSecurityProfile; complianceSecurityProfile != nil {
		if value := complianceSecurityProfile.Value; value != nil {
			enhancedSecurityCompliance[0].ComplianceSecurityProfileEnabled = pointer.From(value) == workspaces.ComplianceSecurityProfileValueEnabled
		}

		complianceSecurityProfileStandards := make([]string, 0)
		for _, complianceStandard := range pointer.From(complianceSecurityProfile.ComplianceStandards) {
			if complianceStandard == string(validate.ComplianceStandardNONE) {
				continue
			}

			complianceSecurityProfileStandards = append(complianceSecurityProfileStandards, complianceStandard)
		}

		enhancedSecurityCompliance[0].ComplianceSecurityProfileStandards = complianceSecurityProfileStandards
	}

	return enhancedSecurityCompliance
}

func (DatabricksServerlessWorkspaceResource) databricksServerlessWorkspaceEnhancedSecurityComplianceConstraint() []string {
	return []string{
		"enhanced_security_compliance.0.automatic_cluster_update_enabled",
		"enhanced_security_compliance.0.compliance_security_profile_enabled",
		"enhanced_security_compliance.0.enhanced_security_monitoring_enabled",
	}
}
