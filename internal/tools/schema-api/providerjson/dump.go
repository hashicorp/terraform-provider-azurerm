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
