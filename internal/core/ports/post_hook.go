package ports

type PostHook interface {
	// Run executes post-generation steps inside the target directory (e.g., go mod tidy, git init).
	Run(targetDir string) error
}
