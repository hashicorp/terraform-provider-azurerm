package account

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountEndpoints struct {
	Catalog  *string `json:"catalog,omitempty"`
	Guardian *string `json:"guardian,omitempty"`
	Scan     *string `json:"scan,omitempty"`
}
