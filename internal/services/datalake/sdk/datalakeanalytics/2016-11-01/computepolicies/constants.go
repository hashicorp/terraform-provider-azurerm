package computepolicies

import "strings"

type AADObjectType string

const (
	AADObjectTypeGroup            AADObjectType = "Group"
	AADObjectTypeServicePrincipal AADObjectType = "ServicePrincipal"
	AADObjectTypeUser             AADObjectType = "User"
)

func PossibleValuesForAADObjectType() []string {
	return []string{
		string(AADObjectTypeGroup),
		string(AADObjectTypeServicePrincipal),
		string(AADObjectTypeUser),
	}
}

func parseAADObjectType(input string) (*AADObjectType, error) {
	vals := map[string]AADObjectType{
		"group":            AADObjectTypeGroup,
		"serviceprincipal": AADObjectTypeServicePrincipal,
		"user":             AADObjectTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AADObjectType(input)
	return &out, nil
}
