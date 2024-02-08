package dataconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobStorageEventType string

const (
	BlobStorageEventTypeMicrosoftPointStoragePointBlobCreated BlobStorageEventType = "Microsoft.Storage.BlobCreated"
	BlobStorageEventTypeMicrosoftPointStoragePointBlobRenamed BlobStorageEventType = "Microsoft.Storage.BlobRenamed"
)

func PossibleValuesForBlobStorageEventType() []string {
	return []string{
		string(BlobStorageEventTypeMicrosoftPointStoragePointBlobCreated),
		string(BlobStorageEventTypeMicrosoftPointStoragePointBlobRenamed),
	}
}

func (s *BlobStorageEventType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlobStorageEventType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlobStorageEventType(input string) (*BlobStorageEventType, error) {
	vals := map[string]BlobStorageEventType{
		"microsoft.storage.blobcreated": BlobStorageEventTypeMicrosoftPointStoragePointBlobCreated,
		"microsoft.storage.blobrenamed": BlobStorageEventTypeMicrosoftPointStoragePointBlobRenamed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlobStorageEventType(input)
	return &out, nil
}

type Compression string

const (
	CompressionGZip Compression = "GZip"
	CompressionNone Compression = "None"
)

func PossibleValuesForCompression() []string {
	return []string{
		string(CompressionGZip),
		string(CompressionNone),
	}
}

func (s *Compression) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCompression(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCompression(input string) (*Compression, error) {
	vals := map[string]Compression{
		"gzip": CompressionGZip,
		"none": CompressionNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Compression(input)
	return &out, nil
}

type DataConnectionKind string

const (
	DataConnectionKindCosmosDb  DataConnectionKind = "CosmosDb"
	DataConnectionKindEventGrid DataConnectionKind = "EventGrid"
	DataConnectionKindEventHub  DataConnectionKind = "EventHub"
	DataConnectionKindIotHub    DataConnectionKind = "IotHub"
)

func PossibleValuesForDataConnectionKind() []string {
	return []string{
		string(DataConnectionKindCosmosDb),
		string(DataConnectionKindEventGrid),
		string(DataConnectionKindEventHub),
		string(DataConnectionKindIotHub),
	}
}

func (s *DataConnectionKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataConnectionKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataConnectionKind(input string) (*DataConnectionKind, error) {
	vals := map[string]DataConnectionKind{
		"cosmosdb":  DataConnectionKindCosmosDb,
		"eventgrid": DataConnectionKindEventGrid,
		"eventhub":  DataConnectionKindEventHub,
		"iothub":    DataConnectionKindIotHub,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataConnectionKind(input)
	return &out, nil
}

type DataConnectionType string

const (
	DataConnectionTypeMicrosoftPointKustoClustersDatabasesDataConnections DataConnectionType = "Microsoft.Kusto/clusters/databases/dataConnections"
)

func PossibleValuesForDataConnectionType() []string {
	return []string{
		string(DataConnectionTypeMicrosoftPointKustoClustersDatabasesDataConnections),
	}
}

func (s *DataConnectionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataConnectionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataConnectionType(input string) (*DataConnectionType, error) {
	vals := map[string]DataConnectionType{
		"microsoft.kusto/clusters/databases/dataconnections": DataConnectionTypeMicrosoftPointKustoClustersDatabasesDataConnections,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataConnectionType(input)
	return &out, nil
}

type DatabaseRouting string

const (
	DatabaseRoutingMulti  DatabaseRouting = "Multi"
	DatabaseRoutingSingle DatabaseRouting = "Single"
)

func PossibleValuesForDatabaseRouting() []string {
	return []string{
		string(DatabaseRoutingMulti),
		string(DatabaseRoutingSingle),
	}
}

func (s *DatabaseRouting) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseRouting(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseRouting(input string) (*DatabaseRouting, error) {
	vals := map[string]DatabaseRouting{
		"multi":  DatabaseRoutingMulti,
		"single": DatabaseRoutingSingle,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseRouting(input)
	return &out, nil
}

type EventGridDataFormat string

const (
	EventGridDataFormatAPACHEAVRO     EventGridDataFormat = "APACHEAVRO"
	EventGridDataFormatAVRO           EventGridDataFormat = "AVRO"
	EventGridDataFormatCSV            EventGridDataFormat = "CSV"
	EventGridDataFormatJSON           EventGridDataFormat = "JSON"
	EventGridDataFormatMULTIJSON      EventGridDataFormat = "MULTIJSON"
	EventGridDataFormatORC            EventGridDataFormat = "ORC"
	EventGridDataFormatPARQUET        EventGridDataFormat = "PARQUET"
	EventGridDataFormatPSV            EventGridDataFormat = "PSV"
	EventGridDataFormatRAW            EventGridDataFormat = "RAW"
	EventGridDataFormatSCSV           EventGridDataFormat = "SCSV"
	EventGridDataFormatSINGLEJSON     EventGridDataFormat = "SINGLEJSON"
	EventGridDataFormatSOHSV          EventGridDataFormat = "SOHSV"
	EventGridDataFormatTSV            EventGridDataFormat = "TSV"
	EventGridDataFormatTSVE           EventGridDataFormat = "TSVE"
	EventGridDataFormatTXT            EventGridDataFormat = "TXT"
	EventGridDataFormatWThreeCLOGFILE EventGridDataFormat = "W3CLOGFILE"
)

func PossibleValuesForEventGridDataFormat() []string {
	return []string{
		string(EventGridDataFormatAPACHEAVRO),
		string(EventGridDataFormatAVRO),
		string(EventGridDataFormatCSV),
		string(EventGridDataFormatJSON),
		string(EventGridDataFormatMULTIJSON),
		string(EventGridDataFormatORC),
		string(EventGridDataFormatPARQUET),
		string(EventGridDataFormatPSV),
		string(EventGridDataFormatRAW),
		string(EventGridDataFormatSCSV),
		string(EventGridDataFormatSINGLEJSON),
		string(EventGridDataFormatSOHSV),
		string(EventGridDataFormatTSV),
		string(EventGridDataFormatTSVE),
		string(EventGridDataFormatTXT),
		string(EventGridDataFormatWThreeCLOGFILE),
	}
}

func (s *EventGridDataFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventGridDataFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEventGridDataFormat(input string) (*EventGridDataFormat, error) {
	vals := map[string]EventGridDataFormat{
		"apacheavro": EventGridDataFormatAPACHEAVRO,
		"avro":       EventGridDataFormatAVRO,
		"csv":        EventGridDataFormatCSV,
		"json":       EventGridDataFormatJSON,
		"multijson":  EventGridDataFormatMULTIJSON,
		"orc":        EventGridDataFormatORC,
		"parquet":    EventGridDataFormatPARQUET,
		"psv":        EventGridDataFormatPSV,
		"raw":        EventGridDataFormatRAW,
		"scsv":       EventGridDataFormatSCSV,
		"singlejson": EventGridDataFormatSINGLEJSON,
		"sohsv":      EventGridDataFormatSOHSV,
		"tsv":        EventGridDataFormatTSV,
		"tsve":       EventGridDataFormatTSVE,
		"txt":        EventGridDataFormatTXT,
		"w3clogfile": EventGridDataFormatWThreeCLOGFILE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventGridDataFormat(input)
	return &out, nil
}

type EventHubDataFormat string

const (
	EventHubDataFormatAPACHEAVRO     EventHubDataFormat = "APACHEAVRO"
	EventHubDataFormatAVRO           EventHubDataFormat = "AVRO"
	EventHubDataFormatCSV            EventHubDataFormat = "CSV"
	EventHubDataFormatJSON           EventHubDataFormat = "JSON"
	EventHubDataFormatMULTIJSON      EventHubDataFormat = "MULTIJSON"
	EventHubDataFormatORC            EventHubDataFormat = "ORC"
	EventHubDataFormatPARQUET        EventHubDataFormat = "PARQUET"
	EventHubDataFormatPSV            EventHubDataFormat = "PSV"
	EventHubDataFormatRAW            EventHubDataFormat = "RAW"
	EventHubDataFormatSCSV           EventHubDataFormat = "SCSV"
	EventHubDataFormatSINGLEJSON     EventHubDataFormat = "SINGLEJSON"
	EventHubDataFormatSOHSV          EventHubDataFormat = "SOHSV"
	EventHubDataFormatTSV            EventHubDataFormat = "TSV"
	EventHubDataFormatTSVE           EventHubDataFormat = "TSVE"
	EventHubDataFormatTXT            EventHubDataFormat = "TXT"
	EventHubDataFormatWThreeCLOGFILE EventHubDataFormat = "W3CLOGFILE"
)

func PossibleValuesForEventHubDataFormat() []string {
	return []string{
		string(EventHubDataFormatAPACHEAVRO),
		string(EventHubDataFormatAVRO),
		string(EventHubDataFormatCSV),
		string(EventHubDataFormatJSON),
		string(EventHubDataFormatMULTIJSON),
		string(EventHubDataFormatORC),
		string(EventHubDataFormatPARQUET),
		string(EventHubDataFormatPSV),
		string(EventHubDataFormatRAW),
		string(EventHubDataFormatSCSV),
		string(EventHubDataFormatSINGLEJSON),
		string(EventHubDataFormatSOHSV),
		string(EventHubDataFormatTSV),
		string(EventHubDataFormatTSVE),
		string(EventHubDataFormatTXT),
		string(EventHubDataFormatWThreeCLOGFILE),
	}
}

func (s *EventHubDataFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEventHubDataFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEventHubDataFormat(input string) (*EventHubDataFormat, error) {
	vals := map[string]EventHubDataFormat{
		"apacheavro": EventHubDataFormatAPACHEAVRO,
		"avro":       EventHubDataFormatAVRO,
		"csv":        EventHubDataFormatCSV,
		"json":       EventHubDataFormatJSON,
		"multijson":  EventHubDataFormatMULTIJSON,
		"orc":        EventHubDataFormatORC,
		"parquet":    EventHubDataFormatPARQUET,
		"psv":        EventHubDataFormatPSV,
		"raw":        EventHubDataFormatRAW,
		"scsv":       EventHubDataFormatSCSV,
		"singlejson": EventHubDataFormatSINGLEJSON,
		"sohsv":      EventHubDataFormatSOHSV,
		"tsv":        EventHubDataFormatTSV,
		"tsve":       EventHubDataFormatTSVE,
		"txt":        EventHubDataFormatTXT,
		"w3clogfile": EventHubDataFormatWThreeCLOGFILE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventHubDataFormat(input)
	return &out, nil
}

type IotHubDataFormat string

const (
	IotHubDataFormatAPACHEAVRO     IotHubDataFormat = "APACHEAVRO"
	IotHubDataFormatAVRO           IotHubDataFormat = "AVRO"
	IotHubDataFormatCSV            IotHubDataFormat = "CSV"
	IotHubDataFormatJSON           IotHubDataFormat = "JSON"
	IotHubDataFormatMULTIJSON      IotHubDataFormat = "MULTIJSON"
	IotHubDataFormatORC            IotHubDataFormat = "ORC"
	IotHubDataFormatPARQUET        IotHubDataFormat = "PARQUET"
	IotHubDataFormatPSV            IotHubDataFormat = "PSV"
	IotHubDataFormatRAW            IotHubDataFormat = "RAW"
	IotHubDataFormatSCSV           IotHubDataFormat = "SCSV"
	IotHubDataFormatSINGLEJSON     IotHubDataFormat = "SINGLEJSON"
	IotHubDataFormatSOHSV          IotHubDataFormat = "SOHSV"
	IotHubDataFormatTSV            IotHubDataFormat = "TSV"
	IotHubDataFormatTSVE           IotHubDataFormat = "TSVE"
	IotHubDataFormatTXT            IotHubDataFormat = "TXT"
	IotHubDataFormatWThreeCLOGFILE IotHubDataFormat = "W3CLOGFILE"
)

func PossibleValuesForIotHubDataFormat() []string {
	return []string{
		string(IotHubDataFormatAPACHEAVRO),
		string(IotHubDataFormatAVRO),
		string(IotHubDataFormatCSV),
		string(IotHubDataFormatJSON),
		string(IotHubDataFormatMULTIJSON),
		string(IotHubDataFormatORC),
		string(IotHubDataFormatPARQUET),
		string(IotHubDataFormatPSV),
		string(IotHubDataFormatRAW),
		string(IotHubDataFormatSCSV),
		string(IotHubDataFormatSINGLEJSON),
		string(IotHubDataFormatSOHSV),
		string(IotHubDataFormatTSV),
		string(IotHubDataFormatTSVE),
		string(IotHubDataFormatTXT),
		string(IotHubDataFormatWThreeCLOGFILE),
	}
}

func (s *IotHubDataFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIotHubDataFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIotHubDataFormat(input string) (*IotHubDataFormat, error) {
	vals := map[string]IotHubDataFormat{
		"apacheavro": IotHubDataFormatAPACHEAVRO,
		"avro":       IotHubDataFormatAVRO,
		"csv":        IotHubDataFormatCSV,
		"json":       IotHubDataFormatJSON,
		"multijson":  IotHubDataFormatMULTIJSON,
		"orc":        IotHubDataFormatORC,
		"parquet":    IotHubDataFormatPARQUET,
		"psv":        IotHubDataFormatPSV,
		"raw":        IotHubDataFormatRAW,
		"scsv":       IotHubDataFormatSCSV,
		"singlejson": IotHubDataFormatSINGLEJSON,
		"sohsv":      IotHubDataFormatSOHSV,
		"tsv":        IotHubDataFormatTSV,
		"tsve":       IotHubDataFormatTSVE,
		"txt":        IotHubDataFormatTXT,
		"w3clogfile": IotHubDataFormatWThreeCLOGFILE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IotHubDataFormat(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateMoving    ProvisioningState = "Moving"
	ProvisioningStateRunning   ProvisioningState = "Running"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateRunning),
		string(ProvisioningStateSucceeded),
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
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"moving":    ProvisioningStateMoving,
		"running":   ProvisioningStateRunning,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type Reason string

const (
	ReasonAlreadyExists Reason = "AlreadyExists"
	ReasonInvalid       Reason = "Invalid"
)

func PossibleValuesForReason() []string {
	return []string{
		string(ReasonAlreadyExists),
		string(ReasonInvalid),
	}
}

func (s *Reason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReason(input string) (*Reason, error) {
	vals := map[string]Reason{
		"alreadyexists": ReasonAlreadyExists,
		"invalid":       ReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Reason(input)
	return &out, nil
}
