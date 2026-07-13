// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package differ

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func (d *Differ) loadFromFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := &providerschema.ProviderWrapper{}
	// TODO - Custom marshalling to fix the type assertions later? meh, works for now...
	if err := json.NewDecoder(f).Decode(buf); err != nil {
		return err
	}
	d.base = buf

	return nil
}

func (d *Differ) loadFromProvider(data *providerschema.ProviderJSON, providerName string) error {
	if s, err := providerschema.ProviderFromRaw(data); err != nil {
		return err
	} else {
		d.current = &providerschema.ProviderWrapper{
			ProviderName:   providerName,
			ProviderSchema: s,
		}
	}
	return nil
}
