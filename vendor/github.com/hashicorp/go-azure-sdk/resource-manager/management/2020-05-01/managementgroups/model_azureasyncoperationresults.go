package managementgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureAsyncOperationResults struct {
	Id         *string                        `json:"id,omitempty"`
	Name       *string                        `json:"name,omitempty"`
	Properties *ManagementGroupInfoProperties `json:"properties,omitempty"`
	Status     *string                        `json:"status,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
