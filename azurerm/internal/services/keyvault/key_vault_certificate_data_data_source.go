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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"golang.org/x/crypto/pkcs12"
)

func dataSourceArmKeyVaultCertificateData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKeyVaultCertificateDataRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateKeyVaultChildName,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Computed

			"certificate_hex": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_expires": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmKeyVaultCertificateDataRead(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)
	version := d.Get("version").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Key %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	cert, err := client.GetCertificate(ctx, keyVaultBaseUri, name, version)
	if err != nil {
		if utils.ResponseWasNotFound(cert.Response) {
			log.Printf("[DEBUG] Certificate %q was not found in Key Vault at URI %q - removing from state", name, keyVaultBaseUri)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Key Vault Certificate: %+v", err)
	}

	if cert.ID == nil || *cert.ID == "" {
		return fmt.Errorf("failure reading Key Vault Certificate ID for %q", name)
	}

	d.SetId(*cert.ID)

	id, err := azure.ParseKeyVaultChildID(*cert.ID)
	if err != nil {
		return err
	}

	d.Set("name", id.Name)

	d.Set("version", id.Version)

	certificateData := ""
	if contents := cert.Cer; contents != nil {
		certificateData = strings.ToUpper(hex.EncodeToString(*contents))
	}
	d.Set("certificate_hex", certificateData)

	timeString, err := cert.Attributes.Expires.MarshalText()
	if err != nil {
		return fmt.Errorf("Error parsing expiry time of certificate: %+v", err)
	}

	t, err := time.Parse(time.RFC3339, string(timeString))
	if err != nil {
		return fmt.Errorf("Error converting text to Time struct: %+v", err)
	}

	d.Set("certificate_expires", t.Format(time.RFC3339))

	// Get PFX
	pfx, err := client.GetSecret(ctx, id.KeyVaultBaseUrl, id.Name, id.Version)
	if err != nil {
		return fmt.Errorf("Error retrieving cert from keyvault: %+v", err)
	}
	pfxBytes, err := base64.StdEncoding.DecodeString(*pfx.Value)
	if err != nil {
		return fmt.Errorf("Error decoding base64 certificate: %+v", err)
	}
	pfxKey, pfxCert, err := pkcs12.Decode(pfxBytes, "")
	if err != nil {
		return fmt.Errorf("Error decoding PFX cert: %+v", err)
	}
	keyX509, err := x509.MarshalPKCS8PrivateKey(pfxKey)
	if err != nil {
		return fmt.Errorf("Error reading key from PFX cert: %+v", err)
	}

	// Encode Key and PEM
	keyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyX509,
	}
	var keyPEM bytes.Buffer
	err = pem.Encode(&keyPEM, keyBlock)
	if err != nil {
		return fmt.Errorf("Error encoding Key Vault Certificate Key: %+v", err)
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: pfxCert.Raw,
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, certBlock)
	if err != nil {
		return fmt.Errorf("Error encoding Key Vault Certificate PEM: %+v", err)
	}

	d.Set("certificate_pem", certPEM.String())
	d.Set("certificate_key", keyPEM.String())

	return tags.FlattenAndSet(d, cert.Tags)
}
