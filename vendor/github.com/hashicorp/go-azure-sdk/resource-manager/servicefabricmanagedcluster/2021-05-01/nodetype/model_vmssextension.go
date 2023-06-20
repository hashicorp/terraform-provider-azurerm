package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMSSExtension struct {
	Name       string                  `json:"name"`
	Properties VMSSExtensionProperties `json:"properties"`
}
