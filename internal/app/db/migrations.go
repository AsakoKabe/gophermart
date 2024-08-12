package db

import "embed"

//go:embed migrations/*.sql
var MigrationsFolder embed.FS
