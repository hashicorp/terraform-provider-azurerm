package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerNameItem struct {
	FullyQualifiedDomainName *string `json:"fullyQualifiedDomainName,omitempty"`
	Name                     *string `json:"name,omitempty"`
}
