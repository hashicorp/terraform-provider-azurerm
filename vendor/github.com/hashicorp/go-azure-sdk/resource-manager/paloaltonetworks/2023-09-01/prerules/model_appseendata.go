package prerules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppSeenData struct {
	AppSeenList []AppSeenInfo `json:"appSeenList"`
	Count       int64         `json:"count"`
}
