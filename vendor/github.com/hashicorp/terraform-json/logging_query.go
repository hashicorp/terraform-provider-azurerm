// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tfjson

const (
	MessageListStart         LogMessageType = "list_start"
	MessageListResourceFound LogMessageType = "list_resource_found"
	MessageListComplete      LogMessageType = "list_complete"
)

// ListStartMessage represents "query" result message of type "list_start"
type ListStartMessage struct {
	baseLogMessage
	ListStart ListStartData `json:"list_start"`
}

type ListStartData struct {
	Address      string         `json:"address"`
	ResourceType string         `json:"resource_type"`
	InputConfig  map[string]any `json:"input_config,omitempty"`
}

// ListResourceFoundMessage represents "query" result message of type "list_resource_found"
type ListResourceFoundMessage struct {
	baseLogMessage
	ListResourceFound ListResourceFoundData `json:"list_resource_found"`
}

type ListResourceFoundData struct {
	Address         string         `json:"address"`
	DisplayName     string         `json:"display_name"`
	Identity        map[string]any `json:"identity"`
	IdentityVersion int64          `json:"identity_version"`
	ResourceType    string         `json:"resource_type"`
	ResourceObject  map[string]any `json:"resource_object,omitempty"`
	Config          string         `json:"config,omitempty"`
	ImportConfig    string         `json:"import_config,omitempty"`
}

// ListCompleteMessage represents "query" result message of type "list_complete"
type ListCompleteMessage struct {
	baseLogMessage
	ListComplete ListCompleteData `json:"list_complete"`
}

type ListCompleteData struct {
	Address      string `json:"address"`
	ResourceType string `json:"resource_type"`
	Total        int    `json:"total"`
}
