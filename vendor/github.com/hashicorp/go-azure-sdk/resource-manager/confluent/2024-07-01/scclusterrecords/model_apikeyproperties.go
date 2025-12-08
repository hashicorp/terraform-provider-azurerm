package scclusterrecords

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type APIKeyProperties struct {
	Metadata *SCMetadataEntity `json:"metadata,omitempty"`
	Spec     *APIKeySpecEntity `json:"spec,omitempty"`
}
