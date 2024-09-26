package version

var (
	// GitReleaseTag is the version tag injected by goreleaser at build time.
	// Note: this is not available when using confpb as a library.
	GitReleaseTag     = "unknown"
	GitReleaseVersion = "unknown"
	GitShortCommit    = "unknown"
)
