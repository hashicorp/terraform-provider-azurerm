package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Cardinality struct {
	BackendChain BackendChain `json:"backendChain"`
	Frontend     Frontend     `json:"frontend"`
}
