package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomService struct {
	Docker               *Docker                         `json:"docker,omitempty"`
	Endpoints            *[]Endpoint                     `json:"endpoints,omitempty"`
	EnvironmentVariables *map[string]EnvironmentVariable `json:"environmentVariables,omitempty"`
	Image                *Image                          `json:"image,omitempty"`
	Name                 *string                         `json:"name,omitempty"`
	Volumes              *[]VolumeDefinition             `json:"volumes,omitempty"`
}
