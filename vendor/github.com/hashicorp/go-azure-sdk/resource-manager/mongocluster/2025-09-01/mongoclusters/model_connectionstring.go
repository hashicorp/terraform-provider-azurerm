package mongoclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionString struct {
	ConnectionString *string `json:"connectionString,omitempty"`
	Description      *string `json:"description,omitempty"`
	Name             *string `json:"name,omitempty"`
}
