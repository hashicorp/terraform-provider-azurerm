package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Application struct {
	Name     *string `json:"name,omitempty"`
	Provider *string `json:"provider,omitempty"`
	Version  *string `json:"version,omitempty"`
}
