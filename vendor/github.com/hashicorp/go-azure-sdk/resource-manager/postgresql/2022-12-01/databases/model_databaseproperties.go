package databases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseProperties struct {
	Charset   *string `json:"charset,omitempty"`
	Collation *string `json:"collation,omitempty"`
}
