package componentsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComponentPurgeBody struct {
	Filters []ComponentPurgeBodyFilters `json:"filters"`
	Table   string                      `json:"table"`
}
