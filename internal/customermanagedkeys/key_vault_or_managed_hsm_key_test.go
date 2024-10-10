package customermanagedkeys_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/customermanagedkeys"
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
		hasVersion        customermanagedkeys.VersionType
		keyVaultFieldName string
		hsmFieldName      string
		hsmEnv            environments.Api
	}
	tests := []struct {
		name    string
		args    args
		want    *customermanagedkeys.KeyVaultOrManagedHSMKey
		wantErr bool
	}{
		{
			name: "success with key_vault_key_id",
			args: args{
				d:                 buildKeyVaultData("key_vault_key_id", "https://test.keyvault.azure.net/keys/test-key-name"),
				keyVaultFieldName: "key_vault_key_id",
			},
			want: &customermanagedkeys.KeyVaultOrManagedHSMKey{
				KeyVaultKeyId: &parse.NestedItemId{
					KeyVaultBaseUrl: "https://test.keyvault.azure.net/",
					NestedItemType:  "keys",
					Name:            "test-key-name",
				},
			},
		},
		{
			name: "fail with wrong item type: cert",
			args: args{
				d:                 buildKeyVaultData("key_vault_key_id", "https://test.keyvault.azure.net/certs/test-key-name"),
				keyVaultFieldName: "key_vault_key_id",
			},
			wantErr: true,
		},
		{
			name: "fail with wrong field name",
			args: args{
				d:                 buildKeyVaultData("key_vault_key_url", "https://test.keyvault.azure.net/keys/test-key-name"),
				keyVaultFieldName: "key_vault_key_id",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "fail with no version provided",
			args: args{
				d:                 buildKeyVaultData("key_vault_key_id", "https://test.keyvault.azure.net/keys/test-key-name3"),
				keyVaultFieldName: "key_vault_key_id",
				hasVersion:        customermanagedkeys.VersionTypeVersioned,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success with managed_hsm_key_id",
			args: args{
				d:            buildHSMData("managed_hsm_key_id", "https://test.managedhsm.azure.net/keys/test-key-name"),
				hsmFieldName: "managed_hsm_key_id",
				hasVersion:   customermanagedkeys.VersionTypeVersionless,
			},
			want: &customermanagedkeys.KeyVaultOrManagedHSMKey{
				ManagedHSMKeyVersionlessId: &hsmParse.ManagedHSMDataPlaneVersionlessKeyId{
					ManagedHSMName: "test",
					DomainSuffix:   "managedhsm.azure.net",
					KeyName:        "test-key-name",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			got, err := customermanagedkeys.ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey(tt.args.d, tt.args.hasVersion, tt.args.keyVaultFieldName, tt.args.hsmFieldName, nil, tt.args.hsmEnv)
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
