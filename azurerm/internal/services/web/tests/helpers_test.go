package tests

import (
	"os"
)

func skipStaticSite() bool {
	return os.Getenv("ARM_TEST_GITHUB_TOKEN") == "" || os.Getenv("ARM_TEST_GITHUB_REPO") == ""
}
