package heatmaps

import "strings"

type HeatMapType string

const (
	HeatMapTypeDefault HeatMapType = "default"
)

func PossibleValuesForHeatMapType() []string {
	return []string{
		string(HeatMapTypeDefault),
	}
}

func parseHeatMapType(input string) (*HeatMapType, error) {
	vals := map[string]HeatMapType{
		"default": HeatMapTypeDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HeatMapType(input)
	return &out, nil
}
