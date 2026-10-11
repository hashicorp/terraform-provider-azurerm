package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePurgeBody struct {
	Filters []WorkspacePurgeBodyFilters `json:"filters"`
	Table   string                      `json:"table"`
}
