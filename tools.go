// Package tools for go mod.
// Add here tool dependencies that are not part of the build.

// +build tools

package tools

import (
	_ "github.com/go-phorce/configen/cmd/configen"
	_ "github.com/go-phorce/cov-report/cmd/cov-report"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/godoc"
	_ "golang.org/x/tools/cmd/gorename"
	_ "golang.org/x/tools/cmd/guru"
	_ "golang.org/x/tools/cmd/stringer"
)
