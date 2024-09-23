package customermanagedkeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	hsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KeyVaultOrManagedHSMKey struct {
	KeyVaultKeyID              *parse.NestedItemId
	ManagedHSMKeyID            *hsmParse.ManagedHSMDataPlaneVersionedKeyId
	ManagedHSMKeyVersionlessID *hsmParse.ManagedHSMDataPlaneVersionlessKeyId
}

func (k *KeyVaultOrManagedHSMKey) ID() string {
	if k == nil {
		return ""
	}

	if k.KeyVaultKeyID != nil {
		return k.KeyVaultKeyID.ID()
	}

	if k.ManagedHSMKeyID != nil {
		return k.ManagedHSMKeyID.ID()
	}

	if k.ManagedHSMKeyVersionlessID != nil {
		return k.ManagedHSMKeyVersionlessID.ID()
	}

	return ""
}

func (k *KeyVaultOrManagedHSMKey) BaseUri() string {
	if k.KeyVaultKeyID != nil {
		return k.KeyVaultKeyID.KeyVaultBaseUrl
	}

	if k.ManagedHSMKeyID != nil {
		return k.ManagedHSMKeyID.BaseUri()
	}

	if k.ManagedHSMKeyVersionlessID != nil {
		return k.ManagedHSMKeyVersionlessID.BaseUri()
	}

	return ""
}

func expandKeyvauleID(keyRaw string, hasVersion *bool) (*parse.NestedItemId, error) {
	if pointer.From(hasVersion) {
		if keyID, err := parse.ParseNestedKeyID(keyRaw); err == nil {
			return keyID, nil
		} else {
			return nil, err
		}
	}

	if keyID, err := parse.ParseOptionallyVersionedNestedKeyID(keyRaw); err == nil {
		return keyID, nil
	} else {
		return nil, err
	}
}

func expandManagedHSMKey(keyRaw string, hasVersion *bool, hsmEnv environments.Api) (*hsmParse.ManagedHSMDataPlaneVersionedKeyId, *hsmParse.ManagedHSMDataPlaneVersionlessKeyId, error) {
	// if specified with hasVersion == True, then it has to be parsed as versionedKeyID
	var domainSuffix *string
	if hsmEnv != nil {
		domainSuffix, _ = hsmEnv.DomainSuffix()
	}
	if hasVersion == nil || *hasVersion {
		versioned, err := hsmParse.ManagedHSMDataPlaneVersionedKeyID(keyRaw, domainSuffix)
		if err == nil {
			return versioned, nil, nil
		}
		// if required versioned but got error
		if pointer.From(hasVersion) {
			return nil, nil, err
		}
	}

	// versionless or optional version
	if versionless, err := hsmParse.ManagedHSMDataPlaneVersionlessKeyID(keyRaw, domainSuffix); err == nil {
		return nil, versionless, nil
	} else {
		return nil, nil, err
	}
}

// hasVersion:
//   - nil: both versioned or versionless are ok
//   - true: must have version
//   - false: must not have vesrion
func ExpandKeyVaultOrManagedHSMOptionallyVersionedKey(d interface{}, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
	return ExpandKeyVaultOrManagedHSMKey(d, nil, hsmEnv)
}

func ExpandKeyVaultOrManagedHSMKey(d interface{}, hasVersion *bool, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
	return ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey(d, hasVersion, "key_vault_key_id", "managed_hsm_key_id", hsmEnv)
}

func ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey(d interface{}, hasVersion *bool, keyVaultFieldName, hsmFieldName string, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
	key := &KeyVaultOrManagedHSMKey{}
	var err error
	var vaultKeyStr, hsmKeyStr string
	if rd, ok := d.(*pluginsdk.ResourceData); ok {
		if keyRaw, ok := rd.GetOk(keyVaultFieldName); ok {
			vaultKeyStr = keyRaw.(string)
		} else if keyRaw, ok = rd.GetOk(hsmFieldName); ok {
			hsmKeyStr = keyRaw.(string)
		}
	} else if obj, ok := d.(map[string]interface{}); ok {
		if keyRaw, ok := obj[keyVaultFieldName]; ok {
			vaultKeyStr = keyRaw.(string)
		} else if keyRaw, ok = obj[hsmFieldName]; ok {
			hsmKeyStr = keyRaw.(string)
		}
	} else {
		return nil, fmt.Errorf("not supported data type to parse CMK: %T", d)
	}

	if vaultKeyStr != "" {
		if key.KeyVaultKeyID, err = expandKeyvauleID(vaultKeyStr, hasVersion); err != nil {
			return nil, err
		}
	} else if hsmKeyStr != "" {
		if key.ManagedHSMKeyID, key.ManagedHSMKeyVersionlessID, err = expandManagedHSMKey(hsmKeyStr, hasVersion, hsmEnv); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("at least one of `%s` or `%s` should be specified", keyVaultFieldName, hsmFieldName)
	}
	return key, err
}

func FlattenKeyVaultOrManagedHSMID(id string, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
	if id == "" {
		return nil, nil
	}

	key := &KeyVaultOrManagedHSMKey{}
	var err error
	key.KeyVaultKeyID, err = parse.ParseOptionallyVersionedNestedItemID(id)
	if err == nil {
		return key, nil
	}

	var domainSuffix *string
	if hsmEnv != nil {
		domainSuffix, _ = hsmEnv.DomainSuffix()
	}
	if key.ManagedHSMKeyID, err = hsmParse.ManagedHSMDataPlaneVersionedKeyID(id, domainSuffix); err == nil {
		return key, nil
	}

	if key.ManagedHSMKeyVersionlessID, err = hsmParse.ManagedHSMDataPlaneVersionlessKeyID(id, domainSuffix); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("cannot parse given id to key vault key nor managed hsm key: %s", id)
}
