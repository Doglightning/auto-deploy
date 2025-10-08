package system

import (
	"github.com/argus-labs/monorepo/pkg/cardinal"
	"rampage/component"
	systemevent "rampage/system_event"
)

type GraveyardSystemState struct {
	cardinal.BaseSystemState
	PlayerDeathSystemEvents cardinal.WithSystemEventReceiver[systemevent.PlayerDeath]
	Graves                  GraveSearch
}

func GraveyardSystem(state *GraveyardSystemState) error {
	for event := range state.PlayerDeathSystemEvents.Iter() {
		_, _ = state.Graves.Create(component.Gravestone{Nickname: event.Nickname})

		state.Logger().Info().Msgf("Created grave stone for player %s", event.Nickname)
	}
	return nil
}
