package parse

// interface to unify key of keyvalt and key of hsm behavior

type VaultKeyID interface {
	ID() string
	String() string
}

func NewVaultKeyID(vaulter Vaulter, keyName, version string) VaultKeyID {
	switch vaulter.Type() {
	case VaultTypeMHSM:
		return NewManagedHSMKeyID(vaulter.GetSubscriptionID(), vaulter.GetResourceGroup(), vaulter.GetName(), keyName, version)
	default:
		return NewKeyID(vaulter.GetSubscriptionID(), vaulter.GetResourceGroup(), vaulter.GetName(), keyName, version)
	}
}
func NewVaultKeyVersionlessID(vaulter Vaulter, keyName string) VaultKeyID {
	switch vaulter.Type() {
	case VaultTypeMHSM:
		return NewManagedHSMKeyID(vaulter.GetSubscriptionID(), vaulter.GetResourceGroup(), vaulter.GetName(), keyName, "")
	default:
		return NewKeyVersionlessID(vaulter.GetSubscriptionID(), vaulter.GetResourceGroup(), vaulter.GetName(), keyName)
	}
}
