package createresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateResourceSupportedProperties struct {
	CreationSupported *bool   `json:"creationSupported,omitempty"`
	Name              *string `json:"name,omitempty"`
}
