// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/fileservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var storageAccountSharePropertiesResourceName = "azurerm_storage_account_share_properties"

func resourceStorageAccountShareProperties() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceStorageAccountSharePropertiesCreate,
		Read:   resourceStorageAccountSharePropertiesRead,
		Update: resourceStorageAccountSharePropertiesUpdate,
		Delete: resourceStorageAccountSharePropertiesDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseStorageAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"cors_rule": helpers.SchemaStorageAccountCorsRule(true),

			"retention_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"days": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      7,
							ValidateFunc: validation.IntBetween(1, 365),
						},
					},
				},
			},

			"smb": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authentication_types": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Kerberos",
									"NTLMv2",
								}, false),
							},
						},

						"channel_encryption_type": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"AES-128-CCM",
									"AES-128-GCM",
									"AES-256-GCM",
								}, false),
							},
						},

						"kerberos_ticket_encryption_type": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"AES-256",
									"RC4-HMAC",
								}, false),
							},
						},

						"multichannel_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"versions": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"SMB2.1",
									"SMB3.0",
									"SMB3.1.1",
								}, false),
							},
						},
					},
				},
			},
		},
	}

	return resource
}

func resourceStorageAccountSharePropertiesCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:CREATE] Calling 'client.GetProperties' for %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	// TODO: add import check here...
	// if !response.WasNotFound(existing.HttpResponse) {
	// 	return tf.ImportAsExistsError(storageAccountResourceName, id.ID())
	// }

	accountTier := pointer.From(model.Sku.Tier)
	accountKind := pointer.From(model.Kind)
	replicationType := strings.ToUpper(strings.Split(string(model.Sku.Name), "_")[1])
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, replicationType)

	if !supportLevel.supportShare {
		return fmt.Errorf("%q are not supported for account kind %q in sku tier %q", storageAccountSharePropertiesResourceName, accountKind, accountTier)
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'storageClient.FindAccount' for %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	accountDetails, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if accountDetails == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	// NOTE: Wait for the data plane share container to become available...
	log.Printf("[DEBUG] [%s:CREATE] Calling 'custompollers.NewDataPlaneFileShareAvailabilityPoller' building File Share Poller: %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	pollerType, err := custompollers.NewDataPlaneFileShareAvailabilityPoller(storageClient, accountDetails)
	if err != nil {
		return fmt.Errorf("building File Share Poller: %+v", err)
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'poller.PollUntilDone' waiting for the File Service to become available: %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for the File Service to become available: %+v", err)
	}

	sharePayload := expandAccountShareResourceProperties(d)

	// The API complains if any multichannel info is sent on non premium file shares. Even if multichannel is set to false
	if accountTier != storageaccounts.SkuTierPremium && sharePayload.Properties != nil && sharePayload.Properties.ProtocolSettings != nil {
		// Error if the user has tried to enable multichannel on a standard tier storage account
		smb := sharePayload.Properties.ProtocolSettings.Smb
		if smb != nil && smb.Multichannel != nil {
			if smb.Multichannel.Enabled != nil && *smb.Multichannel.Enabled {
				return fmt.Errorf("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
			}

			sharePayload.Properties.ProtocolSettings.Smb.Multichannel = nil
		}
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'storageClient.ResourceManager.FileService.SetServiceProperties': %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	if _, err = storageClient.ResourceManager.FileService.SetServiceProperties(ctx, *id, sharePayload); err != nil {
		return fmt.Errorf("creating %q: %+v", storageAccountSharePropertiesResourceName, err)
	}

	d.SetId(id.ID())

	return resourceStorageAccountSharePropertiesRead(d, meta)
}

func resourceStorageAccountSharePropertiesUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:UPDATE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	accountTier := pointer.From(model.Sku.Tier)
	accountKind := pointer.From(model.Kind)
	replicationType := strings.ToUpper(strings.Split(string(model.Sku.Name), "_")[1])
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, replicationType)

	if d.HasChange("cors_rule") || d.HasChange("retention_policy") || d.HasChange("smb") {
		if !supportLevel.supportShare {
			return fmt.Errorf("%q are not supported for account kind %q in sku tier %q", storageAccountSharePropertiesResourceName, accountKind, accountTier)
		}

		sharePayload := expandAccountShareResourceProperties(d)

		// The API complains if any multichannel info is sent on non premium file shares. Even if multichannel is set to false
		if accountTier != storageaccounts.SkuTierPremium && sharePayload.Properties != nil && sharePayload.Properties.ProtocolSettings != nil {
			// Error if the user has tried to enable multichannel on a standard tier storage account
			smb := sharePayload.Properties.ProtocolSettings.Smb
			if smb != nil && smb.Multichannel != nil {
				if smb.Multichannel.Enabled != nil && *smb.Multichannel.Enabled {
					return fmt.Errorf("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
				}

				sharePayload.Properties.ProtocolSettings.Smb.Multichannel = nil
			}
		}

		log.Printf("[DEBUG] [%s:UPDATE] Calling 'storageClient.ResourceManager.FileService.SetServiceProperties': %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
		if _, err = storageClient.ResourceManager.FileService.SetServiceProperties(ctx, *id, sharePayload); err != nil {
			return fmt.Errorf("updating File Share Properties for %s: %+v", *id, err)
		}
	}

	return resourceStorageAccountSharePropertiesRead(d, meta)
}

func resourceStorageAccountSharePropertiesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] [%s:READ] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := resp.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	// we then need to find the storage account
	log.Printf("[DEBUG] [%s:READ] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	log.Printf("[DEBUG] [%s:READ] Calling 'storageClient.ResourceManager.FileService.GetServiceProperties': %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	shareProperties, err := storageClient.ResourceManager.FileService.GetServiceProperties(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving storage account share properties for %s: %+v", *id, err)
	}

	if props := shareProperties.Model.Properties; props != nil {
		if err := d.Set("cors_rule", flattenAccountSharePropertiesCorsRule(props.Cors)); err != nil {
			return fmt.Errorf("setting 'cors_rule': %+v", err)
		}

		if err := d.Set("retention_policy", flattenAccountShareDeleteRetentionPolicy(props.ShareDeleteRetentionPolicy)); err != nil {
			return fmt.Errorf("setting 'retention_policy': %+v", err)
		}

		if err := d.Set("smb", flattenAccountSharePropertiesSMB(props.ProtocolSettings)); err != nil {
			return fmt.Errorf("setting 'smb': %+v", err)
		}
	}

	return nil
}

func resourceStorageAccountSharePropertiesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:DELETE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	// NOTE: Call expand with an empty interface to get an
	// unconfigured block back from the function...
	shareProperties := defaultShareProperties()

	log.Printf("[DEBUG] [%s:DELETE] Calling 'storageClient.ResourceManager.FileService.SetServiceProperties': %s", strings.ToUpper(storageAccountSharePropertiesResourceName), id)
	if _, err = storageClient.ResourceManager.FileService.SetServiceProperties(ctx, *id, shareProperties); err != nil {
		return fmt.Errorf("deleting File Share Properties for %s: %+v", *id, err)
	}

	return nil
}

func expandAccountShareResourceProperties(d *pluginsdk.ResourceData) fileservice.FileServiceProperties {
	props := fileservice.FileServiceProperties{
		Properties: &fileservice.FileServicePropertiesProperties{
			Cors: &fileservice.CorsRules{
				CorsRules: &[]fileservice.CorsRule{},
			},
			ShareDeleteRetentionPolicy: &fileservice.DeleteRetentionPolicy{
				Enabled: pointer.To(false),
			},
			ProtocolSettings: &fileservice.ProtocolSettings{
				Smb: &fileservice.SmbSetting{
					AuthenticationMethods:    pointer.To(""),
					ChannelEncryption:        pointer.To(""),
					KerberosTicketEncryption: pointer.To(""),
					Versions:                 pointer.To(""),
					Multichannel:             nil,
				},
			},
		},
	}

	if corsRule := d.Get("cors_rule").([]interface{}); len(corsRule) > 0 {
		props.Properties.Cors = expandAccountSharePropertiesCorsRule(corsRule)
	}

	if retentionPolicy := d.Get("retention_policy").([]interface{}); len(retentionPolicy) > 0 {
		props.Properties.ShareDeleteRetentionPolicy = expandAccountShareDeleteRetentionPolicy(retentionPolicy)
	}

	if smb := d.Get("smb").([]interface{}); len(smb) > 0 {
		props.Properties.ProtocolSettings = &fileservice.ProtocolSettings{
			Smb: expandAccountSharePropertiesSMB(smb),
		}
	}

	return props
}

// TODO: Remove in v5.0, this is only here for legacy support of existing Storage Accounts...
func expandAccountShareProperties(input []interface{}) fileservice.FileServiceProperties {
	props := defaultShareProperties()

	if len(input) > 0 && input[0] != nil {
		v := input[0].(map[string]interface{})

		props.Properties.ShareDeleteRetentionPolicy = expandAccountShareDeleteRetentionPolicy(v["retention_policy"].([]interface{}))

		props.Properties.Cors = expandAccountSharePropertiesCorsRule(v["cors_rule"].([]interface{}))

		props.Properties.ProtocolSettings = &fileservice.ProtocolSettings{
			Smb: expandAccountSharePropertiesSMB(v["smb"].([]interface{})),
		}
	}

	return props
}

func flattenAccountShareProperties(input *fileservice.FileServiceProperties) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if props := input.Properties; props != nil {
			output = append(output, map[string]interface{}{
				"cors_rule":        flattenAccountSharePropertiesCorsRule(props.Cors),
				"retention_policy": flattenAccountShareDeleteRetentionPolicy(props.ShareDeleteRetentionPolicy),
				"smb":              flattenAccountSharePropertiesSMB(props.ProtocolSettings),
			})
		}
	}

	return output
}

func expandAccountSharePropertiesCorsRule(input []interface{}) *fileservice.CorsRules {
	blobCorsRules := fileservice.CorsRules{}

	if len(input) > 0 {
		corsRules := make([]fileservice.CorsRule, 0)
		for _, raw := range input {
			item := raw.(map[string]interface{})

			allowedMethods := make([]fileservice.AllowedMethods, 0)
			for _, val := range *utils.ExpandStringSlice(item["allowed_methods"].([]interface{})) {
				allowedMethods = append(allowedMethods, fileservice.AllowedMethods(val))
			}
			corsRules = append(corsRules, fileservice.CorsRule{
				AllowedHeaders:  *utils.ExpandStringSlice(item["allowed_headers"].([]interface{})),
				AllowedMethods:  allowedMethods,
				AllowedOrigins:  *utils.ExpandStringSlice(item["allowed_origins"].([]interface{})),
				ExposedHeaders:  *utils.ExpandStringSlice(item["exposed_headers"].([]interface{})),
				MaxAgeInSeconds: int64(item["max_age_in_seconds"].(int)),
			})
		}
		blobCorsRules.CorsRules = &corsRules
	}

	return &blobCorsRules
}

func flattenAccountSharePropertiesCorsRule(input *fileservice.CorsRules) []interface{} {
	corsRules := make([]interface{}, 0)

	if input == nil || input.CorsRules == nil {
		return corsRules
	}

	for _, corsRule := range *input.CorsRules {
		corsRules = append(corsRules, map[string]interface{}{
			"allowed_headers":    corsRule.AllowedHeaders,
			"allowed_methods":    corsRule.AllowedMethods,
			"allowed_origins":    corsRule.AllowedOrigins,
			"exposed_headers":    corsRule.ExposedHeaders,
			"max_age_in_seconds": int(corsRule.MaxAgeInSeconds),
		})
	}

	return corsRules
}

func expandAccountShareDeleteRetentionPolicy(input []interface{}) *fileservice.DeleteRetentionPolicy {
	result := fileservice.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &fileservice.DeleteRetentionPolicy{
		Enabled: pointer.To(true),
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func flattenAccountShareDeleteRetentionPolicy(input *fileservice.DeleteRetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if enabled := input.Enabled; enabled != nil && *enabled {
			days := 0
			if input.Days != nil {
				days = int(*input.Days)
			}

			output = append(output, map[string]interface{}{
				"days": days,
			})
		}
	}

	return output
}

func expandAccountSharePropertiesSMB(input []interface{}) *fileservice.SmbSetting {
	if len(input) == 0 || input[0] == nil {
		return &fileservice.SmbSetting{
			AuthenticationMethods:    pointer.To(""),
			ChannelEncryption:        pointer.To(""),
			KerberosTicketEncryption: pointer.To(""),
			Versions:                 pointer.To(""),
			Multichannel:             nil,
		}
	}

	v := input[0].(map[string]interface{})

	return &fileservice.SmbSetting{
		AuthenticationMethods:    utils.ExpandStringSliceWithDelimiter(v["authentication_types"].(*pluginsdk.Set).List(), ";"),
		ChannelEncryption:        utils.ExpandStringSliceWithDelimiter(v["channel_encryption_type"].(*pluginsdk.Set).List(), ";"),
		KerberosTicketEncryption: utils.ExpandStringSliceWithDelimiter(v["kerberos_ticket_encryption_type"].(*pluginsdk.Set).List(), ";"),
		Versions:                 utils.ExpandStringSliceWithDelimiter(v["versions"].(*pluginsdk.Set).List(), ";"),
		Multichannel: &fileservice.Multichannel{
			Enabled: pointer.To(v["multichannel_enabled"].(bool)),
		},
	}
}

func flattenAccountSharePropertiesSMB(input *fileservice.ProtocolSettings) []interface{} {
	if input == nil || input.Smb == nil {
		return []interface{}{}
	}

	versions := make([]interface{}, 0)
	if input.Smb.Versions != nil {
		versions = utils.FlattenStringSliceWithDelimiter(input.Smb.Versions, ";")
	}

	authenticationMethods := make([]interface{}, 0)
	if input.Smb.AuthenticationMethods != nil {
		authenticationMethods = utils.FlattenStringSliceWithDelimiter(input.Smb.AuthenticationMethods, ";")
	}

	kerberosTicketEncryption := make([]interface{}, 0)
	if input.Smb.KerberosTicketEncryption != nil {
		kerberosTicketEncryption = utils.FlattenStringSliceWithDelimiter(input.Smb.KerberosTicketEncryption, ";")
	}

	channelEncryption := make([]interface{}, 0)
	if input.Smb.ChannelEncryption != nil {
		channelEncryption = utils.FlattenStringSliceWithDelimiter(input.Smb.ChannelEncryption, ";")
	}

	multichannelEnabled := false
	if input.Smb.Multichannel != nil && input.Smb.Multichannel.Enabled != nil {
		multichannelEnabled = *input.Smb.Multichannel.Enabled
	}

	if len(versions) == 0 && len(authenticationMethods) == 0 && len(kerberosTicketEncryption) == 0 && len(channelEncryption) == 0 && (input.Smb.Multichannel == nil || input.Smb.Multichannel.Enabled == nil) {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"authentication_types":            authenticationMethods,
			"channel_encryption_type":         channelEncryption,
			"kerberos_ticket_encryption_type": kerberosTicketEncryption,
			"multichannel_enabled":            multichannelEnabled,
			"versions":                        versions,
		},
	}
}

func defaultShareProperties() fileservice.FileServiceProperties {
	return fileservice.FileServiceProperties{
		Properties: &fileservice.FileServicePropertiesProperties{
			Cors: &fileservice.CorsRules{
				CorsRules: &[]fileservice.CorsRule{},
			},
			ShareDeleteRetentionPolicy: &fileservice.DeleteRetentionPolicy{
				Enabled: pointer.To(false),
			},
		},
	}
}
