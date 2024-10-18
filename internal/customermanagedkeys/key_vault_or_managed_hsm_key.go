package customermanagedkeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	hsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VersionType int

const (
	VersionTypeAny VersionType = iota
	VersionTypeVersioned
	VersionTypeVersionless
)

type KeyVaultOrManagedHSMKey struct {
	KeyVaultKeyId              *parse.NestedItemId
	ManagedHSMKeyId            *hsmParse.ManagedHSMDataPlaneVersionedKeyId
	ManagedHSMKeyVersionlessId *hsmParse.ManagedHSMDataPlaneVersionlessKeyId
}

func (k *KeyVaultOrManagedHSMKey) IsSet() bool {
	return k != nil && (k.KeyVaultKeyId != nil || k.ManagedHSMKeyId != nil || k.ManagedHSMKeyVersionlessId != nil)
}

func (k *KeyVaultOrManagedHSMKey) ID() string {
	if k == nil {
		return ""
	}

	if k.KeyVaultKeyId != nil {
		return k.KeyVaultKeyId.ID()
	}

	if k.ManagedHSMKeyId != nil {
		return k.ManagedHSMKeyId.ID()
	}

	if k.ManagedHSMKeyVersionlessId != nil {
		return k.ManagedHSMKeyVersionlessId.ID()
	}

	return ""
}

func (k *KeyVaultOrManagedHSMKey) KeyVaultKeyID() string {
	if k != nil && k.KeyVaultKeyId != nil {
		return k.KeyVaultKeyId.ID()
	}
	return ""
}

func (k *KeyVaultOrManagedHSMKey) ManagedHSMKeyID() string {
	if k != nil && k.ManagedHSMKeyId != nil {
		return k.ManagedHSMKeyId.ID()
	}

	if k != nil && k.ManagedHSMKeyVersionlessId != nil {
		return k.ManagedHSMKeyVersionlessId.ID()
	}

	return ""
}

func (k *KeyVaultOrManagedHSMKey) BaseUri() string {
	if k.KeyVaultKeyId != nil {
		return k.KeyVaultKeyId.KeyVaultBaseUrl
	}

	if k.ManagedHSMKeyId != nil {
		return k.ManagedHSMKeyId.BaseUri()
	}

	if k.ManagedHSMKeyVersionlessId != nil {
		return k.ManagedHSMKeyVersionlessId.BaseUri()
	}

	return ""
}

func parseKeyvaultID(keyRaw string, requireVersion VersionType, _ environments.Api) (*parse.NestedItemId, error) {
	keyID, err := parse.ParseOptionallyVersionedNestedKeyID(keyRaw)
	if err != nil {
		return nil, err
	}

	if requireVersion == VersionTypeVersioned && keyID.Version == "" {
		return nil, fmt.Errorf("expected a key vault versioned ID but no version information was found in: %q", keyRaw)
	}

	if requireVersion == VersionTypeVersionless && keyID.Version != "" {
		return nil, fmt.Errorf("expected a key vault versionless ID but version information was found in: %q", keyRaw)
	}

	return keyID, nil
}

func parseManagedHSMKey(keyRaw string, requireVersion VersionType, hsmEnv environments.Api) (
	versioned *hsmParse.ManagedHSMDataPlaneVersionedKeyId, versionless *hsmParse.ManagedHSMDataPlaneVersionlessKeyId, err error) {
	// if specified with hasVersion == True, then it has to be parsed as versionedKeyID
	var domainSuffix *string
	if hsmEnv != nil {
		domainSuffix, _ = hsmEnv.DomainSuffix()
	}

	switch requireVersion {
	case VersionTypeAny:
		if versioned, err = hsmParse.ManagedHSMDataPlaneVersionedKeyID(keyRaw, domainSuffix); err != nil {
			if versionless, err = hsmParse.ManagedHSMDataPlaneVersionlessKeyID(keyRaw, domainSuffix); err != nil {
				return nil, nil, fmt.Errorf("parse Managed HSM both versionedID and versionlessID err for %s", keyRaw)
			}
		}
	case VersionTypeVersioned:
		versioned, err = hsmParse.ManagedHSMDataPlaneVersionedKeyID(keyRaw, domainSuffix)
	case VersionTypeVersionless:
		versionless, err = hsmParse.ManagedHSMDataPlaneVersionlessKeyID(keyRaw, domainSuffix)
	}

	return versioned, versionless, err
}

func ExpandKeyVaultOrManagedHSMKey(d interface{}, requireVersion VersionType, keyVaultEnv, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
	return ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey(d, requireVersion, "key_vault_key_id", "managed_hsm_key_id", keyVaultEnv, hsmEnv)
}

// ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey
// d: should be one of *pluginsdk.ResourceData or map[string]interface{}
// if return nil, nil, it means no key_vault_key_id or managed_hsm_key_id is specified
func ExpandKeyVaultOrManagedHSMKeyWithCustomFieldKey(d interface{}, requireVersion VersionType, keyVaultFieldName, hsmFieldName string, keyVaultEnv, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
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
			vaultKeyStr, _ = keyRaw.(string)
		}
		if keyRaw, ok := obj[hsmFieldName]; ok {
			hsmKeyStr, _ = keyRaw.(string)
		}
	} else {
		return nil, fmt.Errorf("not supported data type to parse CMK: %T", d)
	}

	switch {
	case vaultKeyStr != "":
		if key.KeyVaultKeyId, err = parseKeyvaultID(vaultKeyStr, requireVersion, keyVaultEnv); err != nil {
			return nil, err
		}
	case hsmKeyStr != "":
		if key.ManagedHSMKeyId, key.ManagedHSMKeyVersionlessId, err = parseManagedHSMKey(hsmKeyStr, requireVersion, hsmEnv); err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}
	return key, err
}

// FlattenKeyVaultOrManagedHSMID uses `KeyVaultOrManagedHSMKey.SetState()` to save the state, which this function is designed not to do.
func FlattenKeyVaultOrManagedHSMID(id string, keyVaultEnv, hsmEnv environments.Api) (*KeyVaultOrManagedHSMKey, error) {
	_ = keyVaultEnv
	if id == "" {
		return nil, nil
	}

	key := &KeyVaultOrManagedHSMKey{}
	var err error
	key.KeyVaultKeyId, err = parse.ParseOptionallyVersionedNestedKeyID(id)
	if err == nil {
		return key, nil
	}

	var domainSuffix *string
	if hsmEnv != nil {
		domainSuffix, _ = hsmEnv.DomainSuffix()
	}
	if key.ManagedHSMKeyId, err = hsmParse.ManagedHSMDataPlaneVersionedKeyID(id, domainSuffix); err == nil {
		return key, nil
	}

	if key.ManagedHSMKeyVersionlessId, err = hsmParse.ManagedHSMDataPlaneVersionlessKeyID(id, domainSuffix); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("cannot parse given id to key vault key nor managed hsm key: %s", id)
}
