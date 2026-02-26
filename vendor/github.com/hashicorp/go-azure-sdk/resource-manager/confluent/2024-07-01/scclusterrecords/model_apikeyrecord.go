package scclusterrecords

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type APIKeyRecord struct {
	Id         *string           `json:"id,omitempty"`
	Kind       *string           `json:"kind,omitempty"`
	Properties *APIKeyProperties `json:"properties,omitempty"`
}
