package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynapseNotebookActivityTypeProperties struct {
	Conf                     *interface{}                                `json:"conf,omitempty"`
	ConfigurationType        *ConfigurationType                          `json:"configurationType,omitempty"`
	DriverSize               *interface{}                                `json:"driverSize,omitempty"`
	ExecutorSize             *interface{}                                `json:"executorSize,omitempty"`
	Notebook                 SynapseNotebookReference                    `json:"notebook"`
	NumExecutors             *interface{}                                `json:"numExecutors,omitempty"`
	Parameters               *map[string]NotebookParameter               `json:"parameters,omitempty"`
	SparkConfig              *map[string]interface{}                     `json:"sparkConfig,omitempty"`
	SparkPool                *BigDataPoolParametrizationReference        `json:"sparkPool,omitempty"`
	TargetSparkConfiguration *SparkConfigurationParametrizationReference `json:"targetSparkConfiguration,omitempty"`
}
