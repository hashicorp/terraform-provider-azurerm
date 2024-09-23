package customermanagedkeys_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	cmk "github.com/hashicorp/terraform-provider-azurerm/internal/customermanagedkeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	hsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

func buildData(keyVaultKey, keyVualtValue, hsmKey, hsmValue string) interface{} {
	data := map[string]interface{}{}
	if keyVaultKey != "" {
		data[keyVaultKey] = keyVualtValue
	}

	if hsmKey != "" {
		data[hsmKey] = hsmValue
	}

	return data
}

func buildKeyVaultData(key, value string) interface{} {
	return buildData(key, value, "", "")
}

func buildHSMData(key, value string) interface{} {
	return buildData("", "", key, value)
}

func TestExpandKeyVaultOrManagedHSMKeyKey(t *testing.T) {
	type args struct {
		d                 interface{}
		hasVersion        *bool
		keyVaultFieldName string
		hsmFieldName      string
		hsmEnv            environments.Api
	}
	tests := []struct {
		name    string
		args    args
		want    *cmk.KeyVaultOrManagedHSMKey
		wantErr bool
	}{
		{
			args: args{
				d:                 buildKeyVaultData("key_vault_key_id", "https://test.keyvault.azure.net/keys/test-key-name"),
				keyVaultFieldName: "key_vault_key_id",
			},
			want: &cmk.KeyVaultOrManagedHSMKey{
				KeyVaultKeyID: &parse.NestedItemId{
					KeyVaultBaseUrl: "https://test.keyvault.azure.net/",
					NestedItemType:  "keys",
					Name:            "test-key-name",
				},
			},
		},
		{
			args: args{
				d:                 buildKeyVaultData("key_vault_key_id", "https://test.keyvault.azure.net/certs/test-key-name"),
				keyVaultFieldName: "key_vault_key_id",
			},
			wantErr: true,
		},
		{
			args: args{
				d:                 buildKeyVaultData("key_vault_key_url", "https://test.keyvault.azure.net/keys/test-key-name"),
				keyVaultFieldName: "key_vault_key_id",
			},
			wantErr: true,
		},
		{
			args: args{
				d:            buildHSMData("managed_hsm_key_id", "https://test.managedhsm.azure.net/keys/test-key-name"),
				hsmFieldName: "managed_hsm_key_id",
				hasVersion:   pointer.To(false),
			},
			want: &cmk.KeyVaultOrManagedHSMKey{
				ManagedHSMKeyVersionlessID: &hsmParse.ManagedHSMDataPlaneVersionlessKeyId{
					ManagedHSMName: "test",
					DomainSuffix:   "managedhsm.azure.net",
					KeyName:        "test-key-name",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			got, err := cmk.ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey(tt.args.d, tt.args.hasVersion, tt.args.keyVaultFieldName, tt.args.hsmFieldName, tt.args.hsmEnv)
			if (err != nil) != tt.wantErr {
				t2.Errorf("ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t2.Errorf("ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
