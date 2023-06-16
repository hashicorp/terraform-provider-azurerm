package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentSetting struct {
	Name  string  `json:"name"`
	Value *string `json:"value,omitempty"`
}
