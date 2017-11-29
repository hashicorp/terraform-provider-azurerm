package aztesting

import (
	"os"

	tftest "github.com/hashicorp/terraform/helper/resource"
)

// ParaTestEnvVar is used to enable parallelism in testing.
const ParaTestEnvVar = "TF_PARA"

// AzTestT is the interface used to handle the test lifecycle of a test.
//
// Users should just use a *testing.T object, which implements this.
type AzTestT interface {
	tftest.TestT
	Parallel()
}

// Test is a helper to force the acceptance testing harness to run in the
// normal unit test suite. This should only be used for resource that don't
// have any external dependencies.
func Test(t AzTestT, c tftest.TestCase) {
	if os.Getenv(ParaTestEnvVar) != "" {
		t.Parallel()
	}

	tftest.Test(t, c)
}
