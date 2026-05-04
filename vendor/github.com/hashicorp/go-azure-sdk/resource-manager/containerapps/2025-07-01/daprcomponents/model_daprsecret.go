package daprcomponents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaprSecret struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}
