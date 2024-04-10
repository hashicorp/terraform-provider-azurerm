// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"golang.org/x/crypto/pkcs12"
)

func dataSourceKeyVaultCertificateData() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmKeyVaultCertificateDataRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.NestedItemName,
			},

			"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

			"version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			// Computed

			"hex": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"pem": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"expires": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"not_before": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"certificates_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmKeyVaultCertificateDataRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
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
			return fmt.Errorf("the Certificate %q was not found in Key Vault at URI %q", name, *keyVaultBaseUri)
		}

		return fmt.Errorf("reading Key Vault Certificate: %+v", err)
	}

	if cert.ID == nil || *cert.ID == "" {
		return fmt.Errorf("failure reading Key Vault Certificate ID for %q", name)
	}

	id, err := parse.ParseNestedItemID(*cert.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID())

	d.Set("name", id.Name)

	d.Set("version", id.Version)

	certificateData := ""
	if contents := cert.Cer; contents != nil {
		certificateData = strings.ToUpper(hex.EncodeToString(*contents))
	}
	d.Set("hex", certificateData)

	expireString, err := cert.Attributes.Expires.MarshalText()
	if err != nil {
		return fmt.Errorf("parsing expiry time of certificate: %+v", err)
	}

	e, err := time.Parse(time.RFC3339, string(expireString))
	if err != nil {
		return fmt.Errorf("converting text to Time struct: %+v", err)
	}

	d.Set("expires", e.Format(time.RFC3339))

	notBeforeString, err := cert.Attributes.NotBefore.MarshalText()
	if err != nil {
		return fmt.Errorf("parsing not-before time of certificate: %+v", err)
	}

	n, err := time.Parse(time.RFC3339, string(notBeforeString))
	if err != nil {
		return fmt.Errorf("converting text to Time struct: %+v", err)
	}

	d.Set("not_before", n.Format(time.RFC3339))

	// Get PFX
	pfx, err := client.GetSecret(ctx, id.KeyVaultBaseUrl, id.Name, id.Version)
	if err != nil {
		return fmt.Errorf("retrieving certificate %q from keyvault: %+v", id.Name, err)
	}

	var PEMBlocks []*pem.Block

	if *pfx.ContentType == "application/x-pkcs12" {
		bytes, err := base64.StdEncoding.DecodeString(*pfx.Value)
		if err != nil {
			return fmt.Errorf("decoding base64 certificate (%q): %+v", id.Name, err)
		}

		// note PFX passwords are set to an empty string in Key Vault, this include password protected PFX uploads.
		blocks, err := pkcs12.ToPEM(bytes, "")
		if err != nil {
			return fmt.Errorf("decoding certificate (%q): %+v", id.Name, err)
		}
		PEMBlocks = blocks
	} else {
		block, rest := pem.Decode([]byte(*pfx.Value))
		if block == nil {
			return fmt.Errorf("decoding certificate (%q): %+v", id.Name, err)
		}
		PEMBlocks = append(PEMBlocks, block)
		for len(rest) > 0 {
			block, rest = pem.Decode(rest)
			PEMBlocks = append(PEMBlocks, block)
		}
	}

	var pemKey []byte
	var pemCerts [][]byte

	for _, block := range PEMBlocks {
		if strings.Contains(block.Type, "PRIVATE KEY") {
			pemKey = block.Bytes
		}

		if strings.Contains(block.Type, "CERTIFICATE") {
			log.Printf("[DEBUG] Adding Cerrtificate block")
			pemCerts = append(pemCerts, block.Bytes)
		}
	}

	var privateKey interface{}

	if *pfx.ContentType == "application/x-pkcs12" {
		rsakey, err := x509.ParsePKCS1PrivateKey(pemKey)
		if err != nil {
			// try to parse as a EC key
			eckey, err := x509.ParseECPrivateKey(pemKey)
			if err != nil {
				return fmt.Errorf("decoding private key: not RSA or ECDSA type (%q): %+v", id.Name, err)
			}
			privateKey = eckey
		} else {
			privateKey = rsakey
		}
	} else {
		pkey, err := x509.ParsePKCS8PrivateKey(pemKey)
		if err != nil {
			return fmt.Errorf("decoding PKCS8 RSA private key (%q): %+v", id.Name, err)
		}
		privateKey = pkey
	}

	var keyX509 []byte
	var pemKeyHeader string
	if privateKey != nil {
		switch v := privateKey.(type) {
		case *ecdsa.PrivateKey:
			keyX509, err = x509.MarshalECPrivateKey(privateKey.(*ecdsa.PrivateKey))
			if err != nil {
				return fmt.Errorf("marshalling private key type %+v (%q): %+v", v, id.Name, err)
			}
			pemKeyHeader = "EC PRIVATE KEY"
		case *rsa.PrivateKey:
			keyX509 = x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey))
			pemKeyHeader = "RSA PRIVATE KEY"
		default:
			return fmt.Errorf("marshalling private key type %+v (%q): key type is not supported", v, id.Name)
		}
	}

	// Encode Key and PEM
	keyBlock := &pem.Block{
		Type:  pemKeyHeader,
		Bytes: keyX509,
	}

	var keyPEM bytes.Buffer
	err = pem.Encode(&keyPEM, keyBlock)
	if err != nil {
		return fmt.Errorf("encoding Key Vault Certificate Key: %+v", err)
	}

	certs := ""

	for _, pemCert := range pemCerts {
		certBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: pemCert,
		}

		var certPEM bytes.Buffer
		err = pem.Encode(&certPEM, certBlock)
		if err != nil {
			return fmt.Errorf("encoding Key Vault Certificate PEM: %+v", err)
		}
		certs += certPEM.String()
	}

	d.Set("pem", certs)
	d.Set("key", keyPEM.String())
	d.Set("certificates_count", len(pemCerts))

	return tags.FlattenAndSet(d, cert.Tags)
}
