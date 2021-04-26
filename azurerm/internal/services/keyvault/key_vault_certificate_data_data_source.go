package keyvault

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"golang.org/x/crypto/pkcs12"
)

func dataSourceKeyVaultCertificateData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKeyVaultCertificateDataRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NestedItemName,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.VaultID,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Computed

			"hex": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"pem": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"expires": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmKeyVaultCertificateDataRead(d *schema.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := parse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}
	version := d.Get("version").(string)

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Key %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	cert, err := client.GetCertificate(ctx, *keyVaultBaseUri, name, version)
	if err != nil {
		if utils.ResponseWasNotFound(cert.Response) {
			log.Printf("[DEBUG] Certificate %q was not found in Key Vault at URI %q - removing from state", name, *keyVaultBaseUri)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Key Vault Certificate: %+v", err)
	}

	if cert.ID == nil || *cert.ID == "" {
		return fmt.Errorf("failure reading Key Vault Certificate ID for %q", name)
	}

	d.SetId(*cert.ID)

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", id.Name)

	d.Set("version", id.Version)

	certificateData := ""
	if contents := cert.Cer; contents != nil {
		certificateData = strings.ToUpper(hex.EncodeToString(*contents))
	}
	d.Set("hex", certificateData)

	timeString, err := cert.Attributes.Expires.MarshalText()
	if err != nil {
		return fmt.Errorf("parsing expiry time of certificate: %+v", err)
	}

	t, err := time.Parse(time.RFC3339, string(timeString))
	if err != nil {
		return fmt.Errorf("converting text to Time struct: %+v", err)
	}

	d.Set("expires", t.Format(time.RFC3339))

	// Get PFX
	pfx, err := client.GetSecret(ctx, id.KeyVaultBaseUrl, id.Name, id.Version)
	if err != nil {
		return fmt.Errorf("retrieving certificate %q from keyvault: %+v", id.Name, err)
	}

	pfxBytes, err := base64.StdEncoding.DecodeString(*pfx.Value)
	if err != nil {
		return fmt.Errorf("decoding base64 certificate (%q): %+v", id.Name, err)
	}

	// note PFX passwords are set to an empty string in Key Vault, this include password protected PFX uploads.
	pfxKey, pfxCert, err := pkcs12.Decode(pfxBytes, "")
	if err != nil {
		return fmt.Errorf("decoding certificate (%q): %+v", id.Name, err)
	}

	keyX509, err := x509.MarshalPKCS8PrivateKey(pfxKey)
	if err != nil {
		return fmt.Errorf("reading key from certificate (%q): %+v", id.Name, err)
	}

	// Encode Key and PEM
	keyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyX509,
	}

	var keyPEM bytes.Buffer
	err = pem.Encode(&keyPEM, keyBlock)
	if err != nil {
		return fmt.Errorf("encoding Key Vault Certificate Key: %+v", err)
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: pfxCert.Raw,
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, certBlock)
	if err != nil {
		return fmt.Errorf("encoding Key Vault Certificate PEM: %+v", err)
	}

	d.Set("pem", certPEM.String())
	d.Set("key", keyPEM.String())

	return tags.FlattenAndSet(d, cert.Tags)
}
