package deviceupdates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagUpdate struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
