package keyvault

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/sdk/v7.3/keyvault"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKeyVaultKeyRotationPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultKeyRotationPolicyCreate,
		Read:   resourceKeyVaultKeyRotationPolicyRead,
		Update: resourceKeyVaultKeyRotationPolicyUpdate,
		Delete: resourceKeyVaultKeyRotationPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ParseNestedItemID(id) // TODO as it is not a real nested item with a version..
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read:   pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"key_resource_versionless_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.KeyVersionlessID,
			},

			"expire_after": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ISO8601Duration,
				AtLeastOneOf: []string{
					"expire_after",
					"automatic",
				},
				RequiredWith: []string{
					"expire_after",
					"notify_before_expiry",
				},
			},

			"notify_before_expiry": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ISO8601Duration,
				RequiredWith: []string{
					"expire_after",
					"notify_before_expiry",
				},
			},

			"automatic": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"time_after_creation": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.ISO8601Duration,
							AtLeastOneOf: []string{
								"automatic.0.time_after_creation",
								"automatic.0.time_before_expiry",
							},
						},
						"time_before_expiry": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.ISO8601Duration,
							AtLeastOneOf: []string{
								"automatic.0.time_after_creation",
								"automatic.0.time_before_expiry",
							},
						},
					},
				},
			},

			"resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKeyVaultKeyRotationPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Key Rotation Policy creation.")
	keyVersionlessId, err := parse.KeyVersionlessID(d.Get("key_resource_versionless_id").(string))
	if err != nil {
		return err
	}
	keyVaultId := parse.NewVaultID(keyVersionlessId.SubscriptionId, keyVersionlessId.ResourceGroup, keyVersionlessId.VaultName)

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Key %q vault url from id %q: %+v", keyVersionlessId.KeyName, keyVaultId, err)
	}

	existing, err := client.GetKeyRotationPolicy(ctx, *keyVaultBaseUri, keyVersionlessId.KeyName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Key Rotation Policy for Key %q (Key Vault %q): %s", keyVersionlessId.KeyName, *keyVaultBaseUri, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_key_vault_key_rotation_policy", *existing.ID)
	}

	if _, err := client.UpdateKeyRotationPolicy(ctx, *keyVaultBaseUri, keyVersionlessId.KeyName, expandKeyVaultKeyRotationPolicyOptions(d)); err != nil {
		return fmt.Errorf("Creating Key Rotation Policy: %+v", err)
	}

	read, err := client.GetKeyRotationPolicy(ctx, *keyVaultBaseUri, keyVersionlessId.KeyName)
	if err != nil {
		return err
	}

	d.SetId(*read.ID)

	return resourceKeyVaultKeyRotationPolicyRead(d, meta)
}

func resourceKeyVaultKeyRotationPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id()) // TODO as it is not a real nested item with a version..
	if err != nil {
		return err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}

	keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	if _, err = client.UpdateKeyRotationPolicy(ctx, id.KeyVaultBaseUrl, id.Name, expandKeyVaultKeyRotationPolicyOptions(d)); err != nil {
		return err
	}

	return resourceKeyVaultKeyRotationPolicyRead(d, meta)
}

func resourceKeyVaultKeyRotationPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id()) // TODO as it is not a real nested item with a version..
	if err != nil {
		return err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}
	keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	resp, err := client.GetKeyRotationPolicy(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Key Rotation Policy for key %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}

		return err
	}

	expiryTimeSet := false
	if resp.Attributes != nil && resp.Attributes.ExpiryTime != nil {
		d.Set("expire_after", *resp.Attributes.ExpiryTime)
		expiryTimeSet = true
	}

	if resp.LifetimeActions != nil {
		for _, ltAction := range *resp.LifetimeActions {
			action := ltAction.Action
			trigger := ltAction.Trigger
			if action == nil || trigger == nil {
				log.Printf("[DEBUG] Key Rotation Policy for key %q in Key Vault at URI %q has no trigger or action", id.Name, id.KeyVaultBaseUrl)
			}

			if strings.EqualFold(string(action.Type), string(keyvault.Rotate)) {
				autoRotation := make(map[string]interface{}, 0)
				if timeAfterCreate := trigger.TimeAfterCreate; timeAfterCreate != nil {
					autoRotation["time_after_creation"] = *timeAfterCreate
				}
				if timeBeforeExpiry := trigger.TimeBeforeExpiry; timeBeforeExpiry != nil {
					autoRotation["time_before_expiry"] = *timeBeforeExpiry
				}
				d.Set("automatic", []map[string]interface{}{autoRotation})
			}

			if strings.EqualFold(string(action.Type), string(keyvault.Notify)) && trigger.TimeBeforeExpiry != nil && expiryTimeSet {
				d.Set("notify_before_expiry", *trigger.TimeBeforeExpiry)
			}
		}
	}

	d.Set("key_resource_versionless_id", parse.NewKeyVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, id.Name).ID())
	d.Set("resource_id", parse.NewKeyRotationPolicyID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, id.Name).ID())

	return nil
}

func resourceKeyVaultKeyRotationPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id()) // TODO as it is not a real nested item with a version..
	if err != nil {
		return err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}
	keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	kv, err := keyVaultsClient.VaultsClient.Get(ctx, keyVaultId.ResourceGroup, keyVaultId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(kv.Response) {
			log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing Key Policy from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving key vault %q properties: %+v", *keyVaultId, err)
	}

	// Delete is basically emptying the thing
	parameters := keyvault.KeyRotationPolicy{
		LifetimeActions: &[]keyvault.LifetimeActions{},
		Attributes: &keyvault.KeyRotationPolicyAttributes{
			ExpiryTime: nil,
		},
	}
	if _, err = client.UpdateKeyRotationPolicy(ctx, id.KeyVaultBaseUrl, id.Name, parameters); err != nil {
		return err
	}

	return nil
}

func expandKeyVaultKeyRotationPolicyOptions(d *pluginsdk.ResourceData) keyvault.KeyRotationPolicy {
	lifetimeActions := make([]keyvault.LifetimeActions, 0)

	var expiryTime *string = nil // needs to be set to nil if not set
	if rawExpiryTime := d.Get("expire_after").(string); rawExpiryTime != "" {
		expiryTime = utils.String(rawExpiryTime)
	}

	if notifyBeforeExpiry := d.Get("notify_before_expiry").(string); notifyBeforeExpiry != "" {
		lifetimeActionNotify := keyvault.LifetimeActions{
			Trigger: &keyvault.LifetimeActionsTrigger{
				TimeBeforeExpiry: utils.String(notifyBeforeExpiry), // for notify always before expiry
			},
			Action: &keyvault.LifetimeActionsType{
				Type: keyvault.Notify,
			},
		}
		lifetimeActions = append(lifetimeActions, lifetimeActionNotify)
	}

	autoRotation := d.Get("automatic").([]interface{})
	if autoRotation != nil && len(autoRotation) == 1 && autoRotation[0] != nil {
		lifetimeActionRotate := keyvault.LifetimeActions{
			Action: &keyvault.LifetimeActionsType{
				Type: keyvault.Rotate,
			},
			Trigger: &keyvault.LifetimeActionsTrigger{},
		}
		autoRotationRaw := autoRotation[0].(map[string]interface{})

		if v := autoRotationRaw["time_after_creation"]; v != "" {
			timeAfterCreate := v.(string)
			lifetimeActionRotate.Trigger.TimeAfterCreate = &timeAfterCreate
		}

		if v := autoRotationRaw["time_before_expiry"]; v != "" {
			timeBeforeExpiry := v.(string)
			lifetimeActionRotate.Trigger.TimeBeforeExpiry = &timeBeforeExpiry
		}

		lifetimeActions = append(lifetimeActions, lifetimeActionRotate)
	}

	return keyvault.KeyRotationPolicy{
		LifetimeActions: &lifetimeActions,
		Attributes: &keyvault.KeyRotationPolicyAttributes{
			ExpiryTime: expiryTime,
		},
	}
}
