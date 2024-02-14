module github.com/crossworth/pkgs/ptr

go 1.22.0

require (
	github.com/crossworth/pkgs/floats v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/crossworth/pkgs/floats => ../floats
