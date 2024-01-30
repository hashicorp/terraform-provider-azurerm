package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BuildpacksGroupProperties struct {
	Buildpacks *[]BuildpackProperties `json:"buildpacks,omitempty"`
	Name       *string                `json:"name,omitempty"`
}
