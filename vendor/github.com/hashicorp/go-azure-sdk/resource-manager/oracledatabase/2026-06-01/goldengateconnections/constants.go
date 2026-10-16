package goldengateconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourceProvisioningState string

const (
	AzureResourceProvisioningStateCanceled     AzureResourceProvisioningState = "Canceled"
	AzureResourceProvisioningStateFailed       AzureResourceProvisioningState = "Failed"
	AzureResourceProvisioningStateProvisioning AzureResourceProvisioningState = "Provisioning"
	AzureResourceProvisioningStateSucceeded    AzureResourceProvisioningState = "Succeeded"
)

func PossibleValuesForAzureResourceProvisioningState() []string {
	return []string{
		string(AzureResourceProvisioningStateCanceled),
		string(AzureResourceProvisioningStateFailed),
		string(AzureResourceProvisioningStateProvisioning),
		string(AzureResourceProvisioningStateSucceeded),
	}
}

func (s *AzureResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureResourceProvisioningState(input string) (*AzureResourceProvisioningState, error) {
	vals := map[string]AzureResourceProvisioningState{
		"canceled":     AzureResourceProvisioningStateCanceled,
		"failed":       AzureResourceProvisioningStateFailed,
		"provisioning": AzureResourceProvisioningStateProvisioning,
		"succeeded":    AzureResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceProvisioningState(input)
	return &out, nil
}

type ConnectionLifecycleState string

const (
	ConnectionLifecycleStateACTIVE   ConnectionLifecycleState = "ACTIVE"
	ConnectionLifecycleStateCREATING ConnectionLifecycleState = "CREATING"
	ConnectionLifecycleStateDELETED  ConnectionLifecycleState = "DELETED"
	ConnectionLifecycleStateDELETING ConnectionLifecycleState = "DELETING"
	ConnectionLifecycleStateFAILED   ConnectionLifecycleState = "FAILED"
	ConnectionLifecycleStateUPDATING ConnectionLifecycleState = "UPDATING"
)

func PossibleValuesForConnectionLifecycleState() []string {
	return []string{
		string(ConnectionLifecycleStateACTIVE),
		string(ConnectionLifecycleStateCREATING),
		string(ConnectionLifecycleStateDELETED),
		string(ConnectionLifecycleStateDELETING),
		string(ConnectionLifecycleStateFAILED),
		string(ConnectionLifecycleStateUPDATING),
	}
}

func (s *ConnectionLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionLifecycleState(input string) (*ConnectionLifecycleState, error) {
	vals := map[string]ConnectionLifecycleState{
		"active":   ConnectionLifecycleStateACTIVE,
		"creating": ConnectionLifecycleStateCREATING,
		"deleted":  ConnectionLifecycleStateDELETED,
		"deleting": ConnectionLifecycleStateDELETING,
		"failed":   ConnectionLifecycleStateFAILED,
		"updating": ConnectionLifecycleStateUPDATING,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionLifecycleState(input)
	return &out, nil
}

type ConnectionType string

const (
	ConnectionTypeAMAZONKINESIS         ConnectionType = "AMAZON_KINESIS"
	ConnectionTypeAMAZONREDSHIFT        ConnectionType = "AMAZON_REDSHIFT"
	ConnectionTypeAMAZONSThree          ConnectionType = "AMAZON_S3"
	ConnectionTypeAZUREDATALAKESTORAGE  ConnectionType = "AZURE_DATA_LAKE_STORAGE"
	ConnectionTypeAZURESYNAPSEANALYTICS ConnectionType = "AZURE_SYNAPSE_ANALYTICS"
	ConnectionTypeDATABRICKS            ConnectionType = "DATABRICKS"
	ConnectionTypeDBTwo                 ConnectionType = "DB2"
	ConnectionTypeELASTICSEARCH         ConnectionType = "ELASTICSEARCH"
	ConnectionTypeGENERIC               ConnectionType = "GENERIC"
	ConnectionTypeGOLDENGATE            ConnectionType = "GOLDENGATE"
	ConnectionTypeGOOGLEBIGQUERY        ConnectionType = "GOOGLE_BIGQUERY"
	ConnectionTypeGOOGLECLOUDSTORAGE    ConnectionType = "GOOGLE_CLOUD_STORAGE"
	ConnectionTypeGOOGLEPUBSUB          ConnectionType = "GOOGLE_PUBSUB"
	ConnectionTypeHDFS                  ConnectionType = "HDFS"
	ConnectionTypeICEBERG               ConnectionType = "ICEBERG"
	ConnectionTypeJAVAMESSAGESERVICE    ConnectionType = "JAVA_MESSAGE_SERVICE"
	ConnectionTypeKAFKA                 ConnectionType = "KAFKA"
	ConnectionTypeKAFKASCHEMAREGISTRY   ConnectionType = "KAFKA_SCHEMA_REGISTRY"
	ConnectionTypeMICROSOFTFABRIC       ConnectionType = "MICROSOFT_FABRIC"
	ConnectionTypeMICROSOFTSQLSERVER    ConnectionType = "MICROSOFT_SQLSERVER"
	ConnectionTypeMONGODB               ConnectionType = "MONGODB"
	ConnectionTypeMYSQL                 ConnectionType = "MYSQL"
	ConnectionTypeOCIOBJECTSTORAGE      ConnectionType = "OCI_OBJECT_STORAGE"
	ConnectionTypeORACLE                ConnectionType = "ORACLE"
	ConnectionTypeORACLENOSQL           ConnectionType = "ORACLE_NOSQL"
	ConnectionTypePOSTGRESQL            ConnectionType = "POSTGRESQL"
	ConnectionTypeREDIS                 ConnectionType = "REDIS"
	ConnectionTypeSNOWFLAKE             ConnectionType = "SNOWFLAKE"
)

func PossibleValuesForConnectionType() []string {
	return []string{
		string(ConnectionTypeAMAZONKINESIS),
		string(ConnectionTypeAMAZONREDSHIFT),
		string(ConnectionTypeAMAZONSThree),
		string(ConnectionTypeAZUREDATALAKESTORAGE),
		string(ConnectionTypeAZURESYNAPSEANALYTICS),
		string(ConnectionTypeDATABRICKS),
		string(ConnectionTypeDBTwo),
		string(ConnectionTypeELASTICSEARCH),
		string(ConnectionTypeGENERIC),
		string(ConnectionTypeGOLDENGATE),
		string(ConnectionTypeGOOGLEBIGQUERY),
		string(ConnectionTypeGOOGLECLOUDSTORAGE),
		string(ConnectionTypeGOOGLEPUBSUB),
		string(ConnectionTypeHDFS),
		string(ConnectionTypeICEBERG),
		string(ConnectionTypeJAVAMESSAGESERVICE),
		string(ConnectionTypeKAFKA),
		string(ConnectionTypeKAFKASCHEMAREGISTRY),
		string(ConnectionTypeMICROSOFTFABRIC),
		string(ConnectionTypeMICROSOFTSQLSERVER),
		string(ConnectionTypeMONGODB),
		string(ConnectionTypeMYSQL),
		string(ConnectionTypeOCIOBJECTSTORAGE),
		string(ConnectionTypeORACLE),
		string(ConnectionTypeORACLENOSQL),
		string(ConnectionTypePOSTGRESQL),
		string(ConnectionTypeREDIS),
		string(ConnectionTypeSNOWFLAKE),
	}
}

func (s *ConnectionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionType(input string) (*ConnectionType, error) {
	vals := map[string]ConnectionType{
		"amazon_kinesis":          ConnectionTypeAMAZONKINESIS,
		"amazon_redshift":         ConnectionTypeAMAZONREDSHIFT,
		"amazon_s3":               ConnectionTypeAMAZONSThree,
		"azure_data_lake_storage": ConnectionTypeAZUREDATALAKESTORAGE,
		"azure_synapse_analytics": ConnectionTypeAZURESYNAPSEANALYTICS,
		"databricks":              ConnectionTypeDATABRICKS,
		"db2":                     ConnectionTypeDBTwo,
		"elasticsearch":           ConnectionTypeELASTICSEARCH,
		"generic":                 ConnectionTypeGENERIC,
		"goldengate":              ConnectionTypeGOLDENGATE,
		"google_bigquery":         ConnectionTypeGOOGLEBIGQUERY,
		"google_cloud_storage":    ConnectionTypeGOOGLECLOUDSTORAGE,
		"google_pubsub":           ConnectionTypeGOOGLEPUBSUB,
		"hdfs":                    ConnectionTypeHDFS,
		"iceberg":                 ConnectionTypeICEBERG,
		"java_message_service":    ConnectionTypeJAVAMESSAGESERVICE,
		"kafka":                   ConnectionTypeKAFKA,
		"kafka_schema_registry":   ConnectionTypeKAFKASCHEMAREGISTRY,
		"microsoft_fabric":        ConnectionTypeMICROSOFTFABRIC,
		"microsoft_sqlserver":     ConnectionTypeMICROSOFTSQLSERVER,
		"mongodb":                 ConnectionTypeMONGODB,
		"mysql":                   ConnectionTypeMYSQL,
		"oci_object_storage":      ConnectionTypeOCIOBJECTSTORAGE,
		"oracle":                  ConnectionTypeORACLE,
		"oracle_nosql":            ConnectionTypeORACLENOSQL,
		"postgresql":              ConnectionTypePOSTGRESQL,
		"redis":                   ConnectionTypeREDIS,
		"snowflake":               ConnectionTypeSNOWFLAKE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionType(input)
	return &out, nil
}

type GoldenGateConnectionAssignmentLifecycleState string

const (
	GoldenGateConnectionAssignmentLifecycleStateACTIVE   GoldenGateConnectionAssignmentLifecycleState = "ACTIVE"
	GoldenGateConnectionAssignmentLifecycleStateCREATING GoldenGateConnectionAssignmentLifecycleState = "CREATING"
	GoldenGateConnectionAssignmentLifecycleStateDELETED  GoldenGateConnectionAssignmentLifecycleState = "DELETED"
	GoldenGateConnectionAssignmentLifecycleStateDELETING GoldenGateConnectionAssignmentLifecycleState = "DELETING"
	GoldenGateConnectionAssignmentLifecycleStateFAILED   GoldenGateConnectionAssignmentLifecycleState = "FAILED"
	GoldenGateConnectionAssignmentLifecycleStateUPDATING GoldenGateConnectionAssignmentLifecycleState = "UPDATING"
)

func PossibleValuesForGoldenGateConnectionAssignmentLifecycleState() []string {
	return []string{
		string(GoldenGateConnectionAssignmentLifecycleStateACTIVE),
		string(GoldenGateConnectionAssignmentLifecycleStateCREATING),
		string(GoldenGateConnectionAssignmentLifecycleStateDELETED),
		string(GoldenGateConnectionAssignmentLifecycleStateDELETING),
		string(GoldenGateConnectionAssignmentLifecycleStateFAILED),
		string(GoldenGateConnectionAssignmentLifecycleStateUPDATING),
	}
}

func (s *GoldenGateConnectionAssignmentLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGoldenGateConnectionAssignmentLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGoldenGateConnectionAssignmentLifecycleState(input string) (*GoldenGateConnectionAssignmentLifecycleState, error) {
	vals := map[string]GoldenGateConnectionAssignmentLifecycleState{
		"active":   GoldenGateConnectionAssignmentLifecycleStateACTIVE,
		"creating": GoldenGateConnectionAssignmentLifecycleStateCREATING,
		"deleted":  GoldenGateConnectionAssignmentLifecycleStateDELETED,
		"deleting": GoldenGateConnectionAssignmentLifecycleStateDELETING,
		"failed":   GoldenGateConnectionAssignmentLifecycleStateFAILED,
		"updating": GoldenGateConnectionAssignmentLifecycleStateUPDATING,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GoldenGateConnectionAssignmentLifecycleState(input)
	return &out, nil
}

type KafkaConnectionTechnologyType string

const (
	KafkaConnectionTechnologyTypeAPACHEKAFKA    KafkaConnectionTechnologyType = "APACHE_KAFKA"
	KafkaConnectionTechnologyTypeAZUREEVENTHUBS KafkaConnectionTechnologyType = "AZURE_EVENT_HUBS"
	KafkaConnectionTechnologyTypeCONFLUENTKAFKA KafkaConnectionTechnologyType = "CONFLUENT_KAFKA"
	KafkaConnectionTechnologyTypeOCISTREAMING   KafkaConnectionTechnologyType = "OCI_STREAMING"
)

func PossibleValuesForKafkaConnectionTechnologyType() []string {
	return []string{
		string(KafkaConnectionTechnologyTypeAPACHEKAFKA),
		string(KafkaConnectionTechnologyTypeAZUREEVENTHUBS),
		string(KafkaConnectionTechnologyTypeCONFLUENTKAFKA),
		string(KafkaConnectionTechnologyTypeOCISTREAMING),
	}
}

func (s *KafkaConnectionTechnologyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKafkaConnectionTechnologyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKafkaConnectionTechnologyType(input string) (*KafkaConnectionTechnologyType, error) {
	vals := map[string]KafkaConnectionTechnologyType{
		"apache_kafka":     KafkaConnectionTechnologyTypeAPACHEKAFKA,
		"azure_event_hubs": KafkaConnectionTechnologyTypeAZUREEVENTHUBS,
		"confluent_kafka":  KafkaConnectionTechnologyTypeCONFLUENTKAFKA,
		"oci_streaming":    KafkaConnectionTechnologyTypeOCISTREAMING,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KafkaConnectionTechnologyType(input)
	return &out, nil
}

type MicrosoftFabricConnectionTechnologyType string

const (
	MicrosoftFabricConnectionTechnologyTypeMICROSOFTFABRICLAKEHOUSE MicrosoftFabricConnectionTechnologyType = "MICROSOFT_FABRIC_LAKEHOUSE"
	MicrosoftFabricConnectionTechnologyTypeMICROSOFTFABRICMIRROR    MicrosoftFabricConnectionTechnologyType = "MICROSOFT_FABRIC_MIRROR"
)

func PossibleValuesForMicrosoftFabricConnectionTechnologyType() []string {
	return []string{
		string(MicrosoftFabricConnectionTechnologyTypeMICROSOFTFABRICLAKEHOUSE),
		string(MicrosoftFabricConnectionTechnologyTypeMICROSOFTFABRICMIRROR),
	}
}

func (s *MicrosoftFabricConnectionTechnologyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMicrosoftFabricConnectionTechnologyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMicrosoftFabricConnectionTechnologyType(input string) (*MicrosoftFabricConnectionTechnologyType, error) {
	vals := map[string]MicrosoftFabricConnectionTechnologyType{
		"microsoft_fabric_lakehouse": MicrosoftFabricConnectionTechnologyTypeMICROSOFTFABRICLAKEHOUSE,
		"microsoft_fabric_mirror":    MicrosoftFabricConnectionTechnologyTypeMICROSOFTFABRICMIRROR,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MicrosoftFabricConnectionTechnologyType(input)
	return &out, nil
}

type OracleConnectionTechnologyType string

const (
	OracleConnectionTechnologyTypeAMAZONRDSORACLE                       OracleConnectionTechnologyType = "AMAZON_RDS_ORACLE"
	OracleConnectionTechnologyTypeOCIAUTONOMOUSDATABASE                 OracleConnectionTechnologyType = "OCI_AUTONOMOUS_DATABASE"
	OracleConnectionTechnologyTypeORACLEAUTONOMOUSDATABASEATAWS         OracleConnectionTechnologyType = "ORACLE_AUTONOMOUS_DATABASE_AT_AWS"
	OracleConnectionTechnologyTypeORACLEAUTONOMOUSDATABASEATAZURE       OracleConnectionTechnologyType = "ORACLE_AUTONOMOUS_DATABASE_AT_AZURE"
	OracleConnectionTechnologyTypeORACLEAUTONOMOUSDATABASEATGOOGLECLOUD OracleConnectionTechnologyType = "ORACLE_AUTONOMOUS_DATABASE_AT_GOOGLE_CLOUD"
	OracleConnectionTechnologyTypeORACLEDATABASE                        OracleConnectionTechnologyType = "ORACLE_DATABASE"
	OracleConnectionTechnologyTypeORACLEEXADATA                         OracleConnectionTechnologyType = "ORACLE_EXADATA"
	OracleConnectionTechnologyTypeORACLEEXADATADATABASEATAWS            OracleConnectionTechnologyType = "ORACLE_EXADATA_DATABASE_AT_AWS"
	OracleConnectionTechnologyTypeORACLEEXADATADATABASEATAZURE          OracleConnectionTechnologyType = "ORACLE_EXADATA_DATABASE_AT_AZURE"
	OracleConnectionTechnologyTypeORACLEEXADATADATABASEATGOOGLECLOUD    OracleConnectionTechnologyType = "ORACLE_EXADATA_DATABASE_AT_GOOGLE_CLOUD"
)

func PossibleValuesForOracleConnectionTechnologyType() []string {
	return []string{
		string(OracleConnectionTechnologyTypeAMAZONRDSORACLE),
		string(OracleConnectionTechnologyTypeOCIAUTONOMOUSDATABASE),
		string(OracleConnectionTechnologyTypeORACLEAUTONOMOUSDATABASEATAWS),
		string(OracleConnectionTechnologyTypeORACLEAUTONOMOUSDATABASEATAZURE),
		string(OracleConnectionTechnologyTypeORACLEAUTONOMOUSDATABASEATGOOGLECLOUD),
		string(OracleConnectionTechnologyTypeORACLEDATABASE),
		string(OracleConnectionTechnologyTypeORACLEEXADATA),
		string(OracleConnectionTechnologyTypeORACLEEXADATADATABASEATAWS),
		string(OracleConnectionTechnologyTypeORACLEEXADATADATABASEATAZURE),
		string(OracleConnectionTechnologyTypeORACLEEXADATADATABASEATGOOGLECLOUD),
	}
}

func (s *OracleConnectionTechnologyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOracleConnectionTechnologyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOracleConnectionTechnologyType(input string) (*OracleConnectionTechnologyType, error) {
	vals := map[string]OracleConnectionTechnologyType{
		"amazon_rds_oracle":                          OracleConnectionTechnologyTypeAMAZONRDSORACLE,
		"oci_autonomous_database":                    OracleConnectionTechnologyTypeOCIAUTONOMOUSDATABASE,
		"oracle_autonomous_database_at_aws":          OracleConnectionTechnologyTypeORACLEAUTONOMOUSDATABASEATAWS,
		"oracle_autonomous_database_at_azure":        OracleConnectionTechnologyTypeORACLEAUTONOMOUSDATABASEATAZURE,
		"oracle_autonomous_database_at_google_cloud": OracleConnectionTechnologyTypeORACLEAUTONOMOUSDATABASEATGOOGLECLOUD,
		"oracle_database":                            OracleConnectionTechnologyTypeORACLEDATABASE,
		"oracle_exadata":                             OracleConnectionTechnologyTypeORACLEEXADATA,
		"oracle_exadata_database_at_aws":             OracleConnectionTechnologyTypeORACLEEXADATADATABASEATAWS,
		"oracle_exadata_database_at_azure":           OracleConnectionTechnologyTypeORACLEEXADATADATABASEATAZURE,
		"oracle_exadata_database_at_google_cloud":    OracleConnectionTechnologyTypeORACLEEXADATADATABASEATGOOGLECLOUD,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OracleConnectionTechnologyType(input)
	return &out, nil
}

type RoutingMethod string

const (
	RoutingMethodDEDICATEDENDPOINT        RoutingMethod = "DEDICATED_ENDPOINT"
	RoutingMethodSHAREDDEPLOYMENTENDPOINT RoutingMethod = "SHARED_DEPLOYMENT_ENDPOINT"
	RoutingMethodSHAREDSERVICEENDPOINT    RoutingMethod = "SHARED_SERVICE_ENDPOINT"
)

func PossibleValuesForRoutingMethod() []string {
	return []string{
		string(RoutingMethodDEDICATEDENDPOINT),
		string(RoutingMethodSHAREDDEPLOYMENTENDPOINT),
		string(RoutingMethodSHAREDSERVICEENDPOINT),
	}
}

func (s *RoutingMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoutingMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoutingMethod(input string) (*RoutingMethod, error) {
	vals := map[string]RoutingMethod{
		"dedicated_endpoint":         RoutingMethodDEDICATEDENDPOINT,
		"shared_deployment_endpoint": RoutingMethodSHAREDDEPLOYMENTENDPOINT,
		"shared_service_endpoint":    RoutingMethodSHAREDSERVICEENDPOINT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingMethod(input)
	return &out, nil
}

type SessionMode string

const (
	SessionModeDIRECT   SessionMode = "DIRECT"
	SessionModeREDIRECT SessionMode = "REDIRECT"
)

func PossibleValuesForSessionMode() []string {
	return []string{
		string(SessionModeDIRECT),
		string(SessionModeREDIRECT),
	}
}

func (s *SessionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSessionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSessionMode(input string) (*SessionMode, error) {
	vals := map[string]SessionMode{
		"direct":   SessionModeDIRECT,
		"redirect": SessionModeREDIRECT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SessionMode(input)
	return &out, nil
}
