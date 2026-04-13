package siteassets

import "embed"

// Dist contains the production web assets baked into the server binary.
//
//go:embed all:dist
var Dist embed.FS
