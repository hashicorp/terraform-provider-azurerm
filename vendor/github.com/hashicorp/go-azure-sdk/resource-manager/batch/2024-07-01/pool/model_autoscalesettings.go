package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoScaleSettings struct {
	EvaluationInterval *string `json:"evaluationInterval,omitempty"`
	Formula            string  `json:"formula"`
}
