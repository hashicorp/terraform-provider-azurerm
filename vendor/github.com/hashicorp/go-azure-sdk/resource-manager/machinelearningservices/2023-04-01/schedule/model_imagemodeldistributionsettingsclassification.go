package schedule

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageModelDistributionSettingsClassification struct {
	AmsGradient                *string `json:"amsGradient,omitempty"`
	Augmentations              *string `json:"augmentations,omitempty"`
	Beta1                      *string `json:"beta1,omitempty"`
	Beta2                      *string `json:"beta2,omitempty"`
	Distributed                *string `json:"distributed,omitempty"`
	EarlyStopping              *string `json:"earlyStopping,omitempty"`
	EarlyStoppingDelay         *string `json:"earlyStoppingDelay,omitempty"`
	EarlyStoppingPatience      *string `json:"earlyStoppingPatience,omitempty"`
	EnableOnnxNormalization    *string `json:"enableOnnxNormalization,omitempty"`
	EvaluationFrequency        *string `json:"evaluationFrequency,omitempty"`
	GradientAccumulationStep   *string `json:"gradientAccumulationStep,omitempty"`
	LayersToFreeze             *string `json:"layersToFreeze,omitempty"`
	LearningRate               *string `json:"learningRate,omitempty"`
	LearningRateScheduler      *string `json:"learningRateScheduler,omitempty"`
	ModelName                  *string `json:"modelName,omitempty"`
	Momentum                   *string `json:"momentum,omitempty"`
	Nesterov                   *string `json:"nesterov,omitempty"`
	NumberOfEpochs             *string `json:"numberOfEpochs,omitempty"`
	NumberOfWorkers            *string `json:"numberOfWorkers,omitempty"`
	Optimizer                  *string `json:"optimizer,omitempty"`
	RandomSeed                 *string `json:"randomSeed,omitempty"`
	StepLRGamma                *string `json:"stepLRGamma,omitempty"`
	StepLRStepSize             *string `json:"stepLRStepSize,omitempty"`
	TrainingBatchSize          *string `json:"trainingBatchSize,omitempty"`
	TrainingCropSize           *string `json:"trainingCropSize,omitempty"`
	ValidationBatchSize        *string `json:"validationBatchSize,omitempty"`
	ValidationCropSize         *string `json:"validationCropSize,omitempty"`
	ValidationResizeSize       *string `json:"validationResizeSize,omitempty"`
	WarmupCosineLRCycles       *string `json:"warmupCosineLRCycles,omitempty"`
	WarmupCosineLRWarmupEpochs *string `json:"warmupCosineLRWarmupEpochs,omitempty"`
	WeightDecay                *string `json:"weightDecay,omitempty"`
	WeightedLoss               *string `json:"weightedLoss,omitempty"`
}
