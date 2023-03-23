package exports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Export struct {
	ETag       *string           `json:"eTag,omitempty"`
	Id         *string           `json:"id,omitempty"`
	Name       *string           `json:"name,omitempty"`
	Properties *ExportProperties `json:"properties,omitempty"`
	Type       *string           `json:"type,omitempty"`
}
