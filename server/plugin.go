package main

import (
	"net/http"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/api"
	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	"github.com/coltoneshaw/mattermost-plugin-customers/server/bot"

	// "github.com/coltoneshaw/mattermost-plugin-customers/server/command"
	"github.com/coltoneshaw/mattermost-plugin-customers/server/config"
	"github.com/coltoneshaw/mattermost-plugin-customers/server/sqlstore"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	pluginapi "github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/mattermost/mattermost/server/public/pluginapi/cluster"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	handler *api.Handler
	config  *config.ServiceImpl

	botID string
	bot   *bot.Bot

	pluginAPI       *pluginapi.Client
	customerService app.CustomerService
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// ServeHTTP routes incoming HTTP requests to the plugin's REST API.
func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.handler.ServeHTTP(w, r)
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

func (p *Plugin) OnActivate() error {
	pluginAPIClient := pluginapi.NewClient(p.API, p.Driver)
	p.pluginAPI = pluginAPIClient

	p.config = config.NewConfigService(pluginAPIClient, manifest)

	logger := logrus.StandardLogger()
	pluginapi.ConfigureLogrus(logger, pluginAPIClient)

	botID, err := pluginAPIClient.Bot.EnsureBot(&model.Bot{
		Username:    "supporty",
		DisplayName: "Supporty the Support bot",
		Description: "The source of truth for customer info",
		OwnerId:     "supporty",
	},
	// pluginapi.ProfileImagePath("assets/plugin_icon.png"),
	)

	if err != nil {
		return errors.Wrapf(err, "failed to ensure bot")
	}

	err = p.config.UpdateConfiguration(func(c *config.Configuration) {
		c.BotUserID = botID
	})
	if err != nil {
		return errors.Wrapf(err, "failed save bot to config")
	}

	p.botID = botID

	apiClient := sqlstore.NewClient(pluginAPIClient)
	p.bot = bot.New(pluginAPIClient, p.config.GetConfiguration().BotUserID, p.config)

	sqlStore, err := sqlstore.New(apiClient)
	if err != nil {
		return errors.Wrapf(err, "failed creating the SQL store")
	}

	customerStore := sqlstore.NewCustomerStore(apiClient, sqlStore)
	p.handler = api.NewHandler(pluginAPIClient, p.config)

	p.customerService = app.NewCustomerService(customerStore, p.bot, pluginAPIClient)

	// Migrations use the scheduler, so they have to be run after playbookRunService and scheduler have started
	mutex, err := cluster.NewMutex(p.API, "CRM_Customers")
	if err != nil {
		return errors.Wrap(err, "failed creating cluster mutex")
	}
	mutex.Lock()
	if err = sqlStore.RunMigrations(); err != nil {
		mutex.Unlock()
		return errors.Wrap(err, "failed to run migrations")
	}
	mutex.Unlock()

	api.NewCustomerHandler(
		p.handler.APIRouter,
		p.customerService,
		pluginAPIClient,
		p.config,
	)

	return nil
}

func (p *Plugin) MessageHasBeenPosted(_ *plugin.Context, post *model.Post) {
	p.customerService.MessageHasBeenPosted(post)
}

// // ExecuteCommand executes a command that has been previously registered via the RegisterCommand.
// func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
// 	runner := command.NewCommandRunner(c, args, pluginapi.NewClient(p.API, p.Driver), p.bot, p.config)

// 	if err := runner.Execute(); err != nil {
// 		return nil, model.NewAppError("Customers.ExecuteCommand", "app.command.execute.error", nil, err.Error(), http.StatusInternalServerError)
// 	}

// 	return &model.CommandResponse{}, nil
// }

// func (p *Plugin) registerCommands() error {
// 	if err := p.API.RegisterCommand(&model.Command{
// 		Trigger:          commandSupport,
// 		Method:           "G",
// 		PluginId:         PluginId,
// 		AutoComplete:     true,
// 		AutoCompleteDesc: "Supporty over here doing things",
// 		// doesn't seem needed
// 		// AutoCompleteHint: "boom boom support",
// 		DisplayName: "Support name",
// 	}); err != nil {
// 		return errors.Wrapf(err, "failed to register plugin slash")
// 	}

// 	return nil
// }

// func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
// 	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")
// 	switch trigger {
// 	case commandSupport:
// 		return p.executeSupportCommand(args), nil

// 	default:
// 		return &model.CommandResponse{
// 			ResponseType: model.CommandResponseTypeEphemeral,
// 			Text:         fmt.Sprintf("Unknown command: " + args.Command),
// 		}, nil
// 	}
// }

// func (p *Plugin) executeSupportCommand(args *model.CommandArgs) *model.CommandResponse {
// 	p.API.CreatePost(&model.Post{
// 		ChannelId: args.ChannelId,
// 		UserId:    p.botID,
// 		Message:   "bruhhhhh, what's up?",
// 	})
// 	return &model.CommandResponse{}
// }

// func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
// 	if post.UserId == p.botID || post.RootId != "" || len(post.FileIds) == 0 {
// 		return
// 	}

// 	supportPackets, names := PostContainsSupportPackage(p, post)

// 	if len(supportPackets) == 0 {
// 		return
// 	}

// 	p.API.CreatePost(&model.Post{
// 		ChannelId: post.ChannelId,
// 		RootId:    post.Id,
// 		UserId:    p.botID,
// 		Message:   "Uploading support packet for " + strings.Join(names, " ,"),
// 	})

// 	ProcessSupportPackets(p, supportPackets, post)
// }
