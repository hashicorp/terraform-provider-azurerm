// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetAppAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetAppAccountCreate,
		Read:   resourceNetAppAccountRead,
		Update: resourceNetAppAccountUpdate,
		Delete: resourceNetAppAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := netappaccounts.ParseNetAppAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.AccountName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

			"active_directory": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dns_servers": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.IPv4Address,
							},
						},
						"domain": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[(\da-zA-Z-).]{1,255}$`),
								`The domain name must end with a letter or number before dot and start with a letter or number after dot and can not be longer than 255 characters in length.`,
							),
						},
						"smb_server_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[\da-zA-Z]{1,10}$`),
								`The smb server name can not be longer than 10 characters in length.`,
							),
						},
						"username": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"password": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"organizational_unit": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"encryption": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				MaxItems:   1,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceNetAppAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := netappaccounts.NewNetAppAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.AccountsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_netapp_account", id.ID())
		}
	}

	accountParameters := netappaccounts.NetAppAccount{
		Location:   azure.NormalizeLocation(d.Get("location").(string)),
		Properties: &netappaccounts.AccountProperties{},
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	activeDirectoryRaw := d.Get("active_directory")
	if activeDirectoryRaw != nil {
		activeDirectories := activeDirectoryRaw.([]interface{})
		activeDirectoriesExpanded := expandNetAppActiveDirectories(activeDirectories)
		if len(pointer.From(activeDirectoriesExpanded)) > 0 {
			accountParameters.Properties.ActiveDirectories = activeDirectoriesExpanded
		}
	}

	anfAccountIdentityRaw := d.Get("identity")
	if anfAccountIdentityRaw != nil {
		anfAccountIdentity := anfAccountIdentityRaw.([]interface{})
		anfAccountIdentityExpanded, err := identity.ExpandLegacySystemAndUserAssignedMap(anfAccountIdentity)
		if err != nil {
			return err
		}
		if anfAccountIdentity != nil {
			accountParameters.Identity = anfAccountIdentityExpanded
		}
	}

	encryptionRaw := d.Get("encryption")
	if encryptionRaw != nil {
		encryption := encryptionRaw.([]interface{})
		encryptionExpanded, err := expandEncryption(ctx, encryption, keyVaultsClient, resourcesClient, expandAnfAccountIdentityResourceId(accountParameters.Identity))
		if err != nil {
			return err
		}
		if encryptionExpanded != nil {
			accountParameters.Properties.Encryption = encryptionExpanded
		}
	}

	if err := client.AccountsCreateOrUpdateThenPoll(ctx, id, accountParameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Wait for account to complete create
	if err := waitForAccountCreateOrUpdate(ctx, client, id); err != nil {
		return err
	}

	d.SetId(id.ID())
	return resourceNetAppAccountRead(d, meta)
}

func resourceNetAppAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	shouldUpdate := false
	update := netappaccounts.NetAppAccountPatch{
		Properties: &netappaccounts.AccountProperties{},
	}

	if d.HasChange("active_directory") {
		shouldUpdate = true
		activeDirectoriesRaw := d.Get("active_directory").([]interface{})
		activeDirectories := expandNetAppActiveDirectories(activeDirectoriesRaw)
		update.Properties.ActiveDirectories = activeDirectories
	}

	if d.HasChange("tags") {
		shouldUpdate = true
		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	}

	if d.HasChange("identity") {
		shouldUpdate = true
		anfAccountIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		update.Identity = anfAccountIdentity
	}

	if shouldUpdate {
		if err = client.AccountsUpdateThenPoll(ctx, *id, update); err != nil {
			return fmt.Errorf("updating %s: %+v", id.ID(), err)
		}

		// Wait for account to complete update
		if err = waitForAccountCreateOrUpdate(ctx, client, *id); err != nil {
			return err
		}
	}

	return resourceNetAppAccountRead(d, meta)
}

func resourceNetAppAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AccountsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NetAppAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		anfAccountIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", anfAccountIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("encryption", flattenEncryption(model.Properties.Encryption)); err != nil {
			return fmt.Errorf("setting `encryption`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceNetAppAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	if err := client.AccountsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandNetAppActiveDirectories(input []interface{}) *[]netappaccounts.ActiveDirectory {
	results := make([]netappaccounts.ActiveDirectory, 0)
	if input == nil {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		dns := strings.Join(*utils.ExpandStringSlice(v["dns_servers"].([]interface{})), ",")

		result := netappaccounts.ActiveDirectory{
			Dns:                utils.String(dns),
			Domain:             utils.String(v["domain"].(string)),
			OrganizationalUnit: utils.String(v["organizational_unit"].(string)),
			Password:           utils.String(v["password"].(string)),
			SmbServerName:      utils.String(v["smb_server_name"].(string)),
			Username:           utils.String(v["username"].(string)),
		}

		results = append(results, result)
	}
	return &results
}

func expandEncryption(ctx context.Context, input []interface{}, keyVaultsClient *keyVaultClient.Client, resourcesClient *resourcesClient.Client, anfAccountIdentityResourceId string) (*netappaccounts.AccountEncryption, error) {
	defaultEnc := netappaccounts.AccountEncryption{
		KeySource: pointer.To(netappaccounts.KeySourceMicrosoftPointNetApp),
	}

	if len(input) == 0 || input[0] == nil {
		return &defaultEnc, nil
	}

	v := input[0].(map[string]interface{})

	keyId, err := keyVaultParse.ParseOptionallyVersionedNestedKeyID(v["key_vault_key_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
	}

	keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyId.KeyVaultBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieving the resource id the key vault at url %q: %s", keyId.KeyVaultBaseUrl, err)
	}

	parsedKeyVaultID, err := commonids.ParseKeyVaultID(pointer.From(keyVaultID))
	if err != nil {
		return nil, err
	}

	encryptionIdentity := &netappaccounts.EncryptionIdentity{}
	if anfAccountIdentityResourceId != "" {
		encryptionIdentity = &netappaccounts.EncryptionIdentity{
			UserAssignedIdentity: pointer.To(anfAccountIdentityResourceId),
		}
	}

	encryptionProperty := netappaccounts.AccountEncryption{
		Identity:  encryptionIdentity,
		KeySource: pointer.To(netappaccounts.KeySourceMicrosoftPointKeyVault),
		KeyVaultProperties: &netappaccounts.KeyVaultProperties{
			KeyName:            keyId.Name,
			KeyVaultUri:        keyId.KeyVaultBaseUrl,
			KeyVaultResourceId: parsedKeyVaultID.ID(),
		},
	}

	return &encryptionProperty, nil
}

func expandAnfAccountIdentityResourceId(anfAccountIdentity *identity.LegacySystemAndUserAssignedMap) string {
	if anfAccountIdentity == nil || anfAccountIdentity.Type != identity.TypeUserAssigned {
		return ""
	}

	identityIds := anfAccountIdentity.IdentityIds
	if len(identityIds) == 0 {
		return ""
	}

	keys := make([]string, 0, len(identityIds))
	for key := range identityIds {
		keys = append(keys, key)
	}

	if len(keys) > 0 {
		return keys[0]
	}

	return ""
}

func flattenEncryption(encryptionProperties *netappaccounts.AccountEncryption) []interface{} {
	if encryptionProperties == nil || *encryptionProperties.KeySource == netappaccounts.KeySourceMicrosoftPointNetApp {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"key_vault_key_id": encryptionProperties.KeyVaultProperties.KeyVaultUri + "keys/" + encryptionProperties.KeyVaultProperties.KeyName,
		},
	}
}

func waitForAccountCreateOrUpdate(ctx context.Context, client *netappaccounts.NetAppAccountsClient, id netappaccounts.NetAppAccountId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404"},
		Target:                    []string{"200", "202"},
		Refresh:                   netappAccountStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	return nil
}

func netappAccountStateRefreshFunc(ctx context.Context, client *netappaccounts.NetAppAccountsClient, id netappaccounts.NetAppAccountId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.AccountsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %s", id.ID(), err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}
