package brokerauthorization

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StateStoreResourceRule struct {
	KeyType StateStoreResourceKeyTypes          `json:"keyType"`
	Keys    []string                            `json:"keys"`
	Method  StateStoreResourceDefinitionMethods `json:"method"`
}
