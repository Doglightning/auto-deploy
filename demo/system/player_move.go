package system

import (
	"time"

	"github.com/argus-labs/monorepo/pkg/cardinal"
	"demo/component"
	"demo/event"
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

			if msg.ArgusAuthID != tag.ArgusAuthID {
				continue
			}

			isOnline := player.Online.Get().Online

			if !isOnline {
				state.PlayerSpawnEvent.Emit(event.PlayerSpawn{
					ArgusAuthID:   tag.ArgusAuthID,
					ArgusAuthName: tag.ArgusAuthName,
					X:             msg.X,
					Y:             msg.Y,
				})
			}

			player.Position.Set(component.Position{X: int(msg.X), Y: int(msg.Y)})
			player.Online.Set(component.OnlineStatus{Online: true, LastActive: time.Now()})

			state.PlayerMovementEvent.Emit(event.PlayerMovement{
				ArgusAuthID: tag.ArgusAuthID,
				X:           msg.X,
				Y:           msg.Y,
			})

			name := tag.ArgusAuthName

			state.Logger().Info().
				Uint32("entity", uint32(entity.ID)).
				Msgf("Player %s (id: %s) moved to %d, %d", name, tag.ArgusAuthID, msg.X, msg.Y)
		}
	}
	return nil
}
