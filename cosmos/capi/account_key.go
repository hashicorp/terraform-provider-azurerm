package capi

type AccountKey struct {
	Type    string
	Version string
	Key     string
}

func NewAccountKeyWithDefaults(key string) AccountKey {
	return NewAccountKey(key, DefaultCosmosAccountKeyType, DefaultCosmosAccountKeyVersion)
}

func NewAccountKey(key, keyType, keyVersion string) AccountKey {
	return AccountKey{
		Type:    keyType,
		Version: keyVersion,
		Key:     key,
	}
}
