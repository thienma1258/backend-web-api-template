package api

// CommitHash is a commit hash from git which is overridden
// by CI server when building Docker image.
var (
	CommitHash = "dev"
)
