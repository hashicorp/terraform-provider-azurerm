package job

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForecastingTrainingSettings struct {
	AllowedTrainingAlgorithms    *[]ForecastingModels   `json:"allowedTrainingAlgorithms,omitempty"`
	BlockedTrainingAlgorithms    *[]ForecastingModels   `json:"blockedTrainingAlgorithms,omitempty"`
	EnableDnnTraining            *bool                  `json:"enableDnnTraining,omitempty"`
	EnableModelExplainability    *bool                  `json:"enableModelExplainability,omitempty"`
	EnableOnnxCompatibleModels   *bool                  `json:"enableOnnxCompatibleModels,omitempty"`
	EnableStackEnsemble          *bool                  `json:"enableStackEnsemble,omitempty"`
	EnableVoteEnsemble           *bool                  `json:"enableVoteEnsemble,omitempty"`
	EnsembleModelDownloadTimeout *string                `json:"ensembleModelDownloadTimeout,omitempty"`
	StackEnsembleSettings        *StackEnsembleSettings `json:"stackEnsembleSettings,omitempty"`
}
