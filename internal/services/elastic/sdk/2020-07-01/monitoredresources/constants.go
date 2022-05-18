package monitoredresources

import "strings"

type SendingLogs string

const (
	SendingLogsFalse SendingLogs = "False"
	SendingLogsTrue  SendingLogs = "True"
)

func PossibleValuesForSendingLogs() []string {
	return []string{
		string(SendingLogsFalse),
		string(SendingLogsTrue),
	}
}

func parseSendingLogs(input string) (*SendingLogs, error) {
	vals := map[string]SendingLogs{
		"false": SendingLogsFalse,
		"true":  SendingLogsTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SendingLogs(input)
	return &out, nil
}
