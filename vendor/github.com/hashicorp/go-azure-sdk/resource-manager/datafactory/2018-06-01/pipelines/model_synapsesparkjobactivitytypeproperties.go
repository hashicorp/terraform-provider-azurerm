package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynapseSparkJobActivityTypeProperties struct {
	Args                     *[]interface{}                              `json:"args,omitempty"`
	ClassName                *interface{}                                `json:"className,omitempty"`
	Conf                     *interface{}                                `json:"conf,omitempty"`
	ConfigurationType        *ConfigurationType                          `json:"configurationType,omitempty"`
	DriverSize               *interface{}                                `json:"driverSize,omitempty"`
	ExecutorSize             *interface{}                                `json:"executorSize,omitempty"`
	File                     *interface{}                                `json:"file,omitempty"`
	Files                    *[]interface{}                              `json:"files,omitempty"`
	FilesV2                  *[]interface{}                              `json:"filesV2,omitempty"`
	NumExecutors             *int64                                      `json:"numExecutors,omitempty"`
	PythonCodeReference      *[]interface{}                              `json:"pythonCodeReference,omitempty"`
	ScanFolder               *bool                                       `json:"scanFolder,omitempty"`
	SparkConfig              *map[string]interface{}                     `json:"sparkConfig,omitempty"`
	SparkJob                 SynapseSparkJobReference                    `json:"sparkJob"`
	TargetBigDataPool        *BigDataPoolParametrizationReference        `json:"targetBigDataPool,omitempty"`
	TargetSparkConfiguration *SparkConfigurationParametrizationReference `json:"targetSparkConfiguration,omitempty"`
}
