package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func keyPermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.KeyPermissionsBackup),
				string(keyvault.KeyPermissionsCreate),
				string(keyvault.KeyPermissionsDecrypt),
				string(keyvault.KeyPermissionsDelete),
				string(keyvault.KeyPermissionsEncrypt),
				string(keyvault.KeyPermissionsGet),
				string(keyvault.KeyPermissionsImport),
				string(keyvault.KeyPermissionsList),
				string(keyvault.KeyPermissionsPurge),
				string(keyvault.KeyPermissionsRecover),
				string(keyvault.KeyPermissionsRestore),
				string(keyvault.KeyPermissionsSign),
				string(keyvault.KeyPermissionsUnwrapKey),
				string(keyvault.KeyPermissionsUpdate),
				string(keyvault.KeyPermissionsVerify),
				string(keyvault.KeyPermissionsWrapKey),
			}, true),
			DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
		},
	}
}

func secretPermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.SecretPermissionsBackup),
				string(keyvault.SecretPermissionsDelete),
				string(keyvault.SecretPermissionsGet),
				string(keyvault.SecretPermissionsList),
				string(keyvault.SecretPermissionsPurge),
				string(keyvault.SecretPermissionsRecover),
				string(keyvault.SecretPermissionsRestore),
				string(keyvault.SecretPermissionsSet),
			}, true),
			DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
		},
	}
}

func certificatePermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.Create),
				string(keyvault.Delete),
				string(keyvault.Deleteissuers),
				string(keyvault.Get),
				string(keyvault.Getissuers),
				string(keyvault.Import),
				string(keyvault.List),
				string(keyvault.Listissuers),
				string(keyvault.Managecontacts),
				string(keyvault.Manageissuers),
				string(keyvault.Purge),
				string(keyvault.Recover),
				string(keyvault.Setissuers),
				string(keyvault.Update),
			}, true),
			DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
		},
	}
}
