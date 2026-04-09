package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationPackageContact struct {
	ContactName *string `json:"contactName,omitempty"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
}
