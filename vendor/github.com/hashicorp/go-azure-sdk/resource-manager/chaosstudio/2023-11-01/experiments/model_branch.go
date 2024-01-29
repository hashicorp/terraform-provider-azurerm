package experiments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Branch struct {
	Actions []Action `json:"actions"`
	Name    string   `json:"name"`
}
