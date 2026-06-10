package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Activity interface {
	Activity() BaseActivityImpl
}

var _ Activity = BaseActivityImpl{}

type BaseActivityImpl struct {
	DependsOn        *[]ActivityDependency     `json:"dependsOn,omitempty"`
	Description      *string                   `json:"description,omitempty"`
	Name             string                    `json:"name"`
	OnInactiveMarkAs *ActivityOnInactiveMarkAs `json:"onInactiveMarkAs,omitempty"`
	State            *ActivityState            `json:"state,omitempty"`
	Type             string                    `json:"type"`
	UserProperties   *[]UserProperty           `json:"userProperties,omitempty"`
}

func (s BaseActivityImpl) Activity() BaseActivityImpl {
	return s
}

var _ Activity = RawActivityImpl{}

// RawActivityImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawActivityImpl struct {
	activity BaseActivityImpl
	Type     string
	Values   map[string]interface{}
}

func (s RawActivityImpl) Activity() BaseActivityImpl {
	return s.activity
}

func UnmarshalActivityImplementation(input []byte) (Activity, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Activity into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AppendVariable") {
		var out AppendVariableActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AppendVariableActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataExplorerCommand") {
		var out AzureDataExplorerCommandActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataExplorerCommandActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureFunctionActivity") {
		var out AzureFunctionActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFunctionActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMLBatchExecution") {
		var out AzureMLBatchExecutionActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMLBatchExecutionActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMLExecutePipeline") {
		var out AzureMLExecutePipelineActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMLExecutePipelineActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMLUpdateResource") {
		var out AzureMLUpdateResourceActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMLUpdateResourceActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Container") {
		var out ControlActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ControlActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Copy") {
		var out CopyActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CopyActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Custom") {
		var out CustomActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DataLakeAnalyticsU-SQL") {
		var out DataLakeAnalyticsUSQLActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataLakeAnalyticsUSQLActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabricksJob") {
		var out DatabricksJobActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DatabricksJobActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabricksNotebook") {
		var out DatabricksNotebookActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DatabricksNotebookActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabricksSparkJar") {
		var out DatabricksSparkJarActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DatabricksSparkJarActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DatabricksSparkPython") {
		var out DatabricksSparkPythonActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DatabricksSparkPythonActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Delete") {
		var out DeleteActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeleteActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ExecuteDataFlow") {
		var out ExecuteDataFlowActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExecuteDataFlowActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ExecutePipeline") {
		var out ExecutePipelineActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExecutePipelineActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ExecuteSSISPackage") {
		var out ExecuteSSISPackageActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExecuteSSISPackageActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ExecuteWranglingDataflow") {
		var out ExecuteWranglingDataflowActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExecuteWranglingDataflowActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Execution") {
		var out ExecutionActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExecutionActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Fail") {
		var out FailActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FailActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Filter") {
		var out FilterActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FilterActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ForEach") {
		var out ForEachActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ForEachActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GetMetadata") {
		var out GetMetadataActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GetMetadataActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HDInsightHive") {
		var out HDInsightHiveActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HDInsightHiveActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HDInsightMapReduce") {
		var out HDInsightMapReduceActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HDInsightMapReduceActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HDInsightPig") {
		var out HDInsightPigActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HDInsightPigActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HDInsightSpark") {
		var out HDInsightSparkActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HDInsightSparkActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HDInsightStreaming") {
		var out HDInsightStreamingActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HDInsightStreamingActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IfCondition") {
		var out IfConditionActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IfConditionActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Lookup") {
		var out LookupActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LookupActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Script") {
		var out ScriptActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ScriptActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SetVariable") {
		var out SetVariableActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SetVariableActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlServerStoredProcedure") {
		var out SqlServerStoredProcedureActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlServerStoredProcedureActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Switch") {
		var out SwitchActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SwitchActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SynapseNotebook") {
		var out SynapseNotebookActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SynapseNotebookActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SparkJob") {
		var out SynapseSparkJobDefinitionActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SynapseSparkJobDefinitionActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Until") {
		var out UntilActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UntilActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Validation") {
		var out ValidationActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ValidationActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Wait") {
		var out WaitActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WaitActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WebActivity") {
		var out WebActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebActivity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WebHook") {
		var out WebHookActivity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebHookActivity: %+v", err)
		}
		return out, nil
	}

	var parent BaseActivityImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseActivityImpl: %+v", err)
	}

	return RawActivityImpl{
		activity: parent,
		Type:     value,
		Values:   temp,
	}, nil

}
