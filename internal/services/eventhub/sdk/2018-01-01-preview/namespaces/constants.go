package namespaces

import "strings"

type KeySource string

const (
	KeySourceMicrosoftPointKeyVault KeySource = "Microsoft.KeyVault"
)

func PossibleValuesForKeySource() []string {
	return []string{
		"Microsoft.KeyVault",
	}
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"microsoftpointkeyvault": "Microsoft.KeyVault",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := KeySource(v)
	return &out, nil
}

type SkuName string

const (
	SkuNameBasic    SkuName = "Basic"
	SkuNameStandard SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		"Basic",
		"Standard",
	}
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"basic":    "Basic",
		"standard": "Standard",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := SkuName(v)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		"Basic",
		"Standard",
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    "Basic",
		"standard": "Standard",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := SkuTier(v)
	return &out, nil
}
