package checknameavailability

import "strings"

type ResourceType string

const (
	ResourceTypeMicrosoftPointCdnProfilesEndpoints ResourceType = "Microsoft.Cdn/Profiles/Endpoints"
)

func PossibleValuesForResourceType() []string {
	return []string{
		string(ResourceTypeMicrosoftPointCdnProfilesEndpoints),
	}
}

func parseResourceType(input string) (*ResourceType, error) {
	vals := map[string]ResourceType{
		"microsoft.cdn/profiles/endpoints": ResourceTypeMicrosoftPointCdnProfilesEndpoints,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceType(input)
	return &out, nil
}
