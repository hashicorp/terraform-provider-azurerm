package managedhsmkeys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHsmKeyCreateParameters struct {
	Properties ManagedHsmKeyProperties `json:"properties"`
	Tags       *map[string]string      `json:"tags,omitempty"`
}
