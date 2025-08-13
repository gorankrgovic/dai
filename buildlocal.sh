VERSION=v0.1.0
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT=$(git rev-parse HEAD)

go build -trimpath -ldflags "-s -w \
  -X github.com/gorankrgovic/dai/internal/buildinfo.Version=${VERSION} \
  -X github.com/gorankrgovic/dai/internal/buildinfo.Commit=${COMMIT} \
  -X github.com/gorankrgovic/dai/internal/buildinfo.Date=${DATE}" \
  -o dai