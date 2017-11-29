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

// Test is a simple wrapper of Terraform's testing framework.
//
// It respects the all the use conventions defined by Terraform.
//
// If "TF_PARA" is set to some non-empty value, all the parallel test cases will
// be run parallel.
// example usage to run parallel acceptance tests with n threads:
// TF_ACC=1 TF_PARA=1 go test [TEST] [TESTARGS] -v -parallel=n
func Test(t AzTestT, c tftest.TestCase) {
	if os.Getenv(ParaTestEnvVar) != "" {
		t.Parallel()
	}

	tftest.Test(t, c)
}
