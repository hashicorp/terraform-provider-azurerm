package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationExtension struct {
	Publisher *string `json:"publisher,omitempty"`
	Type      *string `json:"type,omitempty"`
}
