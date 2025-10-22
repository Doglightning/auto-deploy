package system

import (
	"time"

	"demo/component"
	"demo/event"

	"github.com/argus-labs/world-engine/pkg/cardinal"
)

type MovePlayerCommand struct {
	cardinal.BaseCommand
	ArgusAuthID string `json:"argus_auth_id"`
	X           uint32 `json:"x"`
	Y           uint32 `json:"y"`
}

func (a MovePlayerCommand) Name() string {
	return "move-player"
}

type MovePlayerSystemState struct {
	cardinal.BaseSystemState
	MovePlayerCommands  cardinal.WithCommand[MovePlayerCommand]
	PlayerSpawnEvent    cardinal.WithEvent[event.PlayerSpawn]
	PlayerMovementEvent cardinal.WithEvent[event.PlayerMovement]
	Players             PlayerSearch
}

func MovePlayerSystem(state *MovePlayerSystemState) error {
	for msg := range state.MovePlayerCommands.Iter() {
		for entity, player := range state.Players.Iter() {
			tag := player.Tag.Get()

			if msg.Payload().ArgusAuthID != tag.ArgusAuthID {
				continue
			}

			isOnline := player.Online.Get().Online

			if !isOnline {
				state.PlayerSpawnEvent.Emit(event.PlayerSpawn{
					ArgusAuthID:   tag.ArgusAuthID,
					ArgusAuthName: tag.ArgusAuthName,
					X:             msg.Payload().X,
					Y:             msg.Payload().Y,
				})
			}

			player.Position.Set(component.Position{X: int(msg.Payload().X), Y: int(msg.Payload().Y)})
			player.Online.Set(component.OnlineStatus{Online: true, LastActive: time.Now()})

			state.PlayerMovementEvent.Emit(event.PlayerMovement{
				ArgusAuthID: tag.ArgusAuthID,
				X:           msg.Payload().X,
				Y:           msg.Payload().Y,
			})

			name := tag.ArgusAuthName

			state.Logger().Info().
				Uint32("entity", uint32(entity.ID)).
				Msgf("Player %s (id: %s) moved to %d, %d", name, tag.ArgusAuthID, msg.Payload().X, msg.Payload().Y)
		}
	}
	return nil
}
