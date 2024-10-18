package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterKey struct {
	Name    *string `json:"name,omitempty"`
	OrderBy *string `json:"orderBy,omitempty"`
}
