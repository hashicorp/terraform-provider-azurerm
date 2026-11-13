package sensitivitylabels

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SensitivityLabelProperties struct {
	ClientClassificationSource *ClientClassificationSource `json:"clientClassificationSource,omitempty"`
	ColumnName                 *string                     `json:"columnName,omitempty"`
	InformationType            *string                     `json:"informationType,omitempty"`
	InformationTypeId          *string                     `json:"informationTypeId,omitempty"`
	IsDisabled                 *bool                       `json:"isDisabled,omitempty"`
	LabelId                    *string                     `json:"labelId,omitempty"`
	LabelName                  *string                     `json:"labelName,omitempty"`
	Rank                       *SensitivityLabelRank       `json:"rank,omitempty"`
	SchemaName                 *string                     `json:"schemaName,omitempty"`
	TableName                  *string                     `json:"tableName,omitempty"`
}
