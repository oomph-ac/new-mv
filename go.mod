module github.com/oomph-ac/new-mv

go 1.22

toolchain go1.22.1

require (
	github.com/df-mc/dragonfly v0.9.18-0.20240814140312-13b68f1ec242
	github.com/df-mc/worldupgrader v1.0.16
	github.com/go-gl/mathgl v1.1.0
	github.com/google/uuid v1.6.0
	github.com/rogpeppe/go-internal v1.12.0
	github.com/samber/lo v1.38.1
	github.com/sandertv/go-raknet v1.14.1
	github.com/sandertv/gophertunnel v1.40.1
	github.com/segmentio/fasthash v1.0.3
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/exp v0.0.0-20240808152545-0cdaa3abc0fa
	golang.org/x/image v0.19.0
	golang.org/x/oauth2 v0.22.0
)

require (
	github.com/brentp/intintmap v0.0.0-20190211203843-30dc0ade9af9 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/df-mc/atomic v1.10.0 // indirect
	github.com/df-mc/goleveldb v1.1.9 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/muhammadmuzzammil1998/jsonc v1.0.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
)

replace github.com/sandertv/go-raknet => github.com/tedacmc/tedac-raknet v0.0.4

replace github.com/sandertv/gophertunnel => ../gophertunnel