package eventhubsclusters

import "strings"

type ClusterSkuName string

const (
	ClusterSkuNameDedicated ClusterSkuName = "Dedicated"
)

func PossibleValuesForClusterSkuName() []string {
	return []string{
		"Dedicated",
	}
}

func parseClusterSkuName(input string) (*ClusterSkuName, error) {
	vals := map[string]ClusterSkuName{
		"dedicated": "Dedicated",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := ClusterSkuName(v)
	return &out, nil
}
