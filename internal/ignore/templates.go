package ignore

// DefaultDaiIgnore returns the default .daiignore template (gitignore syntax).
func DefaultDaiIgnore() string {
	return `# build artefacts and deps
dist/
build/
node_modules/
vendor/

# lock files
*.lock
pnpm-lock.yaml
package-lock.json
composer.lock
poetry.lock

# generated files / maps / snapshots
*.min.js
*.map
*.snap

# binary / assets
*.png
*.jpg
*.jpeg
*.gif
*.webp
*.pdf
*.mp4
*.mp3
*.wav
*.zip

# test fixture
testdata/
fixtures/
`
}
