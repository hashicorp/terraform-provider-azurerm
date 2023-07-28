package schedule

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StackEnsembleSettings struct {
	StackMetaLearnerKWargs          *interface{}          `json:"stackMetaLearnerKWargs,omitempty"`
	StackMetaLearnerTrainPercentage *float64              `json:"stackMetaLearnerTrainPercentage,omitempty"`
	StackMetaLearnerType            *StackMetaLearnerType `json:"stackMetaLearnerType,omitempty"`
}
