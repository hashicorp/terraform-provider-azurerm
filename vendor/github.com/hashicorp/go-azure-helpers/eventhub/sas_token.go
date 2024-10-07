// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	connStringSharedAccessKeyKey     = "SharedAccessKey"
	connStringSharedAccessKeyNameKey = "SharedAccessKeyName"
)

func ComputeEventHubSASToken(sharedAccessKeyName string,
	sharedAccessKey string,
	eventHubUri string,
	expiry string,
) (string, error) {
	uri := url.QueryEscape(eventHubUri)

	expireTime, err := time.Parse(time.RFC3339, expiry)
	if err != nil {
		return "", err
	}
	expireTimestamp := expireTime.Unix()
	expireStr := strconv.FormatInt(expireTimestamp, 10)

	stringToSign := uri + "\n" + expireStr

	key := []byte(sharedAccessKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	sasToken := "sr=" + uri
	sasToken += "&sig=" + url.QueryEscape(signature)
	sasToken += "&se=" + (expireStr)
	sasToken += "&skn=" + (sharedAccessKeyName)

	return sasToken, nil
}

func ComputeEventHubSASConnectionString(sasToken string) string {
	return fmt.Sprintf("SharedAccessSignature %s", sasToken)
}

func ComputeEventHubSASConnectionUrl(endpoint string, entityPath string) (*string, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	var url string
	if entityPath == "" {
		url = strings.TrimRight(endpoint, "/")
	} else {
		url = endpoint + entityPath
	}

	return &url, nil
}

func ParseEventHubSASConnectionString(connString string) (map[string]string, error) {
	// This connection string was for a real Event Hub which has been deleted
	// so its safe to include here for reference to understand the format.
	// Endpoint=sb://example-ehn.servicebus.windows.net/;SharedAccessKeyName=example-ehar;SharedAccessKey=DzGpfdyJda9D/xIkME0FLA66wZnheOBID0s1/rrtlHg=;EntityPath=example-eh
	validKeys := map[string]bool{"Endpoint": true, "SharedAccessKeyName": true, "SharedAccessKey": true, "EntityPath": true}
	// The k-v pairs are separated with semi-colons
	tokens := strings.Split(connString, ";")

	kvp := make(map[string]string)

	for _, atoken := range tokens {
		// The individual k-v are separated by an equals sign.
		kv := strings.SplitN(atoken, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("[ERROR] token `%s` is an invalid key=pair (connection string %s)", atoken, connString)
		}

		key := kv[0]
		val := kv[1]

		if _, present := validKeys[key]; !present {
			return nil, fmt.Errorf("[ERROR] Unknown Key `%s` in connection string %s", key, connString)
		}
		kvp[key] = val
	}

	if _, present := kvp[connStringSharedAccessKeyKey]; !present {
		return nil, fmt.Errorf("[ERROR] Shared Access Key not found in connection string: %s", connString)
	}

	return kvp, nil
}
