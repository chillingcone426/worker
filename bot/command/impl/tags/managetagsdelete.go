package tags

import (
	"github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/common/sentry"
	"github.com/TicketsBot/worker/bot/command"
	"github.com/TicketsBot/worker/bot/command/registry"
	"github.com/TicketsBot/worker/bot/customisation"
	"github.com/TicketsBot/worker/bot/dbclient"
	"github.com/TicketsBot/worker/i18n"
	"github.com/rxdn/gdl/objects/interaction"
)

type ManageTagsDeleteCommand struct {
}

func (ManageTagsDeleteCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "delete",
		Description:     i18n.HelpTagDelete,
		Type:            interaction.ApplicationCommandTypeChatInput,
		Aliases:         []string{"del", "rm", "remove"},
		PermissionLevel: permission.Support,
		Category:        command.Tags,
		Arguments: command.Arguments(
			command.NewRequiredArgument("id", "ID of the tag to delete", interaction.OptionTypeString, i18n.MessageTagDeleteInvalidArguments),
		),
	}
}

func (c ManageTagsDeleteCommand) GetExecutor() interface{} {
	return c.Execute
}

func (ManageTagsDeleteCommand) Execute(ctx registry.CommandContext, tagId string) {
	// TODO: Causes a race condition, just try to delete
	var found bool
	{
		tag, err := dbclient.Client.Tag.Get(ctx.GuildId(), tagId)
		if err != nil {
			sentry.ErrorWithContext(err, ctx.ToErrorContext())
			ctx.Reject()
			return
		}

		found = tag != ""
	}

	if !found {
		ctx.Reject()
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageTagDeleteDoesNotExist, tagId)
		return
	}

	if err := dbclient.Client.Tag.Delete(ctx.GuildId(), tagId); err == nil {
		ctx.Reply(customisation.Green, i18n.MessageTag, i18n.MessageTagDeleteSuccess, tagId)
	} else {
		ctx.HandleError(err)
	}
}
