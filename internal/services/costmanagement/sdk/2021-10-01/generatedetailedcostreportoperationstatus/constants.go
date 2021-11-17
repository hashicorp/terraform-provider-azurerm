package generatedetailedcostreportoperationstatus

import "strings"

type OperationStatusType string

const (
	OperationStatusTypeCompleted       OperationStatusType = "Completed"
	OperationStatusTypeFailed          OperationStatusType = "Failed"
	OperationStatusTypeInProgress      OperationStatusType = "InProgress"
	OperationStatusTypeNoDataFound     OperationStatusType = "NoDataFound"
	OperationStatusTypeQueued          OperationStatusType = "Queued"
	OperationStatusTypeReadyToDownload OperationStatusType = "ReadyToDownload"
	OperationStatusTypeTimedOut        OperationStatusType = "TimedOut"
)

func PossibleValuesForOperationStatusType() []string {
	return []string{
		string(OperationStatusTypeCompleted),
		string(OperationStatusTypeFailed),
		string(OperationStatusTypeInProgress),
		string(OperationStatusTypeNoDataFound),
		string(OperationStatusTypeQueued),
		string(OperationStatusTypeReadyToDownload),
		string(OperationStatusTypeTimedOut),
	}
}

func parseOperationStatusType(input string) (*OperationStatusType, error) {
	vals := map[string]OperationStatusType{
		"completed":       OperationStatusTypeCompleted,
		"failed":          OperationStatusTypeFailed,
		"inprogress":      OperationStatusTypeInProgress,
		"nodatafound":     OperationStatusTypeNoDataFound,
		"queued":          OperationStatusTypeQueued,
		"readytodownload": OperationStatusTypeReadyToDownload,
		"timedout":        OperationStatusTypeTimedOut,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationStatusType(input)
	return &out, nil
}
