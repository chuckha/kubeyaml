package messages

// VersionsResponse is the struct defining what a resposne to /versions will be.
type VersionsResponse struct {
	Versions       []string
	DefaultVersion string
}
