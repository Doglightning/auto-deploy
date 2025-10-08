package system

import (
	"github.com/argus-labs/monorepo/pkg/cardinal"
	"rampage/component"
)

type PlayerSearch = cardinal.Exact[struct {
	Tag    cardinal.Ref[component.PlayerTag]
	Health cardinal.Ref[component.Health]
}]

type GraveSearch = cardinal.Exact[struct {
	Grave cardinal.Ref[component.Gravestone]
}]
