package state

import (
	"github.com/docker/docker/api/types"
	"github.com/enescakir/emoji"
)

const (
	Running    = "running"
	Stopped    = "stopped"
	Exited     = "exited"
	Paused     = "paused"
	Restarting = "restarting"
)

func GetEmoji(container *types.Container) string {
	switch container.State {
	case Running:
		return emoji.OkHand.String()
	case Exited:
		return emoji.Door.String()
	case Paused:
		return emoji.PauseButton.String()
	case Restarting:
		return emoji.ClockwiseVerticalArrows.String()
	case Stopped:
		return emoji.StopSign.String()
	default:
		return emoji.QuestionMark.String()
	}
}
