package accounts

type AddStorageAccountWithAccountParameters struct {
	Name       string                      `json:"name"`
	Properties AddStorageAccountProperties `json:"properties"`
}
