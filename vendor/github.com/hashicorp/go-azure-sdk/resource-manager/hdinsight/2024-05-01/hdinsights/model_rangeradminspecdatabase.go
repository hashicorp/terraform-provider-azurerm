package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RangerAdminSpecDatabase struct {
	Host              string  `json:"host"`
	Name              string  `json:"name"`
	PasswordSecretRef *string `json:"passwordSecretRef,omitempty"`
	Username          *string `json:"username,omitempty"`
}
