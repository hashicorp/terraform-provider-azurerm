package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomain struct {
	Name             string `json:"name"`
	UseSubDomainName *bool  `json:"useSubDomainName,omitempty"`
}
