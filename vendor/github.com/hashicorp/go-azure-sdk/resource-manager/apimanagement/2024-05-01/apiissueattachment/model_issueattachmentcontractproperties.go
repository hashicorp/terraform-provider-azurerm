package apiissueattachment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IssueAttachmentContractProperties struct {
	Content       string `json:"content"`
	ContentFormat string `json:"contentFormat"`
	Title         string `json:"title"`
}
