package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMLWebServiceFile struct {
	FilePath          string                 `json:"filePath"`
	LinkedServiceName LinkedServiceReference `json:"linkedServiceName"`
}
