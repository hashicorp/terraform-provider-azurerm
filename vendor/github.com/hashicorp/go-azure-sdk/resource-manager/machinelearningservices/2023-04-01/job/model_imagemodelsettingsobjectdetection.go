package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageModelSettingsObjectDetection struct {
	AdvancedSettings            *string                `json:"advancedSettings,omitempty"`
	AmsGradient                 *bool                  `json:"amsGradient,omitempty"`
	Augmentations               *string                `json:"augmentations,omitempty"`
	Beta1                       *float64               `json:"beta1,omitempty"`
	Beta2                       *float64               `json:"beta2,omitempty"`
	BoxDetectionsPerImage       *int64                 `json:"boxDetectionsPerImage,omitempty"`
	BoxScoreThreshold           *float64               `json:"boxScoreThreshold,omitempty"`
	CheckpointFrequency         *int64                 `json:"checkpointFrequency,omitempty"`
	CheckpointModel             JobInput               `json:"checkpointModel"`
	CheckpointRunId             *string                `json:"checkpointRunId,omitempty"`
	Distributed                 *bool                  `json:"distributed,omitempty"`
	EarlyStopping               *bool                  `json:"earlyStopping,omitempty"`
	EarlyStoppingDelay          *int64                 `json:"earlyStoppingDelay,omitempty"`
	EarlyStoppingPatience       *int64                 `json:"earlyStoppingPatience,omitempty"`
	EnableOnnxNormalization     *bool                  `json:"enableOnnxNormalization,omitempty"`
	EvaluationFrequency         *int64                 `json:"evaluationFrequency,omitempty"`
	GradientAccumulationStep    *int64                 `json:"gradientAccumulationStep,omitempty"`
	ImageSize                   *int64                 `json:"imageSize,omitempty"`
	LayersToFreeze              *int64                 `json:"layersToFreeze,omitempty"`
	LearningRate                *float64               `json:"learningRate,omitempty"`
	LearningRateScheduler       *LearningRateScheduler `json:"learningRateScheduler,omitempty"`
	MaxSize                     *int64                 `json:"maxSize,omitempty"`
	MinSize                     *int64                 `json:"minSize,omitempty"`
	ModelName                   *string                `json:"modelName,omitempty"`
	ModelSize                   *ModelSize             `json:"modelSize,omitempty"`
	Momentum                    *float64               `json:"momentum,omitempty"`
	MultiScale                  *bool                  `json:"multiScale,omitempty"`
	Nesterov                    *bool                  `json:"nesterov,omitempty"`
	NmsIouThreshold             *float64               `json:"nmsIouThreshold,omitempty"`
	NumberOfEpochs              *int64                 `json:"numberOfEpochs,omitempty"`
	NumberOfWorkers             *int64                 `json:"numberOfWorkers,omitempty"`
	Optimizer                   *StochasticOptimizer   `json:"optimizer,omitempty"`
	RandomSeed                  *int64                 `json:"randomSeed,omitempty"`
	StepLRGamma                 *float64               `json:"stepLRGamma,omitempty"`
	StepLRStepSize              *int64                 `json:"stepLRStepSize,omitempty"`
	TileGridSize                *string                `json:"tileGridSize,omitempty"`
	TileOverlapRatio            *float64               `json:"tileOverlapRatio,omitempty"`
	TilePredictionsNmsThreshold *float64               `json:"tilePredictionsNmsThreshold,omitempty"`
	TrainingBatchSize           *int64                 `json:"trainingBatchSize,omitempty"`
	ValidationBatchSize         *int64                 `json:"validationBatchSize,omitempty"`
	ValidationIouThreshold      *float64               `json:"validationIouThreshold,omitempty"`
	ValidationMetricType        *ValidationMetricType  `json:"validationMetricType,omitempty"`
	WarmupCosineLRCycles        *float64               `json:"warmupCosineLRCycles,omitempty"`
	WarmupCosineLRWarmupEpochs  *int64                 `json:"warmupCosineLRWarmupEpochs,omitempty"`
	WeightDecay                 *float64               `json:"weightDecay,omitempty"`
}

var _ json.Unmarshaler = &ImageModelSettingsObjectDetection{}

func (s *ImageModelSettingsObjectDetection) UnmarshalJSON(bytes []byte) error {
	type alias ImageModelSettingsObjectDetection
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ImageModelSettingsObjectDetection: %+v", err)
	}

	s.AdvancedSettings = decoded.AdvancedSettings
	s.AmsGradient = decoded.AmsGradient
	s.Augmentations = decoded.Augmentations
	s.Beta1 = decoded.Beta1
	s.Beta2 = decoded.Beta2
	s.BoxDetectionsPerImage = decoded.BoxDetectionsPerImage
	s.BoxScoreThreshold = decoded.BoxScoreThreshold
	s.CheckpointFrequency = decoded.CheckpointFrequency
	s.CheckpointRunId = decoded.CheckpointRunId
	s.Distributed = decoded.Distributed
	s.EarlyStopping = decoded.EarlyStopping
	s.EarlyStoppingDelay = decoded.EarlyStoppingDelay
	s.EarlyStoppingPatience = decoded.EarlyStoppingPatience
	s.EnableOnnxNormalization = decoded.EnableOnnxNormalization
	s.EvaluationFrequency = decoded.EvaluationFrequency
	s.GradientAccumulationStep = decoded.GradientAccumulationStep
	s.ImageSize = decoded.ImageSize
	s.LayersToFreeze = decoded.LayersToFreeze
	s.LearningRate = decoded.LearningRate
	s.LearningRateScheduler = decoded.LearningRateScheduler
	s.MaxSize = decoded.MaxSize
	s.MinSize = decoded.MinSize
	s.ModelName = decoded.ModelName
	s.ModelSize = decoded.ModelSize
	s.Momentum = decoded.Momentum
	s.MultiScale = decoded.MultiScale
	s.Nesterov = decoded.Nesterov
	s.NmsIouThreshold = decoded.NmsIouThreshold
	s.NumberOfEpochs = decoded.NumberOfEpochs
	s.NumberOfWorkers = decoded.NumberOfWorkers
	s.Optimizer = decoded.Optimizer
	s.RandomSeed = decoded.RandomSeed
	s.StepLRGamma = decoded.StepLRGamma
	s.StepLRStepSize = decoded.StepLRStepSize
	s.TileGridSize = decoded.TileGridSize
	s.TileOverlapRatio = decoded.TileOverlapRatio
	s.TilePredictionsNmsThreshold = decoded.TilePredictionsNmsThreshold
	s.TrainingBatchSize = decoded.TrainingBatchSize
	s.ValidationBatchSize = decoded.ValidationBatchSize
	s.ValidationIouThreshold = decoded.ValidationIouThreshold
	s.ValidationMetricType = decoded.ValidationMetricType
	s.WarmupCosineLRCycles = decoded.WarmupCosineLRCycles
	s.WarmupCosineLRWarmupEpochs = decoded.WarmupCosineLRWarmupEpochs
	s.WeightDecay = decoded.WeightDecay

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ImageModelSettingsObjectDetection into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["checkpointModel"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CheckpointModel' for 'ImageModelSettingsObjectDetection': %+v", err)
		}
		s.CheckpointModel = impl
	}
	return nil
}
