module main

go 1.21.0

require (
	github.com/opencontainers/image-spec v1.1.0-rc4
	oras.land/oras-go/v2 v2.0.0-00010101000000-000000000000
)

require (
	github.com/opencontainers/go-digest v1.0.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
)

replace oras.land/oras-go/v2 => github.com/Wwwsylvia/oras-go/v2 v2.0.0-20230913113822-65603174f3b3
