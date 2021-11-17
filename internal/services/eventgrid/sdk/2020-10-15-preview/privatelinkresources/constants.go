package privatelinkresources

import "strings"

type ParentType string

const (
	ParentTypeDomains ParentType = "domains"
	ParentTypeTopics  ParentType = "topics"
)

func PossibleValuesForParentType() []string {
	return []string{
		string(ParentTypeDomains),
		string(ParentTypeTopics),
	}
}

func parseParentType(input string) (*ParentType, error) {
	vals := map[string]ParentType{
		"domains": ParentTypeDomains,
		"topics":  ParentTypeTopics,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ParentType(input)
	return &out, nil
}
