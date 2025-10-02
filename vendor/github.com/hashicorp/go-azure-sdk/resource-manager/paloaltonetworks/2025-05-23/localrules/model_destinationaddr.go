package localrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DestinationAddr struct {
	Cidrs       *[]string `json:"cidrs,omitempty"`
	Countries   *[]string `json:"countries,omitempty"`
	Feeds       *[]string `json:"feeds,omitempty"`
	FqdnLists   *[]string `json:"fqdnLists,omitempty"`
	PrefixLists *[]string `json:"prefixLists,omitempty"`
}
