package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SparkProfile struct {
	DefaultStorageUrl *string             `json:"defaultStorageUrl,omitempty"`
	MetastoreSpec     *SparkMetastoreSpec `json:"metastoreSpec,omitempty"`
	UserPluginsSpec   *SparkUserPlugins   `json:"userPluginsSpec,omitempty"`
}
