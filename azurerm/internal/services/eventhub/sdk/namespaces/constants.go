package namespaces

type IdentityType string

const (
	IdentityTypeSystemAssigned IdentityType = "SystemAssigned"
)

type KeySource string

const (
	KeySourceMicrosoftKeyVault KeySource = "Microsoft.KeyVault"
)

type SkuName string

const (
	SkuNameBasic    SkuName = "Basic"
	SkuNameStandard SkuName = "Standard"
)

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierStandard SkuTier = "Standard"
)
