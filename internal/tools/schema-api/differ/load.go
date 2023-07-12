// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package differ

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

func (d *Differ) loadFromFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := &providerjson.ProviderWrapper{}
	// TODO - Custom marshalling to fix the type assertions later? meh, works for now...
	if err := json.NewDecoder(f).Decode(buf); err != nil {
		return err
	}
	d.base = buf

	return nil
}

func (d *Differ) loadFromProvider(data *providerjson.ProviderJSON, providerName string) error {
	if s, err := providerjson.ProviderFromRaw(data); err != nil {
		return err
	} else {
		d.current = &providerjson.ProviderWrapper{
			ProviderName:   providerName,
			ProviderSchema: s,
		}
	}
	return nil
}
