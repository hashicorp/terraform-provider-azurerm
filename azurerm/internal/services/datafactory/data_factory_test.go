package datafactory

import "testing"

func TestAzureRmDataFactoryLinkedServiceConnectionStringDiff(t *testing.T) {
	cases := []struct {
		Old    string
		New    string
		NoDiff bool
	}{
		{
			Old:    "",
			New:    "",
			NoDiff: true,
		},
		{
			Old:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test",
			New:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test;Password=test",
			NoDiff: true,
		},
		{
			Old:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test",
			New:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test",
			NoDiff: true,
		},
		{
			Old:    "Integrated Security=False;Data Source=test2;Initial Catalog=test;User ID=test",
			New:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test;Password=test",
			NoDiff: false,
		},
	}

	for _, tc := range cases {
		noDiff := azureRmDataFactoryLinkedServiceConnectionStringDiff("", tc.Old, tc.New, nil)

		if noDiff != tc.NoDiff {
			t.Fatalf("Expected azureRmDataFactoryLinkedServiceConnectionStringDiff to be '%t' for '%s' '%s' - got '%t'", tc.NoDiff, tc.Old, tc.New, noDiff)
		}
	}
}
