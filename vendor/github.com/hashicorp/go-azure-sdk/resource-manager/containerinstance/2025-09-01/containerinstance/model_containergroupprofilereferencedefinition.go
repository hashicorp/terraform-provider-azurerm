package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupProfileReferenceDefinition struct {
	Id       *string `json:"id,omitempty"`
	Revision *int64  `json:"revision,omitempty"`
}
