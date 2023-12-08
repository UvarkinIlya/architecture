module architecture/modellibrary

go 1.19

replace architecture/logger => ../logger

require (
	architecture/logger v0.0.0-00010101000000-000000000000
	github.com/nats-io/nats.go v1.31.0
)

require (
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/nats-io/nkeys v0.4.5 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
