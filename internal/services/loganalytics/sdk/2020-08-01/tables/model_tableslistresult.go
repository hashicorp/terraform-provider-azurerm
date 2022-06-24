package tables

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TablesListResult struct {
	Value *[]Table `json:"value,omitempty"`
}
