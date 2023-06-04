package modify

import (
	"github.com/kercre123/wire-os/wired/pkg/modify/modifiers"
	"github.com/kercre123/wire-os/wired/pkg/vars"
)

var WireOSModifiers []vars.Modifier = []vars.Modifier{
	{
		Name:        "HigherPerformance",
		Description: "Increases the CPU and RAM frequencies to maximum potential",
		Apply:       modifiers.HigherPerformance_Apply,
		Remove:      modifiers.HigherPerformance_Remove,
		HasInitFunc: true,
		Init:        modifiers.HigherPerformance_Init,
	},
	{
		Name:        "NoSnore",
		Description: "Disables snoring",
		Apply:       modifiers.NoSnore_Apply,
		Remove:      modifiers.NoSnore_Remove,
		HasInitFunc: false,
	},
}
