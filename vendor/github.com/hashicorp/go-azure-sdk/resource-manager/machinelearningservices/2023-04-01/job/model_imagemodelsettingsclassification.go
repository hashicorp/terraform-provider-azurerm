package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageModelSettingsClassification struct {
	AdvancedSettings           *string                `json:"advancedSettings,omitempty"`
	AmsGradient                *bool                  `json:"amsGradient,omitempty"`
	Augmentations              *string                `json:"augmentations,omitempty"`
	Beta1                      *float64               `json:"beta1,omitempty"`
	Beta2                      *float64               `json:"beta2,omitempty"`
	CheckpointFrequency        *int64                 `json:"checkpointFrequency,omitempty"`
	CheckpointModel            JobInput               `json:"checkpointModel"`
	CheckpointRunId            *string                `json:"checkpointRunId,omitempty"`
	Distributed                *bool                  `json:"distributed,omitempty"`
	EarlyStopping              *bool                  `json:"earlyStopping,omitempty"`
	EarlyStoppingDelay         *int64                 `json:"earlyStoppingDelay,omitempty"`
	EarlyStoppingPatience      *int64                 `json:"earlyStoppingPatience,omitempty"`
	EnableOnnxNormalization    *bool                  `json:"enableOnnxNormalization,omitempty"`
	EvaluationFrequency        *int64                 `json:"evaluationFrequency,omitempty"`
	GradientAccumulationStep   *int64                 `json:"gradientAccumulationStep,omitempty"`
	LayersToFreeze             *int64                 `json:"layersToFreeze,omitempty"`
	LearningRate               *float64               `json:"learningRate,omitempty"`
	LearningRateScheduler      *LearningRateScheduler `json:"learningRateScheduler,omitempty"`
	ModelName                  *string                `json:"modelName,omitempty"`
	Momentum                   *float64               `json:"momentum,omitempty"`
	Nesterov                   *bool                  `json:"nesterov,omitempty"`
	NumberOfEpochs             *int64                 `json:"numberOfEpochs,omitempty"`
	NumberOfWorkers            *int64                 `json:"numberOfWorkers,omitempty"`
	Optimizer                  *StochasticOptimizer   `json:"optimizer,omitempty"`
	RandomSeed                 *int64                 `json:"randomSeed,omitempty"`
	StepLRGamma                *float64               `json:"stepLRGamma,omitempty"`
	StepLRStepSize             *int64                 `json:"stepLRStepSize,omitempty"`
	TrainingBatchSize          *int64                 `json:"trainingBatchSize,omitempty"`
	TrainingCropSize           *int64                 `json:"trainingCropSize,omitempty"`
	ValidationBatchSize        *int64                 `json:"validationBatchSize,omitempty"`
	ValidationCropSize         *int64                 `json:"validationCropSize,omitempty"`
	ValidationResizeSize       *int64                 `json:"validationResizeSize,omitempty"`
	WarmupCosineLRCycles       *float64               `json:"warmupCosineLRCycles,omitempty"`
	WarmupCosineLRWarmupEpochs *int64                 `json:"warmupCosineLRWarmupEpochs,omitempty"`
	WeightDecay                *float64               `json:"weightDecay,omitempty"`
	WeightedLoss               *int64                 `json:"weightedLoss,omitempty"`
}

var _ json.Unmarshaler = &ImageModelSettingsClassification{}

func (s *ImageModelSettingsClassification) UnmarshalJSON(bytes []byte) error {
	type alias ImageModelSettingsClassification
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ImageModelSettingsClassification: %+v", err)
	}

	s.AdvancedSettings = decoded.AdvancedSettings
	s.AmsGradient = decoded.AmsGradient
	s.Augmentations = decoded.Augmentations
	s.Beta1 = decoded.Beta1
	s.Beta2 = decoded.Beta2
	s.CheckpointFrequency = decoded.CheckpointFrequency
	s.CheckpointRunId = decoded.CheckpointRunId
	s.Distributed = decoded.Distributed
	s.EarlyStopping = decoded.EarlyStopping
	s.EarlyStoppingDelay = decoded.EarlyStoppingDelay
	s.EarlyStoppingPatience = decoded.EarlyStoppingPatience
	s.EnableOnnxNormalization = decoded.EnableOnnxNormalization
	s.EvaluationFrequency = decoded.EvaluationFrequency
	s.GradientAccumulationStep = decoded.GradientAccumulationStep
	s.LayersToFreeze = decoded.LayersToFreeze
	s.LearningRate = decoded.LearningRate
	s.LearningRateScheduler = decoded.LearningRateScheduler
	s.ModelName = decoded.ModelName
	s.Momentum = decoded.Momentum
	s.Nesterov = decoded.Nesterov
	s.NumberOfEpochs = decoded.NumberOfEpochs
	s.NumberOfWorkers = decoded.NumberOfWorkers
	s.Optimizer = decoded.Optimizer
	s.RandomSeed = decoded.RandomSeed
	s.StepLRGamma = decoded.StepLRGamma
	s.StepLRStepSize = decoded.StepLRStepSize
	s.TrainingBatchSize = decoded.TrainingBatchSize
	s.TrainingCropSize = decoded.TrainingCropSize
	s.ValidationBatchSize = decoded.ValidationBatchSize
	s.ValidationCropSize = decoded.ValidationCropSize
	s.ValidationResizeSize = decoded.ValidationResizeSize
	s.WarmupCosineLRCycles = decoded.WarmupCosineLRCycles
	s.WarmupCosineLRWarmupEpochs = decoded.WarmupCosineLRWarmupEpochs
	s.WeightDecay = decoded.WeightDecay
	s.WeightedLoss = decoded.WeightedLoss

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ImageModelSettingsClassification into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["checkpointModel"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CheckpointModel' for 'ImageModelSettingsClassification': %+v", err)
		}
		s.CheckpointModel = impl
	}
	return nil
}
