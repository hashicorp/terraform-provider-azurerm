package apisbytag

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationTagResourceContractProperties struct {
	ApiName     *string `json:"apiName,omitempty"`
	ApiRevision *string `json:"apiRevision,omitempty"`
	ApiVersion  *string `json:"apiVersion,omitempty"`
	Description *string `json:"description,omitempty"`
	Id          *string `json:"id,omitempty"`
	Method      *string `json:"method,omitempty"`
	Name        *string `json:"name,omitempty"`
	UrlTemplate *string `json:"urlTemplate,omitempty"`
}
