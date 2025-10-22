package system

import (
	"demo/event"

	"github.com/argus-labs/world-engine/pkg/cardinal"
)

type PlayerLeaveCommand struct {
	cardinal.BaseCommand
	ArgusAuthID string `json:"argus_auth_id"`
}

func (p PlayerLeaveCommand) Name() string {
	return "player-leave"
}

type PlayerLeaveSystemState struct {
	cardinal.BaseSystemState
	PlayerLeaveCommands  cardinal.WithCommand[PlayerLeaveCommand]
	PlayerDepartureEvent cardinal.WithEvent[event.PlayerDeparture]
	Players              PlayerSearch
}

// PlayerLeaveSystem is called when a player leaves a quadrant (e.g. to join another quadrant).
func PlayerLeaveSystem(state *PlayerLeaveSystemState) error {
	players := make(map[string]cardinal.Entity)

	for entity, player := range state.Players.Iter() {
		players[player.Tag.Get().ArgusAuthID] = entity
	}

	for msg := range state.PlayerLeaveCommands.Iter() {
		entity, exists := players[msg.Payload().ArgusAuthID]
		if !exists {
			state.Logger().Info().Msgf("Player with ID %s not found", msg.Payload().ArgusAuthID)
			continue
		}

		entity.Destroy()

		state.PlayerDepartureEvent.Emit(event.PlayerDeparture{
			ArgusAuthID: msg.Payload().ArgusAuthID,
		})
	}
	return nil
}
