package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMLUpdateResourceActivityTypeProperties struct {
	TrainedModelFilePath          interface{}            `json:"trainedModelFilePath"`
	TrainedModelLinkedServiceName LinkedServiceReference `json:"trainedModelLinkedServiceName"`
	TrainedModelName              interface{}            `json:"trainedModelName"`
}
