package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionAppConfig struct {
	Deployment          *FunctionsDeployment          `json:"deployment,omitempty"`
	Runtime             *FunctionsRuntime             `json:"runtime,omitempty"`
	ScaleAndConcurrency *FunctionsScaleAndConcurrency `json:"scaleAndConcurrency,omitempty"`
}
