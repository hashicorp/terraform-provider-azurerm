package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeepCreatedCustomDomain struct {
	Name       string                             `json:"name"`
	Properties *DeepCreatedCustomDomainProperties `json:"properties,omitempty"`
}
