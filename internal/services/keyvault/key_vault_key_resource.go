package keyvault

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"golang.org/x/crypto/ssh"
)

func resourceKeyVaultKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultKeyCreate,
		Read:   resourceKeyVaultKeyRead,
		Update: resourceKeyVaultKeyUpdate,
		Delete: resourceKeyVaultKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.ParseNestedItemID(id)
			return err
		}, nestedItemResourceImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			// TODO: Change this back to 5min, once https://github.com/hashicorp/terraform-provider-azurerm/issues/11059 is addressed.
			Read:   pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.VaultID,
			},

			"key_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// turns out Azure's *really* sensitive about the casing of these
				// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.EC),
					string(keyvault.ECHSM),
					string(keyvault.RSA),
					string(keyvault.RSAHSM),
				}, false),
			},

			"key_size": {
				Type:          pluginsdk.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"curve"},
			},

			"key_opts": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					// turns out Azure's *really* sensitive about the casing of these
					// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
					ValidateFunc: validation.StringInSlice([]string{
						string(keyvault.Decrypt),
						string(keyvault.Encrypt),
						string(keyvault.Sign),
						string(keyvault.UnwrapKey),
						string(keyvault.Verify),
						string(keyvault.WrapKey),
					}, false),
				},
			},

			"curve": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
					return old == "SECP256K1" && new == string(keyvault.P256K)
				},
				ValidateFunc: func() pluginsdk.SchemaValidateFunc {
					out := []string{
						string(keyvault.P256),
						string(keyvault.P256K),
						string(keyvault.P384),
						string(keyvault.P521),
					}
					return validation.StringInSlice(out, false)
				}(),
				// TODO: the curve name should probably be mandatory for EC in the future,
				// but handle the diff so that we don't break existing configurations and
				// imported EC keys
				ConflictsWith: []string{"key_size"},
			},

			"not_before_date": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"expiration_date": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			// Computed
			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"n": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"e": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"x": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"y": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_key_pem": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_key_openssh": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceKeyVaultKeyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Key creation.")
	name := d.Get("name").(string)
	keyVaultId, err := parse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Key %q vault url from id %q: %+v", name, *keyVaultId, err)
	}

	existing, err := client.GetKey(ctx, *keyVaultBaseUri, name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Key %q (Key Vault %q): %s", name, *keyVaultBaseUri, err)
		}
	}

	if existing.Key != nil && existing.Key.Kid != nil && *existing.Key.Kid != "" {
		return tf.ImportAsExistsError("azurerm_key_vault_key", *existing.Key.Kid)
	}

	keyType := d.Get("key_type").(string)
	keyOptions := expandKeyVaultKeyOptions(d)
	t := d.Get("tags").(map[string]interface{})

	// TODO: support Importing Keys once this is fixed:
	// https://github.com/Azure/azure-rest-api-specs/issues/1747
	parameters := keyvault.KeyCreateParameters{
		Kty:    keyvault.JSONWebKeyType(keyType),
		KeyOps: keyOptions,
		KeyAttributes: &keyvault.KeyAttributes{
			Enabled: utils.Bool(true),
		},

		Tags: tags.Expand(t),
	}

	if parameters.Kty == keyvault.EC || parameters.Kty == keyvault.ECHSM {
		curveName := d.Get("curve").(string)
		parameters.Curve = keyvault.JSONWebKeyCurveName(curveName)
	} else if parameters.Kty == keyvault.RSA || parameters.Kty == keyvault.RSAHSM {
		keySize, ok := d.GetOk("key_size")
		if !ok {
			return fmt.Errorf("Key size is required when creating an RSA key")
		}
		parameters.KeySize = utils.Int32(int32(keySize.(int)))
	}
	// TODO: support `oct` once this is fixed
	// https://github.com/Azure/azure-rest-api-specs/issues/1739#issuecomment-332236257

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		notBeforeUnixTime := date.UnixTime(notBeforeDate)
		parameters.KeyAttributes.NotBefore = &notBeforeUnixTime
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		expirationUnixTime := date.UnixTime(expirationDate)
		parameters.KeyAttributes.Expires = &expirationUnixTime
	}

	if resp, err := client.CreateKey(ctx, *keyVaultBaseUri, name, parameters); err != nil {
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeys && utils.ResponseWasConflict(resp.Response) {
			recoveredKey, err := client.RecoverDeletedKey(ctx, *keyVaultBaseUri, name)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Recovering Key %q with ID: %q", name, *recoveredKey.Key.Kid)
			if kid := recoveredKey.Key.Kid; kid != nil {
				stateConf := &pluginsdk.StateChangeConf{
					Pending:                   []string{"pending"},
					Target:                    []string{"available"},
					Refresh:                   keyVaultChildItemRefreshFunc(*kid),
					Delay:                     30 * time.Second,
					PollInterval:              10 * time.Second,
					ContinuousTargetOccurence: 10,
					Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
				}

				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for Key Vault Secret %q to become available: %s", name, err)
				}
				log.Printf("[DEBUG] Key %q recovered with ID: %q", name, *kid)
			}
		} else {
			return fmt.Errorf("Creating Key: %+v", err)
		}
	}

	// "" indicates the latest version
	read, err := client.GetKey(ctx, *keyVaultBaseUri, name, "")
	if err != nil {
		return err
	}

	if read.Key == nil || read.Key.Kid == nil {
		return fmt.Errorf("cannot read KeyVault Key '%s' (in key vault '%s')", name, *keyVaultBaseUri)
	}
	keyId, err := parse.ParseNestedItemID(*read.Key.Kid)
	if err != nil {
		return err
	}
	d.SetId(keyId.ID())

	return resourceKeyVaultKeyRead(d, meta)
}

func resourceKeyVaultKeyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := parse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	meta.(*clients.Client).KeyVault.AddToCache(*keyVaultId, id.KeyVaultBaseUrl)

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	keyOptions := expandKeyVaultKeyOptions(d)
	t := d.Get("tags").(map[string]interface{})

	parameters := keyvault.KeyUpdateParameters{
		KeyOps: keyOptions,
		KeyAttributes: &keyvault.KeyAttributes{
			Enabled: utils.Bool(true),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		notBeforeUnixTime := date.UnixTime(notBeforeDate)
		parameters.KeyAttributes.NotBefore = &notBeforeUnixTime
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		expirationUnixTime := date.UnixTime(expirationDate)
		parameters.KeyAttributes.Expires = &expirationUnixTime
	}

	if _, err = client.UpdateKey(ctx, id.KeyVaultBaseUrl, id.Name, "", parameters); err != nil {
		return err
	}

	return resourceKeyVaultKeyRead(d, meta)
}

func resourceKeyVaultKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
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

	resp, err := client.GetKey(ctx, id.KeyVaultBaseUrl, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Key %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("key_vault_id", *keyVaultId)
	d.Set("name", id.Name)

	if key := resp.Key; key != nil {
		d.Set("key_type", string(key.Kty))

		options := flattenKeyVaultKeyOptions(key.KeyOps)
		if err := d.Set("key_opts", options); err != nil {
			return err
		}

		d.Set("n", key.N)
		d.Set("e", key.E)
		d.Set("x", key.X)
		d.Set("y", key.Y)
		if key.N != nil {
			nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
			if err != nil {
				return fmt.Errorf("Could not decode N: %+v", err)
			}
			d.Set("key_size", len(nBytes)*8)
		}

		d.Set("curve", key.Crv)
	}

	if attributes := resp.Attributes; attributes != nil {
		if v := attributes.NotBefore; v != nil {
			d.Set("not_before_date", time.Time(*v).Format(time.RFC3339))
		}

		if v := attributes.Expires; v != nil {
			d.Set("expiration_date", time.Time(*v).Format(time.RFC3339))
		}
	}

	// Computed
	d.Set("version", id.Version)
	d.Set("versionless_id", id.VersionlessID())
	if key := resp.Key; key != nil {
		if key.Kty == keyvault.RSA || key.Kty == keyvault.RSAHSM {
			nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
			if err != nil {
				return fmt.Errorf("failed to decode N: %+v", err)
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(*key.E)
			if err != nil {
				return fmt.Errorf("failed to decode E: %+v", err)
			}
			publicKey := &rsa.PublicKey{
				N: big.NewInt(0).SetBytes(nBytes),
				E: int(big.NewInt(0).SetBytes(eBytes).Uint64()),
			}
			err = readPublicKey(d, publicKey)
			if err != nil {
				return fmt.Errorf("failed to read public key: %+v", err)
			}
		} else if key.Kty == keyvault.EC || key.Kty == keyvault.ECHSM {
			// do ec keys
			xBytes, err := base64.RawURLEncoding.DecodeString(*key.X)
			if err != nil {
				return fmt.Errorf("failed to decode X: %+v", err)
			}
			yBytes, err := base64.RawURLEncoding.DecodeString(*key.Y)
			if err != nil {
				return fmt.Errorf("failed to decode Y: %+v", err)
			}
			publicKey := &ecdsa.PublicKey{
				X: big.NewInt(0).SetBytes(xBytes),
				Y: big.NewInt(0).SetBytes(yBytes),
			}
			switch key.Crv {
			case keyvault.P256:
				publicKey.Curve = elliptic.P256()
			case keyvault.P384:
				publicKey.Curve = elliptic.P384()
			case keyvault.P521:
				publicKey.Curve = elliptic.P521()
			}
			if publicKey.Curve != nil {
				err = readPublicKey(d, publicKey)
				if err != nil {
					return fmt.Errorf("failed to read public key: %+v", err)
				}
			}
		}
	}

	d.Set("resource_id", parse.NewKeyID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, id.Name, id.Version).ID())
	d.Set("resource_versionless_id", parse.NewKeyVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, id.Name).ID())

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceKeyVaultKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
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
			log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving key vault %q properties: %+v", *keyVaultId, err)
	}

	shouldPurge := meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedKeysOnDestroy
	if shouldPurge && kv.Properties != nil && utils.NormaliseNilableBool(kv.Properties.EnablePurgeProtection) {
		log.Printf("[DEBUG] cannot purge key %q because vault %q has purge protection enabled", id.Name, keyVaultId.String())
		shouldPurge = false
	}

	description := fmt.Sprintf("Key %q (Key Vault %q)", id.Name, id.KeyVaultBaseUrl)
	deleter := deleteAndPurgeKey{
		client:      client,
		keyVaultUri: id.KeyVaultBaseUrl,
		name:        id.Name,
	}
	if err := deleteAndOptionallyPurge(ctx, description, shouldPurge, deleter); err != nil {
		return err
	}

	return nil
}

var _ deleteAndPurgeNestedItem = deleteAndPurgeKey{}

type deleteAndPurgeKey struct {
	client      *keyvault.BaseClient
	keyVaultUri string
	name        string
}

func (d deleteAndPurgeKey) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.DeleteKey(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func (d deleteAndPurgeKey) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetKey(ctx, d.keyVaultUri, d.name, "")
	return resp.Response, err
}

func (d deleteAndPurgeKey) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	return d.client.PurgeDeletedKey(ctx, d.keyVaultUri, d.name)
}

func (d deleteAndPurgeKey) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetDeletedKey(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func expandKeyVaultKeyOptions(d *pluginsdk.ResourceData) *[]keyvault.JSONWebKeyOperation {
	options := d.Get("key_opts").([]interface{})
	results := make([]keyvault.JSONWebKeyOperation, 0, len(options))

	for _, option := range options {
		results = append(results, keyvault.JSONWebKeyOperation(option.(string)))
	}

	return &results
}

func flattenKeyVaultKeyOptions(input *[]string) []interface{} {
	results := make([]interface{}, 0, len(*input))

	for _, option := range *input {
		results = append(results, option)
	}

	return results
}

// Credit to Hashicorp modified from https://github.com/hashicorp/terraform-provider-tls/blob/v3.1.0/internal/provider/util.go#L79-L105
func readPublicKey(d *pluginsdk.ResourceData, pubKey interface{}) error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key error: %s", err)
	}
	pubKeyPemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	d.Set("public_key_pem", string(pem.EncodeToMemory(pubKeyPemBlock)))

	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err == nil {
		// Not all EC types can be SSH keys, so we'll produce this only
		// if an appropriate type was selected.
		sshPubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)
		d.Set("public_key_openssh", string(sshPubKeyBytes))
	} else {
		d.Set("public_key_openssh", "")
	}
	return nil
}
