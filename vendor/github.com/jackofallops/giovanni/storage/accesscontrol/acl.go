package accesscontrol

import (
	"strings"
)

type ACL struct {
	Entries []ACE
}

// Validate checks the ACL. Returns nil on success
func (acl *ACL) Validate() error {

	// TODO
	//   - check each user/group is only listed once (per default/non-default)

	for _, v := range acl.Entries {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// ParseACL parses an ACL string
func ParseACL(input string) (ACL, error) {

	aceStrings := strings.Split(input, ",")
	entries := make([]ACE, len(aceStrings))

	for i := 0; i < len(aceStrings); i++ {
		aceString := aceStrings[i]
		entry, err := ParseACE(aceString)
		if err != nil {
			return ACL{}, err
		}
		entries[i] = entry
	}
	return ACL{Entries: entries}, nil
}

// String returns the string form of the ACL - this does not check that it is valid
func (acl *ACL) String() string {

	aceStrings := make([]string, len(acl.Entries))

	for i := 0; i < len(acl.Entries); i++ {
		ace := acl.Entries[i]
		aceStrings[i] = ace.String()
	}
	return strings.Join(aceStrings, ",")
}
