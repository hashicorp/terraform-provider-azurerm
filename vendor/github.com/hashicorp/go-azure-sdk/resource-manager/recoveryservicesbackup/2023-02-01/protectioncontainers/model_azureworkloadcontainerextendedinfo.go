package protectioncontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureWorkloadContainerExtendedInfo struct {
	HostServerName *string                 `json:"hostServerName,omitempty"`
	InquiryInfo    *InquiryInfo            `json:"inquiryInfo,omitempty"`
	NodesList      *[]DistributedNodesInfo `json:"nodesList,omitempty"`
}
