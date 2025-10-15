package main

import (
	"rampage/system"

	"github.com/argus-labs/monorepo/pkg/cardinal"
)

// change
func main() {
	world, err := cardinal.NewWorld(cardinal.WorldOptions{
		TickRate:       1,
		EpochFrequency: 10,
	})
	if err != nil {
		panic(err.Error())
	}

	cardinal.RegisterSystem(world, system.PlayerSpawnerSystem, cardinal.WithHook(cardinal.Init))

	cardinal.RegisterSystem(world, system.CreatePlayerSystem)
	cardinal.RegisterSystem(world, system.RegenSystem)
	cardinal.RegisterSystem(world, system.AttackPlayerSystem)
	cardinal.RegisterSystem(world, system.GraveyardSystem)
	cardinal.RegisterSystem(world, system.CallExternalSystem)

	world.StartGame()
}
