package referencedatasets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReferenceDataSetCreateOrUpdateParameters struct {
	Location   string                             `json:"location"`
	Properties ReferenceDataSetCreationProperties `json:"properties"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
}
