package compute

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func additionalUnattendContentSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		// whilst the SDK supports updating, the API doesn't:
		//   Code="PropertyChangeNotAllowed"
		//   Message="Changing property 'windowsConfiguration.additionalUnattendContent' is not allowed."
		//   Target="windowsConfiguration.additionalUnattendContent
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"content": {
					Type:      schema.TypeString,
					Required:  true,
					ForceNew:  true,
					Sensitive: true,
				},
				"setting": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.AutoLogon),
						string(compute.FirstLogonCommands),
					}, false),
				},
			},
		},
	}
}

func expandAdditionalUnattendContent(input []interface{}) *[]compute.AdditionalUnattendContent {
	output := make([]compute.AdditionalUnattendContent, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		output = append(output, compute.AdditionalUnattendContent{
			SettingName: compute.SettingNames(raw["setting"].(string)),
			Content:     utils.String(raw["content"].(string)),

			// no other possible values
			PassName:      compute.OobeSystem,
			ComponentName: compute.MicrosoftWindowsShellSetup,
		})
	}

	return &output
}

func flattenAdditionalUnattendContent(input *[]compute.AdditionalUnattendContent, d *schema.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	existing := make([]interface{}, 0)
	if v, ok := d.GetOk("additional_unattend_content"); ok {
		existing = v.([]interface{})
	}

	output := make([]interface{}, 0)
	for i, v := range *input {
		// content isn't returned from the API as it's sensitive so we need to look it up
		content := ""
		if len(existing) > i {
			existingVal := existing[i]
			existingRaw, ok := existingVal.(map[string]interface{})
			if ok {
				contentRaw, ok := existingRaw["content"]
				if ok {
					content = contentRaw.(string)
				}
			}
		}

		output = append(output, map[string]interface{}{
			"content": content,
			"setting": string(v.SettingName),
		})
	}

	return output
}

func bootDiagnosticsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// TODO: should this be `storage_account_endpoint`?
				"storage_account_uri": {
					Type:     schema.TypeString,
					Optional: true,
					// TODO: validation
				},
			},
		},
	}
}

func expandBootDiagnostics(input []interface{}) *compute.DiagnosticsProfile {
	if len(input) == 0 || input[0] == nil {
		return &compute.DiagnosticsProfile{
			BootDiagnostics: &compute.BootDiagnostics{
				Enabled:    utils.Bool(false),
				StorageURI: utils.String(""),
			},
		}
	}

	raw := input[0].(map[string]interface{})

	storageAccountURI := raw["storage_account_uri"].(string)

	return &compute.DiagnosticsProfile{
		BootDiagnostics: &compute.BootDiagnostics{
			Enabled:    utils.Bool(true),
			StorageURI: utils.String(storageAccountURI),
		},
	}
}

func flattenBootDiagnostics(input *compute.DiagnosticsProfile) []interface{} {
	if input == nil || input.BootDiagnostics == nil || input.BootDiagnostics.Enabled == nil || !*input.BootDiagnostics.Enabled {
		return []interface{}{}
	}

	storageAccountUri := ""
	if input.BootDiagnostics.StorageURI != nil {
		storageAccountUri = *input.BootDiagnostics.StorageURI
	}

	return []interface{}{
		map[string]interface{}{
			"storage_account_uri": storageAccountUri,
		},
	}
}

func linuxSecretSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// whilst this isn't present in the nested object it's required when this is specified
				"key_vault_id": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: azure.ValidateResourceID, // TODO: more granular validation
				},

				// whilst we /could/ flatten this to `certificate_urls` we're intentionally not to keep this
				// closer to the Windows VMSS resource, which will also take a `store` param
				"certificate": {
					Type:     schema.TypeSet,
					Required: true,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"url": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: azure.ValidateKeyVaultChildId,
							},
						},
					},
				},
			},
		},
	}
}

func expandLinuxSecrets(input []interface{}) *[]compute.VaultSecretGroup {
	output := make([]compute.VaultSecretGroup, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})

		keyVaultId := v["key_vault_id"].(string)
		certificatesRaw := v["certificate"].(*schema.Set).List()
		certificates := make([]compute.VaultCertificate, 0)
		for _, certificateRaw := range certificatesRaw {
			certificateV := certificateRaw.(map[string]interface{})

			url := certificateV["url"].(string)
			certificates = append(certificates, compute.VaultCertificate{
				CertificateURL: utils.String(url),
			})
		}

		output = append(output, compute.VaultSecretGroup{
			SourceVault: &compute.SubResource{
				ID: utils.String(keyVaultId),
			},
			VaultCertificates: &certificates,
		})
	}

	return &output
}

func flattenLinuxSecrets(input *[]compute.VaultSecretGroup) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		keyVaultId := ""
		if v.SourceVault != nil && v.SourceVault.ID != nil {
			keyVaultId = *v.SourceVault.ID
		}

		certificates := make([]interface{}, 0)

		if v.VaultCertificates != nil {
			for _, c := range *v.VaultCertificates {
				if c.CertificateURL == nil {
					continue
				}

				certificates = append(certificates, map[string]interface{}{
					"url": *c.CertificateURL,
				})
			}
		}

		output = append(output, map[string]interface{}{
			"key_vault_id": keyVaultId,
			"certificate":  certificates,
		})
	}

	return output
}

func planSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},

				"product": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},

				"publisher": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func expandPlan(input []interface{}) *compute.Plan {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &compute.Plan{
		Name:      utils.String(raw["name"].(string)),
		Product:   utils.String(raw["product"].(string)),
		Publisher: utils.String(raw["publisher"].(string)),
	}
}

func flattenPlan(input *compute.Plan) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	product := ""
	if input.Product != nil {
		product = *input.Product
	}

	publisher := ""
	if input.Publisher != nil {
		publisher = *input.Publisher
	}

	return []interface{}{
		map[string]interface{}{
			"name":      name,
			"product":   product,
			"publisher": publisher,
		},
	}
}

func sourceImageReferenceSchema(isVirtualMachine bool) *schema.Schema {
	// whilst originally I was hoping we could use the 'id' from `azurerm_platform_image' unfortunately Azure doesn't
	// like this as a value for the 'id' field:
	// Id /...../Versions/16.04.201909091 is not a valid resource reference."
	// as such the image is split into two fields (source_image_id and source_image_reference) to provide better validation
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      isVirtualMachine,
		MaxItems:      1,
		ConflictsWith: []string{"source_image_id"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"publisher": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: isVirtualMachine,
				},
				"offer": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: isVirtualMachine,
				},
				"sku": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: isVirtualMachine,
				},
				"version": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: isVirtualMachine,
				},
			},
		},
	}
}

func expandSourceImageReference(referenceInput []interface{}, imageId string) (*compute.ImageReference, error) {
	if imageId != "" {
		return &compute.ImageReference{
			ID: utils.String(imageId),
		}, nil
	}

	if len(referenceInput) == 0 {
		return nil, fmt.Errorf("Either a `source_image_id` or a `source_image_reference` block must be specified!")
	}

	raw := referenceInput[0].(map[string]interface{})
	return &compute.ImageReference{
		Publisher: utils.String(raw["publisher"].(string)),
		Offer:     utils.String(raw["offer"].(string)),
		Sku:       utils.String(raw["sku"].(string)),
		Version:   utils.String(raw["version"].(string)),
	}, nil
}

func flattenSourceImageReference(input *compute.ImageReference) []interface{} {
	// since the image id is pulled out as a separate field, if that's set we should return an empty block here
	if input == nil || input.ID != nil {
		return []interface{}{}
	}

	var publisher, offer, sku, version string

	if input.Publisher != nil {
		publisher = *input.Publisher
	}
	if input.Offer != nil {
		offer = *input.Offer
	}
	if input.Sku != nil {
		sku = *input.Sku
	}
	if input.Version != nil {
		version = *input.Version
	}

	return []interface{}{
		map[string]interface{}{
			"publisher": publisher,
			"offer":     offer,
			"sku":       sku,
			"version":   version,
		},
	}
}

func winRmListenerSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		// Whilst the SDK allows you to modify this, the API does not:
		//   Code="PropertyChangeNotAllowed"
		//   Message="Changing property 'windowsConfiguration.winRM.listeners' is not allowed."
		//   Target="windowsConfiguration.winRM.listeners"
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"protocol": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.HTTP),
						string(compute.HTTPS),
					}, false),
				},

				"certificate_url": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: azure.ValidateKeyVaultChildId,
				},
			},
		},
	}
}

func expandWinRMListener(input []interface{}) *compute.WinRMConfiguration {
	listeners := make([]compute.WinRMListener, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		listener := compute.WinRMListener{
			Protocol: compute.ProtocolTypes(raw["protocol"].(string)),
		}

		certificateUrl := raw["certificate_url"].(string)
		if certificateUrl != "" {
			listener.CertificateURL = utils.String(certificateUrl)
		}

		listeners = append(listeners, listener)
	}

	return &compute.WinRMConfiguration{
		Listeners: &listeners,
	}
}

func flattenWinRMListener(input *compute.WinRMConfiguration) []interface{} {
	if input == nil || input.Listeners == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input.Listeners {
		certificateUrl := ""
		if v.CertificateURL != nil {
			certificateUrl = *v.CertificateURL
		}

		output = append(output, map[string]interface{}{
			"certificate_url": certificateUrl,
			"protocol":        string(v.Protocol),
		})
	}

	return output
}

func windowsSecretSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// whilst this isn't present in the nested object it's required when this is specified
				"key_vault_id": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: azure.ValidateResourceID,
				},

				"certificate": {
					Type:     schema.TypeSet,
					Required: true,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"store": {
								Type:     schema.TypeString,
								Required: true,
							},
							"url": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: azure.ValidateKeyVaultChildId,
							},
						},
					},
				},
			},
		},
	}
}

func expandWindowsSecrets(input []interface{}) *[]compute.VaultSecretGroup {
	output := make([]compute.VaultSecretGroup, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})

		keyVaultId := v["key_vault_id"].(string)
		certificatesRaw := v["certificate"].(*schema.Set).List()
		certificates := make([]compute.VaultCertificate, 0)
		for _, certificateRaw := range certificatesRaw {
			certificateV := certificateRaw.(map[string]interface{})

			store := certificateV["store"].(string)
			url := certificateV["url"].(string)
			certificates = append(certificates, compute.VaultCertificate{
				CertificateStore: utils.String(store),
				CertificateURL:   utils.String(url),
			})
		}

		output = append(output, compute.VaultSecretGroup{
			SourceVault: &compute.SubResource{
				ID: utils.String(keyVaultId),
			},
			VaultCertificates: &certificates,
		})
	}

	return &output
}

func flattenWindowsSecrets(input *[]compute.VaultSecretGroup) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		keyVaultId := ""
		if v.SourceVault != nil && v.SourceVault.ID != nil {
			keyVaultId = *v.SourceVault.ID
		}

		certificates := make([]interface{}, 0)

		if v.VaultCertificates != nil {
			for _, c := range *v.VaultCertificates {
				store := ""
				if c.CertificateStore != nil {
					store = *c.CertificateStore
				}

				url := ""
				if c.CertificateURL != nil {
					url = *c.CertificateURL
				}

				certificates = append(certificates, map[string]interface{}{
					"store": store,
					"url":   url,
				})
			}
		}

		output = append(output, map[string]interface{}{
			"key_vault_id": keyVaultId,
			"certificate":  certificates,
		})
	}

	return output
}
