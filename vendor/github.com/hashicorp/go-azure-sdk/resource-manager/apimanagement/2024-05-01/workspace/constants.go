package workspace

type Confirmation string

const (
	ConfirmationInvite Confirmation = "invite"
	ConfirmationSignup Confirmation = "signup"
)

type UserState string

const (
	UserStateActive  UserState = "active"
	UserStateBlocked UserState = "blocked"
	UserStatePending UserState = "pending"
)
