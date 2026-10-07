package ipampools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceBasics struct {
	AddressPrefixes *[]string `json:"addressPrefixes,omitempty"`
	ResourceId      *string   `json:"resourceId,omitempty"`
}
