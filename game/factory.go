package game

import (
	"fmt"
)

type Factory struct {
	Resource           Resource
	Level              uint32
	UpgradeSecondsLeft *uint32
}

func (f *Factory) TickSecond(resources *Resources) {
	properties := f.properties()
	if f.Resource == Iron {
		resources.Iron += properties.ironPerSecond
	}
	if f.Resource == Copper {
		resources.Copper += properties.copperPerSecond
	}
	if f.UpgradeSecondsLeft != nil {
		fmt.Printf("left: %v\n", *f.UpgradeSecondsLeft)
		if *f.UpgradeSecondsLeft == 0 {
			f.Level++
			f.UpgradeSecondsLeft = nil
		} else {
			*f.UpgradeSecondsLeft--
		}
	}
}

func (f *Factory) TickMinute(resources *Resources) {
	properties := f.properties()
	if f.Resource == Gold {
		resources.Gold += properties.goldPerMinute
	}
}

// TODO Increase this to 5 later, or maybe just three
const maxLevel = 2

func (f *Factory) Upgrade(resources *Resources) bool {
	if f.Level == maxLevel || f.Upgrading() {
		return false
	}
	properties := f.properties()
	cost := properties.upgradeCost
	upgradeSeconds := properties.upgradeSeconds
	if resources.Iron >= cost.Iron &&
		resources.Copper >= cost.Copper &&
		resources.Gold >= cost.Gold {
		resources.Iron -= cost.Iron
		resources.Copper -= cost.Copper
		resources.Gold -= cost.Gold
		f.UpgradeSecondsLeft = &upgradeSeconds
		return true
	}
	return false
}

func (f *Factory) Upgrading() bool {
	return f.UpgradeSecondsLeft != nil
}

type properties struct {
	ironPerSecond   uint32
	copperPerSecond uint32
	goldPerMinute   uint32
	upgradeSeconds  uint32
	upgradeCost     Resources
}

func (f Factory) properties() properties {
	switch {
	case f.Resource == Iron && f.Level == 1:
		return properties{
			ironPerSecond:  10,
			upgradeSeconds: 15,
			upgradeCost: Resources{
				Iron:   300,
				Copper: 100,
				Gold:   1,
			}}
	case f.Resource == Iron && f.Level == 2:
		return properties{
			ironPerSecond:  20,
			upgradeSeconds: 30,
			upgradeCost: Resources{
				Iron:   800,
				Copper: 250,
				Gold:   2,
			}}
	case f.Resource == Copper && f.Level == 1:
		return properties{
			copperPerSecond: 3,
			upgradeSeconds:  15,
			upgradeCost: Resources{
				Iron:   300,
				Copper: 70,
			}}
	case f.Resource == Copper && f.Level == 2:
		return properties{
			copperPerSecond: 7,
			upgradeSeconds:  30,
			upgradeCost: Resources{
				Iron:   400,
				Copper: 150,
			}}
	case f.Resource == Gold && f.Level == 1:
		return properties{
			goldPerMinute:  2,
			upgradeSeconds: 15,
			upgradeCost: Resources{
				Copper: 100,
				Gold:   2,
			}}
	case f.Resource == Gold && f.Level == 2:
		return properties{
			goldPerMinute:  3,
			upgradeSeconds: 30,
			upgradeCost: Resources{
				Copper: 200,
				Gold:   4,
			}}
	}

	panic(fmt.Sprintf("No case when upgrading resource %v, level %v", f.Resource, f.Level))
}
