package operation

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerRegistryResourceType string

const (
	ContainerRegistryResourceTypeMicrosoftPointContainerRegistryRegistries ContainerRegistryResourceType = "Microsoft.ContainerRegistry/registries"
)

func PossibleValuesForContainerRegistryResourceType() []string {
	return []string{
		string(ContainerRegistryResourceTypeMicrosoftPointContainerRegistryRegistries),
	}
}

func parseContainerRegistryResourceType(input string) (*ContainerRegistryResourceType, error) {
	vals := map[string]ContainerRegistryResourceType{
		"microsoft.containerregistry/registries": ContainerRegistryResourceTypeMicrosoftPointContainerRegistryRegistries,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerRegistryResourceType(input)
	return &out, nil
}
