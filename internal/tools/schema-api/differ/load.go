package differ

import (
	"encoding/json"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
	"os"
)

func (d *Differ) loadFromFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := &providerjson.ProviderWrapper{}
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
