package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateRequest struct {
	Location   string                `json:"location"`
	Name       string                `json:"name"`
	Properties ValidateProperties    `json:"properties"`
	Type       ValidateResourceTypes `json:"type"`
}
