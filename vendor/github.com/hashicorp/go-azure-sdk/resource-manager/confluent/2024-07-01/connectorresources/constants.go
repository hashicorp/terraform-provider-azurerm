package connectorresources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthType string

const (
	AuthTypeKAFKAAPIKEY    AuthType = "KAFKA_API_KEY"
	AuthTypeSERVICEACCOUNT AuthType = "SERVICE_ACCOUNT"
)

func PossibleValuesForAuthType() []string {
	return []string{
		string(AuthTypeKAFKAAPIKEY),
		string(AuthTypeSERVICEACCOUNT),
	}
}

func (s *AuthType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthType(input string) (*AuthType, error) {
	vals := map[string]AuthType{
		"kafka_api_key":   AuthTypeKAFKAAPIKEY,
		"service_account": AuthTypeSERVICEACCOUNT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthType(input)
	return &out, nil
}

type ConnectorClass string

const (
	ConnectorClassAZUREBLOBSINK   ConnectorClass = "AZUREBLOBSINK"
	ConnectorClassAZUREBLOBSOURCE ConnectorClass = "AZUREBLOBSOURCE"
)

func PossibleValuesForConnectorClass() []string {
	return []string{
		string(ConnectorClassAZUREBLOBSINK),
		string(ConnectorClassAZUREBLOBSOURCE),
	}
}

func (s *ConnectorClass) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectorClass(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectorClass(input string) (*ConnectorClass, error) {
	vals := map[string]ConnectorClass{
		"azureblobsink":   ConnectorClassAZUREBLOBSINK,
		"azureblobsource": ConnectorClassAZUREBLOBSOURCE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectorClass(input)
	return &out, nil
}

type ConnectorServiceType string

const (
	ConnectorServiceTypeAzureBlobStorageSinkConnector      ConnectorServiceType = "AzureBlobStorageSinkConnector"
	ConnectorServiceTypeAzureBlobStorageSourceConnector    ConnectorServiceType = "AzureBlobStorageSourceConnector"
	ConnectorServiceTypeAzureCosmosDBSinkConnector         ConnectorServiceType = "AzureCosmosDBSinkConnector"
	ConnectorServiceTypeAzureCosmosDBSourceConnector       ConnectorServiceType = "AzureCosmosDBSourceConnector"
	ConnectorServiceTypeAzureSynapseAnalyticsSinkConnector ConnectorServiceType = "AzureSynapseAnalyticsSinkConnector"
)

func PossibleValuesForConnectorServiceType() []string {
	return []string{
		string(ConnectorServiceTypeAzureBlobStorageSinkConnector),
		string(ConnectorServiceTypeAzureBlobStorageSourceConnector),
		string(ConnectorServiceTypeAzureCosmosDBSinkConnector),
		string(ConnectorServiceTypeAzureCosmosDBSourceConnector),
		string(ConnectorServiceTypeAzureSynapseAnalyticsSinkConnector),
	}
}

func (s *ConnectorServiceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectorServiceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectorServiceType(input string) (*ConnectorServiceType, error) {
	vals := map[string]ConnectorServiceType{
		"azureblobstoragesinkconnector":      ConnectorServiceTypeAzureBlobStorageSinkConnector,
		"azureblobstoragesourceconnector":    ConnectorServiceTypeAzureBlobStorageSourceConnector,
		"azurecosmosdbsinkconnector":         ConnectorServiceTypeAzureCosmosDBSinkConnector,
		"azurecosmosdbsourceconnector":       ConnectorServiceTypeAzureCosmosDBSourceConnector,
		"azuresynapseanalyticssinkconnector": ConnectorServiceTypeAzureSynapseAnalyticsSinkConnector,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectorServiceType(input)
	return &out, nil
}

type ConnectorStatus string

const (
	ConnectorStatusFAILED       ConnectorStatus = "FAILED"
	ConnectorStatusPAUSED       ConnectorStatus = "PAUSED"
	ConnectorStatusPROVISIONING ConnectorStatus = "PROVISIONING"
	ConnectorStatusRUNNING      ConnectorStatus = "RUNNING"
)

func PossibleValuesForConnectorStatus() []string {
	return []string{
		string(ConnectorStatusFAILED),
		string(ConnectorStatusPAUSED),
		string(ConnectorStatusPROVISIONING),
		string(ConnectorStatusRUNNING),
	}
}

func (s *ConnectorStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectorStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectorStatus(input string) (*ConnectorStatus, error) {
	vals := map[string]ConnectorStatus{
		"failed":       ConnectorStatusFAILED,
		"paused":       ConnectorStatusPAUSED,
		"provisioning": ConnectorStatusPROVISIONING,
		"running":      ConnectorStatusRUNNING,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectorStatus(input)
	return &out, nil
}

type ConnectorType string

const (
	ConnectorTypeSINK   ConnectorType = "SINK"
	ConnectorTypeSOURCE ConnectorType = "SOURCE"
)

func PossibleValuesForConnectorType() []string {
	return []string{
		string(ConnectorTypeSINK),
		string(ConnectorTypeSOURCE),
	}
}

func (s *ConnectorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectorType(input string) (*ConnectorType, error) {
	vals := map[string]ConnectorType{
		"sink":   ConnectorTypeSINK,
		"source": ConnectorTypeSOURCE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectorType(input)
	return &out, nil
}

type DataFormatType string

const (
	DataFormatTypeAVRO     DataFormatType = "AVRO"
	DataFormatTypeBYTES    DataFormatType = "BYTES"
	DataFormatTypeJSON     DataFormatType = "JSON"
	DataFormatTypePROTOBUF DataFormatType = "PROTOBUF"
	DataFormatTypeSTRING   DataFormatType = "STRING"
)

func PossibleValuesForDataFormatType() []string {
	return []string{
		string(DataFormatTypeAVRO),
		string(DataFormatTypeBYTES),
		string(DataFormatTypeJSON),
		string(DataFormatTypePROTOBUF),
		string(DataFormatTypeSTRING),
	}
}

func (s *DataFormatType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataFormatType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataFormatType(input string) (*DataFormatType, error) {
	vals := map[string]DataFormatType{
		"avro":     DataFormatTypeAVRO,
		"bytes":    DataFormatTypeBYTES,
		"json":     DataFormatTypeJSON,
		"protobuf": DataFormatTypePROTOBUF,
		"string":   DataFormatTypeSTRING,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataFormatType(input)
	return &out, nil
}

type PartnerConnectorType string

const (
	PartnerConnectorTypeKafkaAzureBlobStorageSink      PartnerConnectorType = "KafkaAzureBlobStorageSink"
	PartnerConnectorTypeKafkaAzureBlobStorageSource    PartnerConnectorType = "KafkaAzureBlobStorageSource"
	PartnerConnectorTypeKafkaAzureCosmosDBSink         PartnerConnectorType = "KafkaAzureCosmosDBSink"
	PartnerConnectorTypeKafkaAzureCosmosDBSource       PartnerConnectorType = "KafkaAzureCosmosDBSource"
	PartnerConnectorTypeKafkaAzureSynapseAnalyticsSink PartnerConnectorType = "KafkaAzureSynapseAnalyticsSink"
)

func PossibleValuesForPartnerConnectorType() []string {
	return []string{
		string(PartnerConnectorTypeKafkaAzureBlobStorageSink),
		string(PartnerConnectorTypeKafkaAzureBlobStorageSource),
		string(PartnerConnectorTypeKafkaAzureCosmosDBSink),
		string(PartnerConnectorTypeKafkaAzureCosmosDBSource),
		string(PartnerConnectorTypeKafkaAzureSynapseAnalyticsSink),
	}
}

func (s *PartnerConnectorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePartnerConnectorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePartnerConnectorType(input string) (*PartnerConnectorType, error) {
	vals := map[string]PartnerConnectorType{
		"kafkaazureblobstoragesink":      PartnerConnectorTypeKafkaAzureBlobStorageSink,
		"kafkaazureblobstoragesource":    PartnerConnectorTypeKafkaAzureBlobStorageSource,
		"kafkaazurecosmosdbsink":         PartnerConnectorTypeKafkaAzureCosmosDBSink,
		"kafkaazurecosmosdbsource":       PartnerConnectorTypeKafkaAzureCosmosDBSource,
		"kafkaazuresynapseanalyticssink": PartnerConnectorTypeKafkaAzureSynapseAnalyticsSink,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerConnectorType(input)
	return &out, nil
}
