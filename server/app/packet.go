package app

import (
	"github.com/mattermost/mattermost/server/public/model"
)

type PacketActionService interface {

	// MessageHasBeenPosted suggests playbooks to the user if triggered
	MessageHasBeenPosted(post *model.Post)
}
