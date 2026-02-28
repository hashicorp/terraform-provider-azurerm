package dataflowendpoint

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerProtocolType string

const (
	BrokerProtocolTypeMqtt       BrokerProtocolType = "Mqtt"
	BrokerProtocolTypeWebSockets BrokerProtocolType = "WebSockets"
)

func PossibleValuesForBrokerProtocolType() []string {
	return []string{
		string(BrokerProtocolTypeMqtt),
		string(BrokerProtocolTypeWebSockets),
	}
}

func (s *BrokerProtocolType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBrokerProtocolType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBrokerProtocolType(input string) (*BrokerProtocolType, error) {
	vals := map[string]BrokerProtocolType{
		"mqtt":       BrokerProtocolTypeMqtt,
		"websockets": BrokerProtocolTypeWebSockets,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BrokerProtocolType(input)
	return &out, nil
}

type CloudEventAttributeType string

const (
	CloudEventAttributeTypeCreateOrRemap CloudEventAttributeType = "CreateOrRemap"
	CloudEventAttributeTypePropagate     CloudEventAttributeType = "Propagate"
)

func PossibleValuesForCloudEventAttributeType() []string {
	return []string{
		string(CloudEventAttributeTypeCreateOrRemap),
		string(CloudEventAttributeTypePropagate),
	}
}

func (s *CloudEventAttributeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCloudEventAttributeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCloudEventAttributeType(input string) (*CloudEventAttributeType, error) {
	vals := map[string]CloudEventAttributeType{
		"createorremap": CloudEventAttributeTypeCreateOrRemap,
		"propagate":     CloudEventAttributeTypePropagate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CloudEventAttributeType(input)
	return &out, nil
}

type DataLakeStorageAuthMethod string

const (
	DataLakeStorageAuthMethodAccessToken                   DataLakeStorageAuthMethod = "AccessToken"
	DataLakeStorageAuthMethodSystemAssignedManagedIdentity DataLakeStorageAuthMethod = "SystemAssignedManagedIdentity"
	DataLakeStorageAuthMethodUserAssignedManagedIdentity   DataLakeStorageAuthMethod = "UserAssignedManagedIdentity"
)

func PossibleValuesForDataLakeStorageAuthMethod() []string {
	return []string{
		string(DataLakeStorageAuthMethodAccessToken),
		string(DataLakeStorageAuthMethodSystemAssignedManagedIdentity),
		string(DataLakeStorageAuthMethodUserAssignedManagedIdentity),
	}
}

func (s *DataLakeStorageAuthMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataLakeStorageAuthMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataLakeStorageAuthMethod(input string) (*DataLakeStorageAuthMethod, error) {
	vals := map[string]DataLakeStorageAuthMethod{
		"accesstoken":                   DataLakeStorageAuthMethodAccessToken,
		"systemassignedmanagedidentity": DataLakeStorageAuthMethodSystemAssignedManagedIdentity,
		"userassignedmanagedidentity":   DataLakeStorageAuthMethodUserAssignedManagedIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataLakeStorageAuthMethod(input)
	return &out, nil
}

type DataflowEndpointAuthenticationSaslType string

const (
	DataflowEndpointAuthenticationSaslTypePlain              DataflowEndpointAuthenticationSaslType = "Plain"
	DataflowEndpointAuthenticationSaslTypeScramShaFiveOneTwo DataflowEndpointAuthenticationSaslType = "ScramSha512"
	DataflowEndpointAuthenticationSaslTypeScramShaTwoFiveSix DataflowEndpointAuthenticationSaslType = "ScramSha256"
)

func PossibleValuesForDataflowEndpointAuthenticationSaslType() []string {
	return []string{
		string(DataflowEndpointAuthenticationSaslTypePlain),
		string(DataflowEndpointAuthenticationSaslTypeScramShaFiveOneTwo),
		string(DataflowEndpointAuthenticationSaslTypeScramShaTwoFiveSix),
	}
}

func (s *DataflowEndpointAuthenticationSaslType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataflowEndpointAuthenticationSaslType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataflowEndpointAuthenticationSaslType(input string) (*DataflowEndpointAuthenticationSaslType, error) {
	vals := map[string]DataflowEndpointAuthenticationSaslType{
		"plain":       DataflowEndpointAuthenticationSaslTypePlain,
		"scramsha512": DataflowEndpointAuthenticationSaslTypeScramShaFiveOneTwo,
		"scramsha256": DataflowEndpointAuthenticationSaslTypeScramShaTwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataflowEndpointAuthenticationSaslType(input)
	return &out, nil
}

type DataflowEndpointFabricPathType string

const (
	DataflowEndpointFabricPathTypeFiles  DataflowEndpointFabricPathType = "Files"
	DataflowEndpointFabricPathTypeTables DataflowEndpointFabricPathType = "Tables"
)

func PossibleValuesForDataflowEndpointFabricPathType() []string {
	return []string{
		string(DataflowEndpointFabricPathTypeFiles),
		string(DataflowEndpointFabricPathTypeTables),
	}
}

func (s *DataflowEndpointFabricPathType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataflowEndpointFabricPathType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataflowEndpointFabricPathType(input string) (*DataflowEndpointFabricPathType, error) {
	vals := map[string]DataflowEndpointFabricPathType{
		"files":  DataflowEndpointFabricPathTypeFiles,
		"tables": DataflowEndpointFabricPathTypeTables,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataflowEndpointFabricPathType(input)
	return &out, nil
}

type DataflowEndpointKafkaAcks string

const (
	DataflowEndpointKafkaAcksAll  DataflowEndpointKafkaAcks = "All"
	DataflowEndpointKafkaAcksOne  DataflowEndpointKafkaAcks = "One"
	DataflowEndpointKafkaAcksZero DataflowEndpointKafkaAcks = "Zero"
)

func PossibleValuesForDataflowEndpointKafkaAcks() []string {
	return []string{
		string(DataflowEndpointKafkaAcksAll),
		string(DataflowEndpointKafkaAcksOne),
		string(DataflowEndpointKafkaAcksZero),
	}
}

func (s *DataflowEndpointKafkaAcks) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataflowEndpointKafkaAcks(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataflowEndpointKafkaAcks(input string) (*DataflowEndpointKafkaAcks, error) {
	vals := map[string]DataflowEndpointKafkaAcks{
		"all":  DataflowEndpointKafkaAcksAll,
		"one":  DataflowEndpointKafkaAcksOne,
		"zero": DataflowEndpointKafkaAcksZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataflowEndpointKafkaAcks(input)
	return &out, nil
}

type DataflowEndpointKafkaCompression string

const (
	DataflowEndpointKafkaCompressionGzip   DataflowEndpointKafkaCompression = "Gzip"
	DataflowEndpointKafkaCompressionLzFour DataflowEndpointKafkaCompression = "Lz4"
	DataflowEndpointKafkaCompressionNone   DataflowEndpointKafkaCompression = "None"
	DataflowEndpointKafkaCompressionSnappy DataflowEndpointKafkaCompression = "Snappy"
)

func PossibleValuesForDataflowEndpointKafkaCompression() []string {
	return []string{
		string(DataflowEndpointKafkaCompressionGzip),
		string(DataflowEndpointKafkaCompressionLzFour),
		string(DataflowEndpointKafkaCompressionNone),
		string(DataflowEndpointKafkaCompressionSnappy),
	}
}

func (s *DataflowEndpointKafkaCompression) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataflowEndpointKafkaCompression(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataflowEndpointKafkaCompression(input string) (*DataflowEndpointKafkaCompression, error) {
	vals := map[string]DataflowEndpointKafkaCompression{
		"gzip":   DataflowEndpointKafkaCompressionGzip,
		"lz4":    DataflowEndpointKafkaCompressionLzFour,
		"none":   DataflowEndpointKafkaCompressionNone,
		"snappy": DataflowEndpointKafkaCompressionSnappy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataflowEndpointKafkaCompression(input)
	return &out, nil
}

type DataflowEndpointKafkaPartitionStrategy string

const (
	DataflowEndpointKafkaPartitionStrategyDefault  DataflowEndpointKafkaPartitionStrategy = "Default"
	DataflowEndpointKafkaPartitionStrategyProperty DataflowEndpointKafkaPartitionStrategy = "Property"
	DataflowEndpointKafkaPartitionStrategyStatic   DataflowEndpointKafkaPartitionStrategy = "Static"
	DataflowEndpointKafkaPartitionStrategyTopic    DataflowEndpointKafkaPartitionStrategy = "Topic"
)

func PossibleValuesForDataflowEndpointKafkaPartitionStrategy() []string {
	return []string{
		string(DataflowEndpointKafkaPartitionStrategyDefault),
		string(DataflowEndpointKafkaPartitionStrategyProperty),
		string(DataflowEndpointKafkaPartitionStrategyStatic),
		string(DataflowEndpointKafkaPartitionStrategyTopic),
	}
}

func (s *DataflowEndpointKafkaPartitionStrategy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataflowEndpointKafkaPartitionStrategy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataflowEndpointKafkaPartitionStrategy(input string) (*DataflowEndpointKafkaPartitionStrategy, error) {
	vals := map[string]DataflowEndpointKafkaPartitionStrategy{
		"default":  DataflowEndpointKafkaPartitionStrategyDefault,
		"property": DataflowEndpointKafkaPartitionStrategyProperty,
		"static":   DataflowEndpointKafkaPartitionStrategyStatic,
		"topic":    DataflowEndpointKafkaPartitionStrategyTopic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataflowEndpointKafkaPartitionStrategy(input)
	return &out, nil
}

type EndpointType string

const (
	EndpointTypeDataExplorer    EndpointType = "DataExplorer"
	EndpointTypeDataLakeStorage EndpointType = "DataLakeStorage"
	EndpointTypeFabricOneLake   EndpointType = "FabricOneLake"
	EndpointTypeKafka           EndpointType = "Kafka"
	EndpointTypeLocalStorage    EndpointType = "LocalStorage"
	EndpointTypeMqtt            EndpointType = "Mqtt"
)

func PossibleValuesForEndpointType() []string {
	return []string{
		string(EndpointTypeDataExplorer),
		string(EndpointTypeDataLakeStorage),
		string(EndpointTypeFabricOneLake),
		string(EndpointTypeKafka),
		string(EndpointTypeLocalStorage),
		string(EndpointTypeMqtt),
	}
}

func (s *EndpointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEndpointType(input string) (*EndpointType, error) {
	vals := map[string]EndpointType{
		"dataexplorer":    EndpointTypeDataExplorer,
		"datalakestorage": EndpointTypeDataLakeStorage,
		"fabriconelake":   EndpointTypeFabricOneLake,
		"kafka":           EndpointTypeKafka,
		"localstorage":    EndpointTypeLocalStorage,
		"mqtt":            EndpointTypeMqtt,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointType(input)
	return &out, nil
}

type ExtendedLocationType string

const (
	ExtendedLocationTypeCustomLocation ExtendedLocationType = "CustomLocation"
)

func PossibleValuesForExtendedLocationType() []string {
	return []string{
		string(ExtendedLocationTypeCustomLocation),
	}
}

func (s *ExtendedLocationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtendedLocationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtendedLocationType(input string) (*ExtendedLocationType, error) {
	vals := map[string]ExtendedLocationType{
		"customlocation": ExtendedLocationTypeCustomLocation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationType(input)
	return &out, nil
}

type KafkaAuthMethod string

const (
	KafkaAuthMethodAnonymous                     KafkaAuthMethod = "Anonymous"
	KafkaAuthMethodSasl                          KafkaAuthMethod = "Sasl"
	KafkaAuthMethodSystemAssignedManagedIdentity KafkaAuthMethod = "SystemAssignedManagedIdentity"
	KafkaAuthMethodUserAssignedManagedIdentity   KafkaAuthMethod = "UserAssignedManagedIdentity"
	KafkaAuthMethodXFiveZeroNineCertificate      KafkaAuthMethod = "X509Certificate"
)

func PossibleValuesForKafkaAuthMethod() []string {
	return []string{
		string(KafkaAuthMethodAnonymous),
		string(KafkaAuthMethodSasl),
		string(KafkaAuthMethodSystemAssignedManagedIdentity),
		string(KafkaAuthMethodUserAssignedManagedIdentity),
		string(KafkaAuthMethodXFiveZeroNineCertificate),
	}
}

func (s *KafkaAuthMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKafkaAuthMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKafkaAuthMethod(input string) (*KafkaAuthMethod, error) {
	vals := map[string]KafkaAuthMethod{
		"anonymous":                     KafkaAuthMethodAnonymous,
		"sasl":                          KafkaAuthMethodSasl,
		"systemassignedmanagedidentity": KafkaAuthMethodSystemAssignedManagedIdentity,
		"userassignedmanagedidentity":   KafkaAuthMethodUserAssignedManagedIdentity,
		"x509certificate":               KafkaAuthMethodXFiveZeroNineCertificate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KafkaAuthMethod(input)
	return &out, nil
}

type ManagedIdentityMethod string

const (
	ManagedIdentityMethodSystemAssignedManagedIdentity ManagedIdentityMethod = "SystemAssignedManagedIdentity"
	ManagedIdentityMethodUserAssignedManagedIdentity   ManagedIdentityMethod = "UserAssignedManagedIdentity"
)

func PossibleValuesForManagedIdentityMethod() []string {
	return []string{
		string(ManagedIdentityMethodSystemAssignedManagedIdentity),
		string(ManagedIdentityMethodUserAssignedManagedIdentity),
	}
}

func (s *ManagedIdentityMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedIdentityMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedIdentityMethod(input string) (*ManagedIdentityMethod, error) {
	vals := map[string]ManagedIdentityMethod{
		"systemassignedmanagedidentity": ManagedIdentityMethodSystemAssignedManagedIdentity,
		"userassignedmanagedidentity":   ManagedIdentityMethodUserAssignedManagedIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedIdentityMethod(input)
	return &out, nil
}

type MqttAuthMethod string

const (
	MqttAuthMethodAnonymous                     MqttAuthMethod = "Anonymous"
	MqttAuthMethodServiceAccountToken           MqttAuthMethod = "ServiceAccountToken"
	MqttAuthMethodSystemAssignedManagedIdentity MqttAuthMethod = "SystemAssignedManagedIdentity"
	MqttAuthMethodUserAssignedManagedIdentity   MqttAuthMethod = "UserAssignedManagedIdentity"
	MqttAuthMethodXFiveZeroNineCertificate      MqttAuthMethod = "X509Certificate"
)

func PossibleValuesForMqttAuthMethod() []string {
	return []string{
		string(MqttAuthMethodAnonymous),
		string(MqttAuthMethodServiceAccountToken),
		string(MqttAuthMethodSystemAssignedManagedIdentity),
		string(MqttAuthMethodUserAssignedManagedIdentity),
		string(MqttAuthMethodXFiveZeroNineCertificate),
	}
}

func (s *MqttAuthMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMqttAuthMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMqttAuthMethod(input string) (*MqttAuthMethod, error) {
	vals := map[string]MqttAuthMethod{
		"anonymous":                     MqttAuthMethodAnonymous,
		"serviceaccounttoken":           MqttAuthMethodServiceAccountToken,
		"systemassignedmanagedidentity": MqttAuthMethodSystemAssignedManagedIdentity,
		"userassignedmanagedidentity":   MqttAuthMethodUserAssignedManagedIdentity,
		"x509certificate":               MqttAuthMethodXFiveZeroNineCertificate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MqttAuthMethod(input)
	return &out, nil
}

type MqttRetainType string

const (
	MqttRetainTypeKeep  MqttRetainType = "Keep"
	MqttRetainTypeNever MqttRetainType = "Never"
)

func PossibleValuesForMqttRetainType() []string {
	return []string{
		string(MqttRetainTypeKeep),
		string(MqttRetainTypeNever),
	}
}

func (s *MqttRetainType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMqttRetainType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMqttRetainType(input string) (*MqttRetainType, error) {
	vals := map[string]MqttRetainType{
		"keep":  MqttRetainTypeKeep,
		"never": MqttRetainTypeNever,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MqttRetainType(input)
	return &out, nil
}

type OperationalMode string

const (
	OperationalModeDisabled OperationalMode = "Disabled"
	OperationalModeEnabled  OperationalMode = "Enabled"
)

func PossibleValuesForOperationalMode() []string {
	return []string{
		string(OperationalModeDisabled),
		string(OperationalModeEnabled),
	}
}

func (s *OperationalMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationalMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationalMode(input string) (*OperationalMode, error) {
	vals := map[string]OperationalMode{
		"disabled": OperationalModeDisabled,
		"enabled":  OperationalModeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationalMode(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
