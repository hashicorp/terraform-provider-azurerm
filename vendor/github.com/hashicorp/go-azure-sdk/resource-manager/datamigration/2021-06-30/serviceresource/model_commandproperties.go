package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommandProperties struct {
	CommandType string        `json:"commandType"`
	Errors      *[]ODataError `json:"errors,omitempty"`
	State       *CommandState `json:"state,omitempty"`
}
