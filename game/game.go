package game

import (
	"fmt"
	"time"
)

type User string

type Resource uint32

const (
	Iron Resource = iota
	Copper
	Gold
)

type Resources struct {
	Iron   uint32
	Copper uint32
	Gold   uint32
}

type State struct {
	User      User
	Resources *Resources
	Factories []*Factory
}

func (s State) clone() State {
	resourcesCopy := *s.Resources
	factories := []*Factory{}
	for _, f := range s.Factories {
		fCopy := *f
		if f.UpgradeSecondsLeft != nil {
			upgradeSecondsCopy := *f.UpgradeSecondsLeft
			fCopy.UpgradeSecondsLeft = &upgradeSecondsCopy
		}
		factories = append(factories, &fCopy)
	}
	return State{
		User:      s.User,
		Resources: &resourcesCopy,
		Factories: factories,
	}
}

type StateMessage struct {
	Resp chan<- State
}

type UpgradeMessage struct {
	Resp     chan<- bool
	Resource Resource
}

type Game struct {
	StateChan   chan<- StateMessage
	UpgradeChan chan<- UpgradeMessage
}

func Start(user User) Game {
	state := State{
		User:      user,
		Resources: &Resources{},
		Factories: []*Factory{
			{
				Resource: Iron,
				Level:    1,
			},
			{
				Resource: Copper,
				Level:    1,
			},
			{
				Resource: Gold,
				Level:    1,
			},
		},
	}
	return Continue(state)
}

func Continue(state State) Game {
	// TODO For state message, probably copy the factories slice
	stateChan := make(chan StateMessage)
	upgradeChan := make(chan UpgradeMessage)

	go func() {
		secondTicker := time.NewTicker(time.Second)
		minuteTicker := time.NewTicker(time.Minute)

	mainLoop:
		for {
			select {
			case <-secondTicker.C:
				for _, factory := range state.Factories {
					factory.TickSecond(state.Resources)
				}
				fmt.Printf("%v\n", state.Resources)
				for _, f := range state.Factories {
					fmt.Printf("%v ", f.Level)
				}
				fmt.Println()
			case <-minuteTicker.C:
				for _, factory := range state.Factories {
					factory.TickMinute(state.Resources)
				}
			case upgrade := <-upgradeChan:
				for _, factory := range state.Factories {
					if factory.Resource == upgrade.Resource {
						upgrade.Resp <- factory.Upgrade(state.Resources)
						continue mainLoop
					}
				}
				upgrade.Resp <- false
			case m := <-stateChan:
				// Deep-copy the state to make sure we don't get any concurrency issues later.
				m.Resp <- state.clone()
			}

			// TODO Serialize the state
		}
	}()

	return Game{
		StateChan:   stateChan,
		UpgradeChan: upgradeChan,
	}
}
