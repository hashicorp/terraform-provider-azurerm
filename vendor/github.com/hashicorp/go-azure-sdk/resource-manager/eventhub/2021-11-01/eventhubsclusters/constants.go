package eventhubsclusters

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterSkuName string

const (
	ClusterSkuNameDedicated ClusterSkuName = "Dedicated"
)

func PossibleValuesForClusterSkuName() []string {
	return []string{
		string(ClusterSkuNameDedicated),
	}
}

func parseClusterSkuName(input string) (*ClusterSkuName, error) {
	vals := map[string]ClusterSkuName{
		"dedicated": ClusterSkuNameDedicated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterSkuName(input)
	return &out, nil
}
