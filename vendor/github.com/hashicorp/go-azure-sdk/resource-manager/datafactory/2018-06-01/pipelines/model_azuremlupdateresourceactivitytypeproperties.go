package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMLUpdateResourceActivityTypeProperties struct {
	TrainedModelFilePath          string                 `json:"trainedModelFilePath"`
	TrainedModelLinkedServiceName LinkedServiceReference `json:"trainedModelLinkedServiceName"`
	TrainedModelName              string                 `json:"trainedModelName"`
}
