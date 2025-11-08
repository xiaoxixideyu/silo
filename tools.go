//go:build tools
// +build tools

package main

import (
	_ "accessor"

	_ "github.com/google/wire/cmd/wire"
	_ "github.com/jmattheis/goverter/cmd/goverter"
	_ "github.com/swaggo/swag"
	_ "golang.org/x/tools/cmd/stringer"
)
