// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AdditionalUnattendContentSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		// whilst the SDK supports updating, the API doesn't:
		//   Code="PropertyChangeNotAllowed"
		//   Message="Changing property 'windowsConfiguration.additionalUnattendContent' is not allowed."
		//   Target="windowsConfiguration.additionalUnattendContent
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"content": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					ForceNew:  true,
					Sensitive: true,
				},
				"setting": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForSettingNames(), false),
				},
			},
		},
	}
}

func AdditionalUnattendContentSchemaVM() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		// whilst the SDK supports updating, the API doesn't:
		//   Code="PropertyChangeNotAllowed"
		//   Message="Changing property 'windowsConfiguration.additionalUnattendContent' is not allowed."
		//   Target="windowsConfiguration.additionalUnattendContent
		ForceNew: true,
		ConflictsWith: []string{
			"os_managed_disk_id",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"content": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					ForceNew:  true,
					Sensitive: true,
				},
				"setting": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForSettingNames(), false),
				},
			},
		},
	}
}

func ExpandAdditionalUnattendContent(input []any) *[]virtualmachines.AdditionalUnattendContent {
	output := make([]virtualmachines.AdditionalUnattendContent, 0)

	for _, v := range input {
		raw := v.(map[string]any)

		output = append(output, virtualmachines.AdditionalUnattendContent{
			SettingName: pointer.ToEnum[virtualmachines.SettingNames](raw["setting"].(string)),
			Content:     pointer.To(raw["content"].(string)),

			// no other possible values
			PassName:      pointer.To(virtualmachines.PassNamesOobeSystem),
			ComponentName: pointer.To(virtualmachines.ComponentNamesMicrosoftNegativeWindowsNegativeShellNegativeSetup),
		})
	}

	return &output
}

func ExpandAdditionalUnattendContentVMSS(input []any) *[]virtualmachinescalesets.AdditionalUnattendContent {
	output := make([]virtualmachinescalesets.AdditionalUnattendContent, 0)

	for _, v := range input {
		raw := v.(map[string]any)

		output = append(output, virtualmachinescalesets.AdditionalUnattendContent{
			SettingName: pointer.ToEnum[virtualmachinescalesets.SettingNames](raw["setting"].(string)),
			Content:     pointer.To(raw["content"].(string)),

			// no other possible values
			PassName:      pointer.To(virtualmachinescalesets.PassNameOobeSystem),
			ComponentName: pointer.To(virtualmachinescalesets.ComponentNameMicrosoftNegativeWindowsNegativeShellNegativeSetup),
		})
	}

	return &output
}

func FlattenAdditionalUnattendContent(input *[]virtualmachines.AdditionalUnattendContent, d *pluginsdk.ResourceData) []any {
	if input == nil {
		return []any{}
	}

	existing := make([]any, 0)
	if v, ok := d.GetOk("additional_unattend_content"); ok {
		existing = v.([]any)
	}

	output := make([]any, 0)
	for i, v := range *input {
		// content isn't returned from the API as it's sensitive so we need to look it up
		content := ""
		if len(existing) > i {
			existingVal := existing[i]
			existingRaw, ok := existingVal.(map[string]any)
			if ok {
				contentRaw, ok := existingRaw["content"]
				if ok {
					content = contentRaw.(string)
				}
			}
		}

		output = append(output, map[string]any{
			"content": content,
			"setting": pointer.From(v.SettingName),
		})
	}

	return output
}

func FlattenAdditionalUnattendContentVMSS(input *[]virtualmachinescalesets.AdditionalUnattendContent, d *pluginsdk.ResourceData) []any {
	if input == nil {
		return []any{}
	}

	existing := make([]any, 0)
	if v, ok := d.GetOk("additional_unattend_content"); ok {
		existing = v.([]any)
	}

	output := make([]any, 0)
	for i, v := range *input {
		// content isn't returned from the API as it's sensitive so we need to look it up
		content := ""
		if len(existing) > i {
			existingVal := existing[i]
			existingRaw, ok := existingVal.(map[string]any)
			if ok {
				contentRaw, ok := existingRaw["content"]
				if ok {
					content = contentRaw.(string)
				}
			}
		}

		output = append(output, map[string]any{
			"content": content,
			"setting": pointer.From(v.SettingName),
		})
	}

	return output
}

func BootDiagnosticsSchema() *pluginsdk.Schema {
	// lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// TODO: should this be `storage_account_endpoint`?
				"storage_account_uri": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// TODO: validation
				},
			},
		},
	}
}

func ExpandBootDiagnostics(input []any) *virtualmachines.DiagnosticsProfile {
	if len(input) == 0 {
		return &virtualmachines.DiagnosticsProfile{
			BootDiagnostics: &virtualmachines.BootDiagnostics{
				Enabled:    pointer.To(false),
				StorageUri: pointer.To(""),
			},
		}
	}

	// this serves the managed boot diagnostics, in this case we only have this empty block without `storage_account_uri` set
	if input[0] == nil {
		return &virtualmachines.DiagnosticsProfile{
			BootDiagnostics: &virtualmachines.BootDiagnostics{
				Enabled:    pointer.To(true),
				StorageUri: pointer.To(""),
			},
		}
	}

	raw := input[0].(map[string]any)

	storageAccountUri := raw["storage_account_uri"].(string)

	return &virtualmachines.DiagnosticsProfile{
		BootDiagnostics: &virtualmachines.BootDiagnostics{
			Enabled:    pointer.To(true),
			StorageUri: pointer.To(storageAccountUri),
		},
	}
}

func ExpandBootDiagnosticsVMSS(input []any) *virtualmachinescalesets.DiagnosticsProfile {
	if len(input) == 0 {
		return &virtualmachinescalesets.DiagnosticsProfile{
			BootDiagnostics: &virtualmachinescalesets.BootDiagnostics{
				Enabled:    pointer.To(false),
				StorageUri: pointer.To(""),
			},
		}
	}

	// this serves the managed boot diagnostics, in this case we only have this empty block without `storage_account_uri` set
	if input[0] == nil {
		return &virtualmachinescalesets.DiagnosticsProfile{
			BootDiagnostics: &virtualmachinescalesets.BootDiagnostics{
				Enabled:    pointer.To(true),
				StorageUri: pointer.To(""),
			},
		}
	}

	raw := input[0].(map[string]any)

	storageAccountUri := raw["storage_account_uri"].(string)

	return &virtualmachinescalesets.DiagnosticsProfile{
		BootDiagnostics: &virtualmachinescalesets.BootDiagnostics{
			Enabled:    pointer.To(true),
			StorageUri: pointer.To(storageAccountUri),
		},
	}
}

func FlattenBootDiagnostics(input *virtualmachines.DiagnosticsProfile) []any {
	if input == nil || input.BootDiagnostics == nil || input.BootDiagnostics.Enabled == nil || !*input.BootDiagnostics.Enabled {
		return []any{}
	}

	storageAccountUri := pointer.From(input.BootDiagnostics.StorageUri)

	return []any{
		map[string]any{
			"storage_account_uri": storageAccountUri,
		},
	}
}

func FlattenBootDiagnosticsVMSS(input *virtualmachinescalesets.DiagnosticsProfile) []any {
	if input == nil || input.BootDiagnostics == nil || input.BootDiagnostics.Enabled == nil || !*input.BootDiagnostics.Enabled {
		return []any{}
	}

	storageAccountUri := pointer.From(input.BootDiagnostics.StorageUri)

	return []any{
		map[string]any{
			"storage_account_uri": storageAccountUri,
		},
	}
}

func LinuxSecretSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// whilst this isn't present in the nested object it's required when this is specified
				"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

				// whilst we /could/ flatten this to `certificate_urls` we're intentionally not to keep this
				// closer to the Windows VMSS resource, which will also take a `store` param
				"certificate": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"url": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret),
							},
						},
					},
				},
			},
		},
	}
}

func ExpandLinuxSecrets(input []any) *[]virtualmachines.VaultSecretGroup {
	output := make([]virtualmachines.VaultSecretGroup, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		keyVaultId := v["key_vault_id"].(string)
		certificatesRaw := v["certificate"].(*pluginsdk.Set).List()
		certificates := make([]virtualmachines.VaultCertificate, 0)
		for _, certificateRaw := range certificatesRaw {
			certificateV := certificateRaw.(map[string]any)

			url := certificateV["url"].(string)
			certificates = append(certificates, virtualmachines.VaultCertificate{
				CertificateURL: pointer.To(url),
			})
		}

		output = append(output, virtualmachines.VaultSecretGroup{
			SourceVault: &virtualmachines.SubResource{
				Id: pointer.To(keyVaultId),
			},
			VaultCertificates: &certificates,
		})
	}

	return &output
}

func ExpandLinuxSecretsVMSS(input []any) *[]virtualmachinescalesets.VaultSecretGroup {
	output := make([]virtualmachinescalesets.VaultSecretGroup, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		keyVaultId := v["key_vault_id"].(string)
		certificatesRaw := v["certificate"].(*pluginsdk.Set).List()
		certificates := make([]virtualmachinescalesets.VaultCertificate, 0)
		for _, certificateRaw := range certificatesRaw {
			certificateV := certificateRaw.(map[string]any)

			url := certificateV["url"].(string)
			certificates = append(certificates, virtualmachinescalesets.VaultCertificate{
				CertificateURL: pointer.To(url),
			})
		}

		output = append(output, virtualmachinescalesets.VaultSecretGroup{
			SourceVault: &virtualmachinescalesets.SubResource{
				Id: pointer.To(keyVaultId),
			},
			VaultCertificates: &certificates,
		})
	}

	return &output
}

func FlattenLinuxSecrets(input *[]virtualmachines.VaultSecretGroup) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input {
		keyVaultId := ""
		if v.SourceVault != nil && v.SourceVault.Id != nil {
			keyVaultId = *v.SourceVault.Id
		}

		certificates := make([]any, 0)

		if v.VaultCertificates != nil {
			for _, c := range *v.VaultCertificates {
				if c.CertificateURL == nil {
					continue
				}

				certificates = append(certificates, map[string]any{
					"url": *c.CertificateURL,
				})
			}
		}

		output = append(output, map[string]any{
			"key_vault_id": keyVaultId,
			"certificate":  certificates,
		})
	}

	return output
}

func FlattenLinuxSecretsVMSS(input *[]virtualmachinescalesets.VaultSecretGroup) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input {
		keyVaultId := ""
		if v.SourceVault != nil && v.SourceVault.Id != nil {
			keyVaultId = *v.SourceVault.Id
		}

		certificates := make([]any, 0)

		if v.VaultCertificates != nil {
			for _, c := range *v.VaultCertificates {
				if c.CertificateURL == nil {
					continue
				}

				certificates = append(certificates, map[string]any{
					"url": *c.CertificateURL,
				})
			}
		}

		output = append(output, map[string]any{
			"key_vault_id": keyVaultId,
			"certificate":  certificates,
		})
	}

	return output
}

func PlanSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},

				"product": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},

				"publisher": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func ExpandPlan(input []any) *virtualmachines.Plan {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)

	return &virtualmachines.Plan{
		Name:      pointer.To(raw["name"].(string)),
		Product:   pointer.To(raw["product"].(string)),
		Publisher: pointer.To(raw["publisher"].(string)),
	}
}

func ExpandPlanVMSS(input []any) *virtualmachinescalesets.Plan {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)

	return &virtualmachinescalesets.Plan{
		Name:      pointer.To(raw["name"].(string)),
		Product:   pointer.To(raw["product"].(string)),
		Publisher: pointer.To(raw["publisher"].(string)),
	}
}

func FlattenPlan(input *virtualmachines.Plan) []any {
	if input == nil {
		return []any{}
	}

	name := pointer.From(input.Name)

	product := pointer.From(input.Product)

	publisher := pointer.From(input.Publisher)

	return []any{
		map[string]any{
			"name":      name,
			"product":   product,
			"publisher": publisher,
		},
	}
}

func FlattenPlanVMSS(input *virtualmachinescalesets.Plan) []any {
	if input == nil {
		return []any{}
	}

	name := pointer.From(input.Name)

	product := pointer.From(input.Product)

	publisher := pointer.From(input.Publisher)

	return []any{
		map[string]any{
			"name":      name,
			"product":   product,
			"publisher": publisher,
		},
	}
}

func SourceImageReferenceSchemaVM() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		ExactlyOneOf: []string{
			"os_managed_disk_id",
			"source_image_id",
			"source_image_reference",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"offer": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"sku": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SourceImageReferenceSchema(isVirtualMachine bool) *pluginsdk.Schema {
	// whilst originally I was hoping we could use the 'id' from `azurerm_platform_image' unfortunately Azure doesn't
	// like this as a value for the 'id' field:
	// Id /...../Versions/16.04.201909091 is not a valid resource reference."
	// as such the image is split into two fields (source_image_id and source_image_reference) to provide better validation
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: isVirtualMachine,
		MaxItems: 1,
		ExactlyOneOf: []string{
			"source_image_id",
			"source_image_reference",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     isVirtualMachine,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"offer": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     isVirtualMachine,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"sku": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     isVirtualMachine,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     isVirtualMachine,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SourceImageReferenceSchemaOrchestratedVMSS() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		ConflictsWith: []string{
			"source_image_id",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"offer": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"sku": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func IsValidHotPatchSourceImageReference(referenceInput []any, imageId string) bool {
	if imageId != "" {
		return false
	}

	if len(referenceInput) == 0 {
		return false
	}

	raw := referenceInput[0].(map[string]any)
	pub := raw["publisher"].(string)
	offer := raw["offer"].(string)
	sku := raw["sku"].(string)

	supportedSkus := []string{
		"2022-datacenter-azure-edition-core",
		"2022-datacenter-azure-edition-core-smalldisk",
		"2022-datacenter-azure-edition-hotpatch",
		"2022-datacenter-azure-edition-hotpatch-smalldisk",
		"2025-datacenter-azure-edition",
		"2025-datacenter-azure-edition-smalldisk",
		"2025-datacenter-azure-edition-core",
		"2025-datacenter-azure-edition-core-smalldisk",
	}

	if pub == "MicrosoftWindowsServer" && offer == "WindowsServer" && slices.Contains(supportedSkus, sku) {
		return true
	}

	return false
}

func ExpandSourceImageReference(referenceInput []any, imageId string) *virtualmachines.ImageReference {
	if imageId != "" {
		// With Version            : "/communityGalleries/publicGalleryName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/communityGalleries/publicGalleryName/images/myGalleryImageName"
		if _, errors := validation.Any(validate.CommunityGalleryImageID, validate.CommunityGalleryImageVersionID)(imageId, "source_image_id"); len(errors) == 0 {
			return &virtualmachines.ImageReference{
				CommunityGalleryImageId: pointer.To(imageId),
			}
		}

		// With Version            : "/sharedGalleries/galleryUniqueName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/sharedGalleries/galleryUniqueName/images/myGalleryImageName"
		if _, errors := validation.Any(validate.SharedGalleryImageID, validate.SharedGalleryImageVersionID)(imageId, "source_image_id"); len(errors) == 0 {
			return &virtualmachines.ImageReference{
				SharedGalleryImageId: pointer.To(imageId),
			}
		}

		return &virtualmachines.ImageReference{
			Id: pointer.To(imageId),
		}
	}

	raw := referenceInput[0].(map[string]any)
	return &virtualmachines.ImageReference{
		Publisher: pointer.To(raw["publisher"].(string)),
		Offer:     pointer.To(raw["offer"].(string)),
		Sku:       pointer.To(raw["sku"].(string)),
		Version:   pointer.To(raw["version"].(string)),
	}
}

func ExpandSourceImageReferenceVMSS(referenceInput []any, imageId string) *virtualmachinescalesets.ImageReference {
	if imageId != "" {
		// With Version            : "/communityGalleries/publicGalleryName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/communityGalleries/publicGalleryName/images/myGalleryImageName"
		if _, errors := validation.Any(validate.CommunityGalleryImageID, validate.CommunityGalleryImageVersionID)(imageId, "source_image_id"); len(errors) == 0 {
			return &virtualmachinescalesets.ImageReference{
				CommunityGalleryImageId: pointer.To(imageId),
			}
		}

		// With Version            : "/sharedGalleries/galleryUniqueName/images/myGalleryImageName/versions/(major.minor.patch | latest)"
		// Versionless(e.g. latest): "/sharedGalleries/galleryUniqueName/images/myGalleryImageName"
		if _, errors := validation.Any(validate.SharedGalleryImageID, validate.SharedGalleryImageVersionID)(imageId, "source_image_id"); len(errors) == 0 {
			return &virtualmachinescalesets.ImageReference{
				SharedGalleryImageId: pointer.To(imageId),
			}
		}

		return &virtualmachinescalesets.ImageReference{
			Id: pointer.To(imageId),
		}
	}

	raw := referenceInput[0].(map[string]any)
	return &virtualmachinescalesets.ImageReference{
		Publisher: pointer.To(raw["publisher"].(string)),
		Offer:     pointer.To(raw["offer"].(string)),
		Sku:       pointer.To(raw["sku"].(string)),
		Version:   pointer.To(raw["version"].(string)),
	}
}

func FlattenSourceImageReference(input *virtualmachines.ImageReference, hasImageId bool) []any {
	// since the image id is pulled out as a separate field, if that's set we should return an empty block here
	if input == nil || hasImageId {
		return []any{}
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

	return []any{
		map[string]any{
			"publisher": publisher,
			"offer":     offer,
			"sku":       sku,
			"version":   version,
		},
	}
}

func FlattenSourceImageReferenceVMSS(input *virtualmachinescalesets.ImageReference, hasImageId bool) []any {
	// since the image id is pulled out as a separate field, if that's set we should return an empty block here
	if input == nil || hasImageId {
		return []any{}
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

	return []any{
		map[string]any{
			"publisher": publisher,
			"offer":     offer,
			"sku":       sku,
			"version":   version,
		},
	}
}

func WinRmListenerSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		// Whilst the SDK allows you to modify this, the API does not:
		//   Code="PropertyChangeNotAllowed"
		//   Message="Changing property 'windowsConfiguration.winRM.listeners' is not allowed."
		//   Target="windowsConfiguration.winRM.listeners"
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"protocol": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForProtocolTypes(), false),
				},

				"certificate_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret),
				},
			},
		},
	}
}

func WinRmListenerSchemaVM() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		// Whilst the SDK allows you to modify this, the API does not:
		//   Code="PropertyChangeNotAllowed"
		//   Message="Changing property 'windowsConfiguration.winRM.listeners' is not allowed."
		//   Target="windowsConfiguration.winRM.listeners"
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"protocol": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(virtualmachines.PossibleValuesForProtocolTypes(), false),
				},

				"certificate_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret),
				},
			},
		},
		ConflictsWith: []string{
			"os_managed_disk_id",
		},
	}
}

func ExpandWinRMListener(input []any) *virtualmachines.WinRMConfiguration {
	listeners := make([]virtualmachines.WinRMListener, 0)

	for _, v := range input {
		raw := v.(map[string]any)

		listener := virtualmachines.WinRMListener{
			Protocol: pointer.ToEnum[virtualmachines.ProtocolTypes](raw["protocol"].(string)),
		}

		certificateUrl := raw["certificate_url"].(string)
		if certificateUrl != "" {
			listener.CertificateURL = pointer.To(certificateUrl)
		}

		listeners = append(listeners, listener)
	}

	return &virtualmachines.WinRMConfiguration{
		Listeners: &listeners,
	}
}

func ExpandWinRMListenerVMSS(input []any) *virtualmachinescalesets.WinRMConfiguration {
	listeners := make([]virtualmachinescalesets.WinRMListener, 0)

	for _, v := range input {
		raw := v.(map[string]any)

		listener := virtualmachinescalesets.WinRMListener{
			Protocol: pointer.ToEnum[virtualmachinescalesets.ProtocolTypes](raw["protocol"].(string)),
		}

		certificateUrl := raw["certificate_url"].(string)
		if certificateUrl != "" {
			listener.CertificateURL = pointer.To(certificateUrl)
		}

		listeners = append(listeners, listener)
	}

	return &virtualmachinescalesets.WinRMConfiguration{
		Listeners: &listeners,
	}
}

func FlattenWinRMListener(input *virtualmachines.WinRMConfiguration) []any {
	if input == nil || input.Listeners == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input.Listeners {
		certificateUrl := pointer.From(v.CertificateURL)

		output = append(output, map[string]any{
			"certificate_url": certificateUrl,
			"protocol":        pointer.From(v.Protocol),
		})
	}

	return output
}

func FlattenWinRMListenerVMSS(input *virtualmachinescalesets.WinRMConfiguration) []any {
	if input == nil || input.Listeners == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input.Listeners {
		certificateUrl := pointer.From(v.CertificateURL)

		output = append(output, map[string]any{
			"certificate_url": certificateUrl,
			"protocol":        pointer.From(v.Protocol),
		})
	}

	return output
}

func WindowsSecretSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// whilst this isn't present in the nested object it's required when this is specified
				"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

				"certificate": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"store": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
							"url": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret),
							},
						},
					},
				},
			},
		},
	}
}

func WindowsSecretSchemaVM() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// whilst this isn't present in the nested object it's required when this is specified
				"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

				"certificate": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"store": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
							"url": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret),
							},
						},
					},
				},
			},
		},
		ConflictsWith: []string{
			"os_managed_disk_id",
		},
	}
}

func ExpandWindowsSecrets(input []any) *[]virtualmachines.VaultSecretGroup {
	output := make([]virtualmachines.VaultSecretGroup, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		keyVaultId := v["key_vault_id"].(string)
		certificatesRaw := v["certificate"].(*pluginsdk.Set).List()
		certificates := make([]virtualmachines.VaultCertificate, 0)
		for _, certificateRaw := range certificatesRaw {
			certificateV := certificateRaw.(map[string]any)

			store := certificateV["store"].(string)
			url := certificateV["url"].(string)
			certificates = append(certificates, virtualmachines.VaultCertificate{
				CertificateStore: pointer.To(store),
				CertificateURL:   pointer.To(url),
			})
		}

		output = append(output, virtualmachines.VaultSecretGroup{
			SourceVault: &virtualmachines.SubResource{
				Id: pointer.To(keyVaultId),
			},
			VaultCertificates: &certificates,
		})
	}

	return &output
}

func ExpandWindowsSecretsVMSS(input []any) *[]virtualmachinescalesets.VaultSecretGroup {
	output := make([]virtualmachinescalesets.VaultSecretGroup, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		keyVaultId := v["key_vault_id"].(string)
		certificatesRaw := v["certificate"].(*pluginsdk.Set).List()
		certificates := make([]virtualmachinescalesets.VaultCertificate, 0)
		for _, certificateRaw := range certificatesRaw {
			certificateV := certificateRaw.(map[string]any)

			store := certificateV["store"].(string)
			url := certificateV["url"].(string)
			certificates = append(certificates, virtualmachinescalesets.VaultCertificate{
				CertificateStore: pointer.To(store),
				CertificateURL:   pointer.To(url),
			})
		}

		output = append(output, virtualmachinescalesets.VaultSecretGroup{
			SourceVault: &virtualmachinescalesets.SubResource{
				Id: pointer.To(keyVaultId),
			},
			VaultCertificates: &certificates,
		})
	}

	return &output
}

func FlattenWindowsSecrets(input *[]virtualmachines.VaultSecretGroup) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input {
		keyVaultId := ""
		if v.SourceVault != nil && v.SourceVault.Id != nil {
			keyVaultId = *v.SourceVault.Id
		}

		certificates := make([]any, 0)

		if v.VaultCertificates != nil {
			for _, c := range *v.VaultCertificates {
				store := pointer.From(c.CertificateStore)

				url := pointer.From(c.CertificateURL)

				certificates = append(certificates, map[string]any{
					"store": store,
					"url":   url,
				})
			}
		}

		output = append(output, map[string]any{
			"key_vault_id": keyVaultId,
			"certificate":  certificates,
		})
	}

	return output
}

func FlattenWindowsSecretsVMSS(input *[]virtualmachinescalesets.VaultSecretGroup) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, v := range *input {
		keyVaultId := ""
		if v.SourceVault != nil && v.SourceVault.Id != nil {
			keyVaultId = *v.SourceVault.Id
		}

		certificates := make([]any, 0)

		if v.VaultCertificates != nil {
			for _, c := range *v.VaultCertificates {
				store := pointer.From(c.CertificateStore)

				url := pointer.From(c.CertificateURL)

				certificates = append(certificates, map[string]any{
					"store": store,
					"url":   url,
				})
			}
		}

		output = append(output, map[string]any{
			"key_vault_id": keyVaultId,
			"certificate":  certificates,
		})
	}

	return output
}
