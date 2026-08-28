package webassets

import "embed"

// FS includes the templates and static assets needed at runtime.
//
//go:embed static/* templates/*.html templates/partials/*.html
var FS embed.FS
