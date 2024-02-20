package tags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagsPatchResource struct {
	Operation  *TagsPatchOperation `json:"operation,omitempty"`
	Properties *Tags               `json:"properties,omitempty"`
}
