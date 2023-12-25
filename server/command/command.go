package command

// import (
// 	"strings"

// 	"github.com/pkg/errors"

// 	"github.com/coltoneshaw/mattermost-plugin-customers/server/bot"
// 	"github.com/coltoneshaw/mattermost-plugin-customers/server/config"
// 	"github.com/mattermost/mattermost/server/public/model"
// 	"github.com/mattermost/mattermost/server/public/plugin"
// 	pluginapi "github.com/mattermost/mattermost/server/public/pluginapi"
// )

// const helpText = "###### Customer Info Plugin - Slash Command Help\n" +
// 	"* `/customer ticket [zendesk #]` - Update data for a ticket \n" +
// 	"\n"

// const confirmPrompt = "CONFIRM"

// const availableCommands = "Available commands: ticket"

// // Register is a function that allows the runner to register commands with the mattermost server.
// type Register func(*model.Command) error

// // RegisterCommands should be called by the plugin to register all necessary commands
// func RegisterCommands(registerFunc Register) error {
// 	return registerFunc(getCommand())
// }

// func getCommand() *model.Command {
// 	return &model.Command{
// 		Trigger:          "customer",
// 		DisplayName:      "Customer",
// 		Description:      "Customer",
// 		AutoComplete:     true,
// 		AutoCompleteDesc: "Available commands: ticket",
// 		AutoCompleteHint: "[command]",
// 		AutocompleteData: getAutocompleteData(),
// 	}
// }

// func getAutocompleteData() *model.AutocompleteData {
// 	command := model.NewAutocompleteData("customer", "[command]", availableCommands)

// 	run := model.NewAutocompleteData("ticket", "", "downloads ticket data to be updated")
// 	command.AddCommand(run)

// 	return command
// }

// // Runner handles commands.
// type Runner struct {
// 	context       *plugin.Context
// 	args          *model.CommandArgs
// 	pluginAPI     *pluginapi.Client
// 	poster        bot.Poster
// 	configService config.Service
// }

// // NewCommandRunner creates a command runner.
// func NewCommandRunner(ctx *plugin.Context,
// 	args *model.CommandArgs,
// 	api *pluginapi.Client,
// 	poster bot.Poster,
// 	configService config.Service,
// ) *Runner {
// 	return &Runner{
// 		context:       ctx,
// 		args:          args,
// 		pluginAPI:     api,
// 		poster:        poster,
// 		configService: configService,
// 	}
// }
// func (r *Runner) isValid() error {
// 	if r.context == nil || r.args == nil || r.pluginAPI == nil {
// 		return errors.New("invalid arguments to command.Runner")
// 	}
// 	return nil
// }

// func (r *Runner) postCommandResponse(text string) {
// 	post := &model.Post{
// 		Message: text,
// 	}
// 	r.poster.EphemeralPost(r.args.UserId, r.args.ChannelId, post)
// }

// func (r *Runner) actionRun(args []string) {

// }

// // Execute should be called by the plugin when a command invocation is received from the Mattermost server.
// func (r *Runner) Execute() error {
// 	if err := r.isValid(); err != nil {
// 		return err
// 	}

// 	split := strings.Fields(r.args.Command)
// 	command := split[0]
// 	parameters := []string{}
// 	cmd := ""
// 	if len(split) > 1 {
// 		cmd = split[1]
// 	}
// 	if len(split) > 2 {
// 		parameters = split[2:]
// 	}

// 	if command != "/customer" {
// 		return nil
// 	}

// 	if cmd == "run" {
// 		r.actionRun(parameters)
// 	}

// 	return nil
// }
