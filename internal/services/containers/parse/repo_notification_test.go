package parse

import "testing"

func TestRepositoryNotification(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *repositoryNotification
	}{

		{
			// empty
			Input: "",
			Error: true,
		},
		{
			// name only
			Input: "foo",
			Error: true,
		},
		{
			// digest but no action
			Input: "foo@sha256:d88ff149d60584cd1dab334761d8b971d318e4417e488bc6201e95719f339b58",
			Error: true,
		},
		{
			// tag but no action
			Input: "foo:latest",
			Error: true,
		},
		{
			// no digest or tag
			Input: "foo:delete",
			Expected: &repositoryNotification{
				Artifact: Artifact{
					Name: "foo",
				},
				Action: RepositoryNotificationActionDelete,
			},
		},
		{
			// digest
			Input: "foo@sha256:d88ff149d60584cd1dab334761d8b971d318e4417e488bc6201e95719f339b58:delete",
			Expected: &repositoryNotification{
				Artifact: Artifact{
					Name:   "foo",
					Digest: "sha256:d88ff149d60584cd1dab334761d8b971d318e4417e488bc6201e95719f339b58",
				},
				Action: RepositoryNotificationActionDelete,
			},
		},
		{
			// tag
			Input: "foo:latest:delete",
			Expected: &repositoryNotification{
				Artifact: Artifact{
					Name: "foo",
					Tag:  "latest",
				},
				Action: RepositoryNotificationActionDelete,
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := RepositoryNotification(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if *actual != *v.Expected {
			t.Fatalf("Expected %v but got %v", *v.Expected, *actual)
		}
	}
}
