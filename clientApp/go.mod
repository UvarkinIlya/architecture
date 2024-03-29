module architecture/clientApp

go 1.19

require (
	architecture/logger v0.0.0-00010101000000-000000000000
	architecture/modellibrary v0.0.0-00010101000000-000000000000
)

replace architecture/logger => ../logger

replace architecture/modellibrary => ../modellibrary

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
