package integrationruntimeenableinteractivequery

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeDataProxyProperties struct {
	ConnectVia           *EntityReference `json:"connectVia,omitempty"`
	Path                 *string          `json:"path,omitempty"`
	StagingLinkedService *EntityReference `json:"stagingLinkedService,omitempty"`
}
