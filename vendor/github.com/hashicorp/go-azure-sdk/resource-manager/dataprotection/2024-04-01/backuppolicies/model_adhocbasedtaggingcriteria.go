package backuppolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdhocBasedTaggingCriteria struct {
	TagInfo *RetentionTag `json:"tagInfo,omitempty"`
}
