package migrations

import "embed"

// Files contains the versioned schema, shipped inside the gateway binary.
//
//go:embed *.sql
var Files embed.FS

const LatestVersion int64 = 2
