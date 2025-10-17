package scclusterrecords

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type APIKeyResourceEntity struct {
	Environment  *string `json:"environment,omitempty"`
	Id           *string `json:"id,omitempty"`
	Kind         *string `json:"kind,omitempty"`
	Related      *string `json:"related,omitempty"`
	ResourceName *string `json:"resourceName,omitempty"`
}
