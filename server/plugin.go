package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	botID string
}

const (
	// todo - can i better find this ID?
	PluginId = "com.mattermost.plugin-starter-template"
)

const (
	commandSupport = "supporty"
)

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, p.API.GetConfig())
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

func (p *Plugin) OnActivate() error {
	BotInfo := &model.Bot{
		Username:    "supporty",
		DisplayName: "Supporty the Support bot",
		Description: "The source of truth for customer info",
	}

	botID, err := p.API.EnsureBotUser(BotInfo)

	if err != nil {
		return errors.Wrapf(err, "failed to ensure bot")
	}

	p.botID = botID

	p.API.CreatePost(&model.Post{
		UserId:    botID,
		ChannelId: "jqqzh76e43bc9jmuxq1ycee6rr",
		Message:   "hey, bud.",
	})

	if err := p.registerCommands(); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	return nil
}

func (p *Plugin) registerCommands() error {
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandSupport,
		Method:           "G",
		PluginId:         PluginId,
		AutoComplete:     true,
		AutoCompleteDesc: "Supporty over here doing things",
		// doesn't seem needed
		// AutoCompleteHint: "boom boom support",
		DisplayName: "Support name",
	}); err != nil {
		return errors.Wrapf(err, "failed to register plugin slash")
	}

	return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")
	switch trigger {
	case commandSupport:
		return p.executeSupportCommand(args), nil

	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown command: " + args.Command),
		}, nil
	}
}

func (p *Plugin) executeSupportCommand(args *model.CommandArgs) *model.CommandResponse {
	p.API.CreatePost(&model.Post{
		ChannelId: args.ChannelId,
		UserId:    p.botID,
		Message:   "bruhhhhh, what's up?",
	})
	return &model.CommandResponse{}
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	if post.UserId == p.botID || post.RootId != "" || len(post.FileIds) == 0 {
		return
	}

	supportPackets, names := PostContainsSupportPackage(p, post)

	if len(supportPackets) == 0 {
		return
	}

	p.API.CreatePost(&model.Post{
		ChannelId: post.ChannelId,
		RootId:    post.Id,
		UserId:    p.botID,
		Message:   "Uploading support packet for " + strings.Join(names, " ,"),
	})

	ProcessSupportPackets(p, supportPackets, post)
}
