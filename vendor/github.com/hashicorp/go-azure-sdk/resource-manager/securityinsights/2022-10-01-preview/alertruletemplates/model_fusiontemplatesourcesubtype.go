package alertruletemplates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FusionTemplateSourceSubType struct {
	SeverityFilter           FusionTemplateSubTypeSeverityFilter `json:"severityFilter"`
	SourceSubTypeDisplayName *string                             `json:"sourceSubTypeDisplayName,omitempty"`
	SourceSubTypeName        string                              `json:"sourceSubTypeName"`
}
