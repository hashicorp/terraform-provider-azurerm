package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynapseSparkJobActivityTypeProperties struct {
	Args                     *[]interface{}                              `json:"args,omitempty"`
	ClassName                *string                                     `json:"className,omitempty"`
	Conf                     *interface{}                                `json:"conf,omitempty"`
	ConfigurationType        *ConfigurationType                          `json:"configurationType,omitempty"`
	DriverSize               *string                                     `json:"driverSize,omitempty"`
	ExecutorSize             *string                                     `json:"executorSize,omitempty"`
	File                     *string                                     `json:"file,omitempty"`
	Files                    *[]string                                   `json:"files,omitempty"`
	FilesV2                  *[]string                                   `json:"filesV2,omitempty"`
	NumExecutors             *int64                                      `json:"numExecutors,omitempty"`
	PythonCodeReference      *[]string                                   `json:"pythonCodeReference,omitempty"`
	ScanFolder               *bool                                       `json:"scanFolder,omitempty"`
	SparkConfig              *map[string]string                          `json:"sparkConfig,omitempty"`
	SparkJob                 SynapseSparkJobReference                    `json:"sparkJob"`
	TargetBigDataPool        *BigDataPoolParametrizationReference        `json:"targetBigDataPool,omitempty"`
	TargetSparkConfiguration *SparkConfigurationParametrizationReference `json:"targetSparkConfiguration,omitempty"`
}
