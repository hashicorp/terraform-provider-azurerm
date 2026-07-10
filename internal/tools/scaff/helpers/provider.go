package helpers

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func ProviderName() (providerName *string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not determine current working directory to evaluate provider name: %+v", err)
	}
	ps := "/"
	if runtime.GOOS == "windows" {
		ps = "\\"
	}
	wdParts := strings.Split(cwd, ps)
	fullProviderName := wdParts[len(wdParts)-1]

	n := strings.TrimPrefix(fullProviderName, "terraform-provider-")

	providerName = &n

	return
}
