package accesscontrol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestACLValidate_WithValidACL(t *testing.T) {
	acl := ACL{
		Entries: []ACE{
			{
				TagType:     TagTypeUser,
				Permissions: "rwx",
			},
			{
				TagType:     TagTypeGroup,
				Permissions: "r-x",
			},
			{
				TagType:     TagTypeOther,
				Permissions: "---",
			},
		},
	}
	assert.Nil(t, acl.Validate(), "Expected ACE to Validate successfully")
}

func TestACLValidate_WithInvalidACL(t *testing.T) {
	acl := ACL{
		Entries: []ACE{
			{
				TagType:     TagType("wibble"),
				Permissions: "rwx",
			},
			{
				TagType:     TagTypeGroup,
				Permissions: "r-x",
			},
			{
				TagType:     TagTypeOther,
				Permissions: "---",
			},
		},
	}
	assert.Error(t, acl.Validate(), "Expected ACE to fail to Validate")
}

func TestACLString(t *testing.T) {
	acl := ACL{
		Entries: []ACE{
			{
				TagType:     TagTypeUser,
				Permissions: "rwx",
			},
			{
				TagType:     TagTypeGroup,
				Permissions: "r-x",
			},
			{
				TagType:     TagTypeOther,
				Permissions: "---",
			},
		},
	}
	assert.Equal(t, "user::rwx,group::r-x,other::---", acl.String())
}

func TestACLParse_WithValidACL(t *testing.T) {
	expected := ACL{
		Entries: []ACE{
			{
				TagType:     TagTypeUser,
				Permissions: "rwx",
			},
			{
				TagType:     TagTypeGroup,
				Permissions: "r-x",
			},
			{
				TagType:     TagTypeOther,
				Permissions: "---",
			},
		},
	}
	acl, err := ParseACL("user::rwx,group::r-x,other::---")
	assert.Nil(t, err, "Expected ACE to parse successfully")
	assert.Equal(t, expected, acl)
}

func TestACLParse_WithInvalidACL(t *testing.T) {
	_, err := ParseACL("user:rwx,group::r-x,other::---")
	assert.EqualError(t, err, "ACE string should have either 3 or 4 parts")
}
