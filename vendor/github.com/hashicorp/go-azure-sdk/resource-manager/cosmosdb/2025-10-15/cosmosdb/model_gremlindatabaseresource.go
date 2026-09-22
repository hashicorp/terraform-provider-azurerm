package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GremlinDatabaseResource struct {
	CreateMode        *CreateMode            `json:"createMode,omitempty"`
	Id                string                 `json:"id"`
	RestoreParameters *RestoreParametersBase `json:"restoreParameters,omitempty"`
}
