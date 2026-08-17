package runtimeenvironment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuntimeEnvironmentUpdateProperties struct {
	DefaultPackages *map[string]string `json:"defaultPackages,omitempty"`
}
