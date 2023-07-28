package job

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageModelDistributionSettingsObjectDetection struct {
	AmsGradient                 *string `json:"amsGradient,omitempty"`
	Augmentations               *string `json:"augmentations,omitempty"`
	Beta1                       *string `json:"beta1,omitempty"`
	Beta2                       *string `json:"beta2,omitempty"`
	BoxDetectionsPerImage       *string `json:"boxDetectionsPerImage,omitempty"`
	BoxScoreThreshold           *string `json:"boxScoreThreshold,omitempty"`
	Distributed                 *string `json:"distributed,omitempty"`
	EarlyStopping               *string `json:"earlyStopping,omitempty"`
	EarlyStoppingDelay          *string `json:"earlyStoppingDelay,omitempty"`
	EarlyStoppingPatience       *string `json:"earlyStoppingPatience,omitempty"`
	EnableOnnxNormalization     *string `json:"enableOnnxNormalization,omitempty"`
	EvaluationFrequency         *string `json:"evaluationFrequency,omitempty"`
	GradientAccumulationStep    *string `json:"gradientAccumulationStep,omitempty"`
	ImageSize                   *string `json:"imageSize,omitempty"`
	LayersToFreeze              *string `json:"layersToFreeze,omitempty"`
	LearningRate                *string `json:"learningRate,omitempty"`
	LearningRateScheduler       *string `json:"learningRateScheduler,omitempty"`
	MaxSize                     *string `json:"maxSize,omitempty"`
	MinSize                     *string `json:"minSize,omitempty"`
	ModelName                   *string `json:"modelName,omitempty"`
	ModelSize                   *string `json:"modelSize,omitempty"`
	Momentum                    *string `json:"momentum,omitempty"`
	MultiScale                  *string `json:"multiScale,omitempty"`
	Nesterov                    *string `json:"nesterov,omitempty"`
	NmsIouThreshold             *string `json:"nmsIouThreshold,omitempty"`
	NumberOfEpochs              *string `json:"numberOfEpochs,omitempty"`
	NumberOfWorkers             *string `json:"numberOfWorkers,omitempty"`
	Optimizer                   *string `json:"optimizer,omitempty"`
	RandomSeed                  *string `json:"randomSeed,omitempty"`
	StepLRGamma                 *string `json:"stepLRGamma,omitempty"`
	StepLRStepSize              *string `json:"stepLRStepSize,omitempty"`
	TileGridSize                *string `json:"tileGridSize,omitempty"`
	TileOverlapRatio            *string `json:"tileOverlapRatio,omitempty"`
	TilePredictionsNmsThreshold *string `json:"tilePredictionsNmsThreshold,omitempty"`
	TrainingBatchSize           *string `json:"trainingBatchSize,omitempty"`
	ValidationBatchSize         *string `json:"validationBatchSize,omitempty"`
	ValidationIouThreshold      *string `json:"validationIouThreshold,omitempty"`
	ValidationMetricType        *string `json:"validationMetricType,omitempty"`
	WarmupCosineLRCycles        *string `json:"warmupCosineLRCycles,omitempty"`
	WarmupCosineLRWarmupEpochs  *string `json:"warmupCosineLRWarmupEpochs,omitempty"`
	WeightDecay                 *string `json:"weightDecay,omitempty"`
}
