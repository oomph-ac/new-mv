module github.com/oomph-ac/new-mv

go 1.22.1

toolchain go1.23.2

require (
	github.com/df-mc/dragonfly v0.9.18-0.20240814140312-13b68f1ec242
	github.com/df-mc/worldupgrader v1.0.18
	github.com/google/uuid v1.6.0
	github.com/rogpeppe/go-internal v1.12.0
	github.com/samber/lo v1.38.1
	github.com/sandertv/go-raknet v1.14.2
	github.com/sandertv/gophertunnel v1.42.0
	github.com/segmentio/fasthash v1.0.3
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c
	golang.org/x/image v0.21.0
)

require (
	github.com/brentp/intintmap v0.0.0-20190211203843-30dc0ade9af9 // indirect
	github.com/df-mc/goleveldb v1.1.9 // indirect
	github.com/gameparrot/goquery v0.2.0 // indirect
	github.com/go-gl/mathgl v1.1.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/muhammadmuzzammil1998/jsonc v1.0.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/oauth2 v0.23.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)

replace github.com/df-mc/dragonfly v0.9.18-0.20240814140312-13b68f1ec242 => ../dragonfly

replace github.com/sandertv/go-raknet => ../go-raknet

replace github.com/sandertv/gophertunnel => ../gophertunnel
