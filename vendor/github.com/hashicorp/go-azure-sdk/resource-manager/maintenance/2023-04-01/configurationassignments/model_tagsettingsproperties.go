package configurationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagSettingsProperties struct {
	FilterOperator *TagOperators        `json:"filterOperator,omitempty"`
	Tags           *map[string][]string `json:"tags,omitempty"`
}
