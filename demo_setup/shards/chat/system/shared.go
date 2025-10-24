package system

import (
	"demo_setup/shards/chat/component"

	"github.com/argus-labs/world-engine/pkg/cardinal"
)

type ChatSearch = cardinal.Exact[struct {
	UserTag cardinal.Ref[component.UserTag]
	Chat    cardinal.Ref[component.Chat]
}]
