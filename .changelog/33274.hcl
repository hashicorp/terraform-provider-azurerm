
change "other-fix" {
  body = "`go-azure-sdk` - `Delete` operations now poll on asynchronous operation URLs if returned by the API instead of only checking for a `404` on the resource URL, ensuring deletion errors are reported to the user"
}
change "dependency" {
  body = "dependencies: `go-azure-sdk` - update to `v0.20260901.1173158`"
}
