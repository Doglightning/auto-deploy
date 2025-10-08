package system

import (
	"time"

	"github.com/argus-labs/monorepo/pkg/cardinal"
	"demo/component"
	"demo/event"
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
			ArgusAuthID:   msg.ArgusAuthID,
			ArgusAuthName: msg.ArgusAuthName,
			X:             msg.X,
			Y:             msg.Y,
		})

		if playerSet.Exists(msg.ArgusAuthID) {
			state.Logger().Info().Msgf("Player with ID %s already exists, skipping creation", msg.ArgusAuthID)
			continue
		}

		id, err := state.Players.Create(
			component.PlayerTag{ArgusAuthID: msg.ArgusAuthID, ArgusAuthName: msg.ArgusAuthName},
			component.Position{X: int(msg.X), Y: int(msg.Y)},
			component.OnlineStatus{Online: true, LastActive: time.Now()},
		)
		playerSet.Add(msg.ArgusAuthID)

		if err != nil {
			// If we return the error, Cardinal will shutdown, so just log it.
			state.Logger().Error().Err(err).Msg("error creating entity")
			continue
		}

		state.Logger().Info().
			Uint32("entity", uint32(id)).
			Msgf("Created player %s (id: %s)", msg.ArgusAuthName, msg.ArgusAuthID)
	}
	return nil
}
