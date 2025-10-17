package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegionSpecEntity struct {
	Cloud      *string   `json:"cloud,omitempty"`
	Name       *string   `json:"name,omitempty"`
	Packages   *[]string `json:"packages,omitempty"`
	RegionName *string   `json:"regionName,omitempty"`
}
