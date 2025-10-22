package system

import (
	"time"

	"demo/component"
	"demo/event"

	"github.com/argus-labs/world-engine/pkg/cardinal"
)

type PlayerSpawnCommand struct {
	cardinal.BaseCommand
	ArgusAuthID   string `json:"argus_auth_id"`
	ArgusAuthName string `json:"argus_auth_name"`
	X             uint32 `json:"x"`
	Y             uint32 `json:"y"`
}

func (a PlayerSpawnCommand) Name() string {
	return "player-spawn"
}

type SpawnPlayerSystemState struct {
	cardinal.BaseSystemState
	SpawnPlayerCommands cardinal.WithCommand[PlayerSpawnCommand]
	PlayerSpawnEvent    cardinal.WithEvent[event.PlayerSpawn]
	Players             PlayerSearch
}

func PlayerSpawnSystem(state *SpawnPlayerSystemState) error {
	for msg := range state.SpawnPlayerCommands.Iter() {
		// Regardless of whether the player exists or not, we emit a spawn event
		// Because the act of spawning is also creating (if they don’t already exist)
		state.PlayerSpawnEvent.Emit(event.PlayerSpawn{
			ArgusAuthID:   msg.Payload().ArgusAuthID,
			ArgusAuthName: msg.Payload().ArgusAuthName,
			X:             msg.Payload().X,
			Y:             msg.Payload().Y,
		})

		if playerSet.Exists(msg.Payload().ArgusAuthID) {
			state.Logger().Info().Msgf("Player with ID %s already exists, skipping creation", msg.Payload().ArgusAuthID)
			continue
		}

		id, err := state.Players.Create(
			component.PlayerTag{ArgusAuthID: msg.Payload().ArgusAuthID, ArgusAuthName: msg.Payload().ArgusAuthName},
			component.Position{X: int(msg.Payload().X), Y: int(msg.Payload().Y)},
			component.OnlineStatus{Online: true, LastActive: time.Now()},
		)
		playerSet.Add(msg.Payload().ArgusAuthID)

		if err != nil {
			// If we return the error, Cardinal will shutdown, so just log it.
			state.Logger().Error().Err(err).Msg("error creating entity")
			continue
		}

		state.Logger().Info().
			Uint32("entity", uint32(id)).
			Msgf("Created player %s (id: %s)", msg.Payload().ArgusAuthName, msg.Payload().ArgusAuthID)
	}
	return nil
}
