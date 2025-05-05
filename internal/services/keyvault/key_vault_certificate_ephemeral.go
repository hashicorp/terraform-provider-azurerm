// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"golang.org/x/crypto/pkcs12"
)

var _ sdk.EphemeralResource = &KeyVaultCertificateEphemeralResource{}

func NewKeyVaultCertificateEphemeralResource() ephemeral.EphemeralResource {
	return &KeyVaultCertificateEphemeralResource{}
}

type KeyVaultCertificateEphemeralResource struct {
	sdk.EphemeralResourceMetadata
}

type KeyVaultCertificateEphemeralResourceModel struct {
	Name             types.String `tfsdk:"name"`
	KeyVaultID       types.String `tfsdk:"key_vault_id"`
	Version          types.String `tfsdk:"version"`
	Hex              types.String `tfsdk:"hex"`
	Pem              types.String `tfsdk:"pem"`
	Key              types.String `tfsdk:"key"`
	ExpirationDate   types.String `tfsdk:"expiration_date"`
	NotBeforeDate    types.String `tfsdk:"not_before_date"`
	CertificateCount types.Int64  `tfsdk:"certificate_count"`
}

func (e *KeyVaultCertificateEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "azurerm_key_vault_certificate"
}

func (e *KeyVaultCertificateEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	e.Defaults(req, resp)
}

func (e *KeyVaultCertificateEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.StringIsNotEmpty,
					},
				},
			},

			"key_vault_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateKeyVaultID,
					},
				},
			},

			"version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.StringIsNotEmpty,
					},
				},
			},

			"hex": schema.StringAttribute{
				Computed: true,
			},

			"pem": schema.StringAttribute{
				Computed: true,
			},

			"key": schema.StringAttribute{
				Computed: true,
			},

			"expiration_date": schema.StringAttribute{
				Computed: true,
			},

			"not_before_date": schema.StringAttribute{
				Computed: true,
			},

			"certificate_count": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (e *KeyVaultCertificateEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	keyVaultsClient := e.Client.KeyVault
	client := e.Client.KeyVault.ManagementClient
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	var data KeyVaultCertificateEphemeralResourceModel

	if ok := e.DecodeOpen(ctx, req, resp, &data); !ok {
		return
	}

	keyVaultID, err := commonids.ParseKeyVaultID(data.KeyVaultID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "", err)
		return
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("looking up base uri for certificate %q in %s", data.Name.ValueString(), keyVaultID), err)
		return
	}

	response, err := client.GetCertificate(ctx, *keyVaultBaseUri, data.Name.ValueString(), data.Version.ValueString())
	if err != nil {
		if utils.ResponseWasNotFound(response.Response) {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("certificate %q does not exist in %s", data.Name.ValueString(), keyVaultID), err)
			return
		}
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving certificate %q from %s", data.Name.ValueString(), keyVaultID), err)
		return
	}

	id, err := parse.ParseNestedItemID(*response.ID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "", err)
		return
	}
	data.Version = types.StringValue(id.Version)

	if attributes := response.Attributes; attributes != nil {
		if expires := attributes.Expires; expires != nil {
			data.ExpirationDate = types.StringValue(time.Time(*expires).Format(time.RFC3339))
		}

		if notBefore := attributes.NotBefore; notBefore != nil {
			data.NotBeforeDate = types.StringValue(time.Time(*notBefore).Format(time.RFC3339))
		}
	}

	certificateData := ""
	if response.Cer != nil {
		certificateData = strings.ToUpper(hex.EncodeToString(*response.Cer))
	}

	data.Hex = types.StringValue(certificateData)

	pfx, err := client.GetSecret(ctx, id.KeyVaultBaseUrl, id.Name, id.Version)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving certificate %q from %s", data.Name.ValueString(), keyVaultID), err)
		return
	}

	pemBlocks := make([]*pem.Block, 0)

	if *pfx.ContentType == "application/x-pkcs12" {
		bytes, err := base64.StdEncoding.DecodeString(*pfx.Value)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("decoding base64 certificate %q", id.Name), err)
			return
		}

		blocks, err := pkcs12.ToPEM(bytes, "")
		if err != nil {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("decoding certificate %q", id.Name), err)
			return
		}
		pemBlocks = blocks
	} else {
		block, rest := pem.Decode([]byte(*pfx.Value))
		if block == nil {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("decoding %q", id.Name), err)
			return
		}
		pemBlocks = append(pemBlocks, block)
		for len(rest) > 0 {
			block, rest = pem.Decode(rest)
			pemBlocks = append(pemBlocks, block)
		}
	}

	var pemKey []byte
	var pemCerts [][]byte

	for _, block := range pemBlocks {
		if strings.Contains(block.Type, "PRIVATE KEY") {
			pemKey = block.Bytes
		}

		if strings.Contains(block.Type, "CERTIFICATE") {
			pemCerts = append(pemCerts, block.Bytes)
		}
	}

	var privateKey interface{}

	if *pfx.ContentType == "application/x-pkcs12" {
		rsakey, err := x509.ParsePKCS1PrivateKey(pemKey)
		if err != nil {
			eckey, err := x509.ParseECPrivateKey(pemKey)
			if err != nil {
				sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("decoding private key %q not RSA or ECDSA type", id.Name), err)
				return
			}
			privateKey = eckey
		} else {
			privateKey = rsakey
		}
	} else {
		pkey, err := x509.ParsePKCS8PrivateKey(pemKey)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("decoding PKCS8 RSA private key %q", id.Name), err)
			return
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
				sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("marshalling private key %q of type %+v", id.Name, v), err)
				return
			}
			pemKeyHeader = "EC PRIVATE KEY"
		case *rsa.PrivateKey:
			keyX509 = x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey))
			pemKeyHeader = "RSA PRIVATE KEY"
		default:
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("marshalling private key %q: key type %+v is not supported", id.Name, v), err)
			return
		}
	}

	keyBlock := &pem.Block{
		Type:  pemKeyHeader,
		Bytes: keyX509,
	}

	var keyPEM bytes.Buffer
	err = pem.Encode(&keyPEM, keyBlock)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("encoding key for %q", id.Name), err)
		return
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
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("encoding PEM for %q", id.Name), err)
			return
		}
		certs += certPEM.String()
	}

	data.Pem = types.StringValue(certs)
	data.Key = types.StringValue(keyPEM.String())
	data.CertificateCount = types.Int64Value(int64(len(pemCerts)))

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
