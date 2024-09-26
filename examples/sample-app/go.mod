module sample-app

go 1.23.0

require (
	github.com/marnixbouhuis/confpb v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.34.3-0.20240816073751-94ecbc261689
)

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/marnixbouhuis/confpb => ./../../
