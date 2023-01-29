package deploymentscripts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentVariable struct {
	Name        string  `json:"name"`
	SecureValue *string `json:"secureValue,omitempty"`
	Value       *string `json:"value,omitempty"`
}
