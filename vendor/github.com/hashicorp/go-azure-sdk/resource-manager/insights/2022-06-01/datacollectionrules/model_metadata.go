package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Metadata struct {
	ProvisionedBy           *string `json:"provisionedBy,omitempty"`
	ProvisionedByResourceId *string `json:"provisionedByResourceId,omitempty"`
}
