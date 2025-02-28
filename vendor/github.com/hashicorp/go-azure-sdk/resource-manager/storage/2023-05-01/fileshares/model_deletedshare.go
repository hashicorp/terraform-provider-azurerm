package fileshares

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedShare struct {
	DeletedShareName    string `json:"deletedShareName"`
	DeletedShareVersion string `json:"deletedShareVersion"`
}
