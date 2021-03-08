package lock

import (
	_ "embed"
)

var (

	//go:embed lock.lua
	LOCK_LUA string

	//go:embed release.lua
	RELEASE_LUA string
)
