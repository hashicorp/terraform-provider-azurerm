// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package providerjson

import (
	"encoding/json"
	"os"
)

func DumpWithWrapper(wrapper *ProviderWrapper, data *ProviderJSON) error {
	if s, err := ProviderFromRaw(data); err != nil {
		return err
	} else {
		wrapper.ProviderSchema = s
	}

	if err := json.NewEncoder(os.Stdout).Encode(wrapper); err != nil {
		return err
	}

	return nil
}

func WriteWithWrapper(wrapper *ProviderWrapper, data *ProviderJSON, filename string) error {
	if s, err := ProviderFromRaw(data); err != nil {
		return err
	} else {
		wrapper.ProviderSchema = s
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(wrapper); err != nil {
		return err
	}

	return nil
}
