package job

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableVerticalFeaturizationSettings struct {
	BlockedTransformers    *[]BlockedTransformers          `json:"blockedTransformers,omitempty"`
	ColumnNameAndTypes     *map[string]string              `json:"columnNameAndTypes,omitempty"`
	DatasetLanguage        *string                         `json:"datasetLanguage,omitempty"`
	EnableDnnFeaturization *bool                           `json:"enableDnnFeaturization,omitempty"`
	Mode                   *FeaturizationMode              `json:"mode,omitempty"`
	TransformerParams      *map[string][]ColumnTransformer `json:"transformerParams,omitempty"`
}
