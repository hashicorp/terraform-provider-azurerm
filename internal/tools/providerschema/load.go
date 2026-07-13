// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package providerschema

import (
	"encoding/json"
	"os"
)

// LoadWrapperFromFile decodes an exported provider schema dump (as written by
// the schema-api -export command) into a ProviderWrapper.
func LoadWrapperFromFile(path string) (*ProviderWrapper, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	wrapper := &ProviderWrapper{}
	if err := json.NewDecoder(f).Decode(wrapper); err != nil {
		return nil, err
	}

	return wrapper, nil
}
