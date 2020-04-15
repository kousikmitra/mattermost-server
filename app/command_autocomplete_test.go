// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
)

func TestParseStaticListArgument(t *testing.T) {
	items := []model.AutocompleteListItem{
		{
			Hint:     "[hint]",
			Item:     "on",
			HelpText: "help",
		},
	}
	fixedArgs := &model.AutocompleteStaticListArg{PossibleArguments: items}

	argument := &model.AutocompleteArg{
		Name:     "", //positional
		HelpText: "some_help",
		Type:     model.AutocompleteArgTypeStaticList,
		Data:     fixedArgs,
	}
	found, _, _, suggestions := parseStaticListArgument(argument, "", "") //TODO understand this!
	assert.True(t, found)
	assert.Equal(t, []model.AutocompleteSuggestion{{Complete: "on", Suggestion: "on", Hint: "[hint]", Description: "help"}}, suggestions)

	found, _, _, suggestions = parseStaticListArgument(argument, "", "o")
	assert.True(t, found)
	assert.Equal(t, []model.AutocompleteSuggestion{{Complete: "on", Suggestion: "on", Hint: "[hint]", Description: "help"}}, suggestions)

	found, parsed, toBeParsed, _ := parseStaticListArgument(argument, "", "on ")
	assert.False(t, found)
	assert.Equal(t, "on ", parsed)
	assert.Equal(t, "", toBeParsed)

	found, parsed, toBeParsed, _ = parseStaticListArgument(argument, "", "on some")
	assert.False(t, found)
	assert.Equal(t, "on ", parsed)
	assert.Equal(t, "some", toBeParsed)

	fixedArgs.PossibleArguments = append(fixedArgs.PossibleArguments,
		model.AutocompleteListItem{Hint: "[hint]", Item: "off", HelpText: "help"})

	found, _, _, suggestions = parseStaticListArgument(argument, "", "o")
	assert.True(t, found)
	assert.Equal(t, []model.AutocompleteSuggestion{{Complete: "on", Suggestion: "on", Hint: "[hint]", Description: "help"}, {Complete: "off", Suggestion: "off", Hint: "[hint]", Description: "help"}}, suggestions)

	found, _, _, suggestions = parseStaticListArgument(argument, "", "of")
	assert.True(t, found)
	assert.Equal(t, []model.AutocompleteSuggestion{{Complete: "off", Suggestion: "off", Hint: "[hint]", Description: "help"}}, suggestions)

	found, _, _, suggestions = parseStaticListArgument(argument, "", "o some")
	assert.True(t, found)
	assert.Len(t, suggestions, 0)

	found, parsed, toBeParsed, _ = parseStaticListArgument(argument, "", "off some")
	assert.False(t, found)
	assert.Equal(t, "off ", parsed)
	assert.Equal(t, "some", toBeParsed)

	fixedArgs.PossibleArguments = append(fixedArgs.PossibleArguments,
		model.AutocompleteListItem{Hint: "[hint]", Item: "onon", HelpText: "help"})

	found, _, _, suggestions = parseStaticListArgument(argument, "", "on")
	assert.True(t, found)
	assert.Equal(t, []model.AutocompleteSuggestion{{Complete: "on", Suggestion: "on", Hint: "[hint]", Description: "help"}, {Complete: "onon", Suggestion: "onon", Hint: "[hint]", Description: "help"}}, suggestions)

	found, _, _, suggestions = parseStaticListArgument(argument, "bla ", "ono")
	assert.True(t, found)
	assert.Equal(t, []model.AutocompleteSuggestion{{Complete: "bla onon", Suggestion: "onon", Hint: "[hint]", Description: "help"}}, suggestions)

	found, parsed, toBeParsed, _ = parseStaticListArgument(argument, "", "on some")
	assert.False(t, found)
	assert.Equal(t, "on ", parsed)
	assert.Equal(t, "some", toBeParsed)

	found, parsed, toBeParsed, _ = parseStaticListArgument(argument, "", "onon some")
	assert.False(t, found)
	assert.Equal(t, "onon ", parsed)
	assert.Equal(t, "some", toBeParsed)
}

func TestParseInputTextArgument(t *testing.T) {
	argument := &model.AutocompleteArg{
		Name:     "", //positional
		HelpText: "some_help",
		Type:     model.AutocompleteArgTypeText,
		Data:     &model.AutocompleteTextArg{Hint: "hint", Pattern: "pat"},
	}

	found, _, _, suggestion := parseInputTextArgument(argument, "", "")
	assert.True(t, found)
	assert.Equal(t, model.AutocompleteSuggestion{Complete: "", Suggestion: "", Hint: "hint", Description: "some_help"}, suggestion)

	found, _, _, suggestion = parseInputTextArgument(argument, "", " ")
	assert.True(t, found)
	assert.Equal(t, model.AutocompleteSuggestion{Complete: " ", Suggestion: "", Hint: "hint", Description: "some_help"}, suggestion)

	found, _, _, suggestion = parseInputTextArgument(argument, "", "abc")
	assert.True(t, found)
	assert.Equal(t, model.AutocompleteSuggestion{Complete: "abc", Suggestion: "", Hint: "hint", Description: "some_help"}, suggestion)

	found, _, _, suggestion = parseInputTextArgument(argument, "", "\"abc dfd df ")
	assert.True(t, found)
	assert.Equal(t, model.AutocompleteSuggestion{Complete: "\"abc dfd df ", Suggestion: "", Hint: "hint", Description: "some_help"}, suggestion)

	found, parsed, toBeParsed, _ := parseInputTextArgument(argument, "", "abc efg ")
	assert.False(t, found)
	assert.Equal(t, "abc ", parsed)
	assert.Equal(t, "efg ", toBeParsed)

	found, parsed, toBeParsed, _ = parseInputTextArgument(argument, "", "abc ")
	assert.False(t, found)
	assert.Equal(t, "abc ", parsed)
	assert.Equal(t, "", toBeParsed)

	found, parsed, toBeParsed, _ = parseInputTextArgument(argument, "", "\"abc def\" abc")
	assert.False(t, found)
	assert.Equal(t, "\"abc def\" ", parsed)
	assert.Equal(t, "abc", toBeParsed)

	found, parsed, toBeParsed, _ = parseInputTextArgument(argument, "", "\"abc def\"")
	assert.False(t, found)
	assert.Equal(t, "\"abc def\"", parsed)
	assert.Equal(t, "", toBeParsed)
}

func TestParseNamedArguments(t *testing.T) {
	argument := &model.AutocompleteArg{
		Name:     "name", //named
		HelpText: "some_help",
		Type:     model.AutocompleteArgTypeText,
		Data:     &model.AutocompleteTextArg{Hint: "hint", Pattern: "pat"},
	}

	found, _, _, suggestion := parseNamedArgument(argument, "", "")
	assert.True(t, found)
	assert.Equal(t, model.AutocompleteSuggestion{Complete: "--name ", Suggestion: "--name", Hint: "hint", Description: "some_help"}, suggestion)

	found, _, _, suggestion = parseNamedArgument(argument, "", " ")
	assert.True(t, found)
	assert.Equal(t, model.AutocompleteSuggestion{Complete: " --name ", Suggestion: "--name", Hint: "hint", Description: "some_help"}, suggestion)

	found, parsed, toBeParsed, _ := parseNamedArgument(argument, "", "abc")
	assert.False(t, found)
	assert.Equal(t, "abc", parsed)
	assert.Equal(t, "", toBeParsed)

	found, parsed, toBeParsed, suggestion = parseNamedArgument(argument, "", "-")
	assert.True(t, found)
	assert.Equal(t, "-", parsed)
	assert.Equal(t, "", toBeParsed)
	assert.Equal(t, model.AutocompleteSuggestion{Complete: "--name ", Suggestion: "--name", Hint: "hint", Description: "some_help"}, suggestion)

	found, parsed, toBeParsed, suggestion = parseNamedArgument(argument, "", " -")
	assert.True(t, found)
	assert.Equal(t, " -", parsed)
	assert.Equal(t, "", toBeParsed)
	assert.Equal(t, model.AutocompleteSuggestion{Complete: " --name ", Suggestion: "--name", Hint: "hint", Description: "some_help"}, suggestion)

	found, parsed, toBeParsed, suggestion = parseNamedArgument(argument, "", "--name")
	assert.True(t, found)
	assert.Equal(t, "--name", parsed)
	assert.Equal(t, "", toBeParsed)
	assert.Equal(t, model.AutocompleteSuggestion{Complete: "--name ", Suggestion: "--name", Hint: "hint", Description: "some_help"}, suggestion)

	found, parsed, toBeParsed, suggestion = parseNamedArgument(argument, "", "--name bla")
	assert.False(t, found)
	assert.Equal(t, "--name ", parsed)
	assert.Equal(t, "bla", toBeParsed)

	found, parsed, toBeParsed, _ = parseNamedArgument(argument, "", "--name bla gla")
	assert.False(t, found)
	assert.Equal(t, "--name ", parsed)
	assert.Equal(t, "bla gla", toBeParsed)

	found, parsed, toBeParsed, _ = parseNamedArgument(argument, "", "--name \"bla gla\"")
	assert.False(t, found)
	assert.Equal(t, "--name ", parsed)
	assert.Equal(t, "\"bla gla\"", toBeParsed)

	found, parsed, toBeParsed, _ = parseNamedArgument(argument, "", "--name \"bla gla\" ")
	assert.False(t, found)
	assert.Equal(t, "--name ", parsed)
	assert.Equal(t, "\"bla gla\" ", toBeParsed)

	found, parsed, toBeParsed, _ = parseNamedArgument(argument, "", "bla")
	assert.False(t, found)
	assert.Equal(t, "bla", parsed)
	assert.Equal(t, "", toBeParsed)

}

func TestSuggestions(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	jira := createJiraAutocompleteData()

	suggestions := th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "ji", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, jira.Trigger, suggestions[0].Complete)
	assert.Equal(t, jira.Trigger, suggestions[0].Suggestion)
	assert.Equal(t, "[command]", suggestions[0].Hint)
	assert.Equal(t, jira.HelpText, suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira crea", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "jira create", suggestions[0].Complete)
	assert.Equal(t, "create", suggestions[0].Suggestion)
	assert.Equal(t, "[issue text]", suggestions[0].Hint)
	assert.Equal(t, "Create a new Issue", suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira c", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 2)
	assert.Equal(t, "jira create", suggestions[1].Complete)
	assert.Equal(t, "create", suggestions[1].Suggestion)
	assert.Equal(t, "[issue text]", suggestions[1].Hint)
	assert.Equal(t, "Create a new Issue", suggestions[1].Description)
	assert.Equal(t, "jira connect", suggestions[0].Complete)
	assert.Equal(t, "connect", suggestions[0].Suggestion)
	assert.Equal(t, "[url]", suggestions[0].Hint)
	assert.Equal(t, "Connect your Mattermost account to your Jira account", suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira create ", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "jira create ", suggestions[0].Complete)
	assert.Equal(t, "", suggestions[0].Suggestion)
	assert.Equal(t, "[text]", suggestions[0].Hint)
	assert.Equal(t, "This text is optional, will be inserted into the description field", suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira create some", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "jira create some", suggestions[0].Complete)
	assert.Equal(t, "", suggestions[0].Suggestion)
	assert.Equal(t, "[text]", suggestions[0].Hint)
	assert.Equal(t, "This text is optional, will be inserted into the description field", suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira create some text ", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 0)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "invalid command", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 0)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira settings notifications o", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 2)
	assert.Equal(t, "jira settings notifications on", suggestions[0].Complete)
	assert.Equal(t, "on", suggestions[0].Suggestion)
	assert.Equal(t, "Turn notifications on", suggestions[0].Hint)
	assert.Equal(t, "", suggestions[0].Description)
	assert.Equal(t, "jira settings notifications off", suggestions[1].Complete)
	assert.Equal(t, "off", suggestions[1].Suggestion)
	assert.Equal(t, "Turn notifications off", suggestions[1].Hint)
	assert.Equal(t, "", suggestions[1].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira ", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 11)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira ", model.SYSTEM_USER_ROLE_ID)
	assert.Len(t, suggestions, 9)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira create \"some issue text", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "jira create \"some issue text", suggestions[0].Complete)
	assert.Equal(t, "", suggestions[0].Suggestion)
	assert.Equal(t, "[text]", suggestions[0].Hint)
	assert.Equal(t, "This text is optional, will be inserted into the description field", suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira timezone ", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "jira timezone --zone ", suggestions[0].Complete)
	assert.Equal(t, "--zone", suggestions[0].Suggestion)
	assert.Equal(t, "[UTC+07:00]", suggestions[0].Hint)
	assert.Equal(t, "Set timezone", suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira timezone --", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "jira timezone --zone ", suggestions[0].Complete)
	assert.Equal(t, "--zone", suggestions[0].Suggestion)
	assert.Equal(t, "[UTC+07:00]", suggestions[0].Hint)
	assert.Equal(t, "Set timezone", suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira timezone --zone ", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "jira timezone --zone ", suggestions[0].Complete)
	assert.Equal(t, "", suggestions[0].Suggestion)
	assert.Equal(t, "[UTC+07:00]", suggestions[0].Hint)
	assert.Equal(t, "Set timezone", suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira timezone --zone bla", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "jira timezone --zone bla", suggestions[0].Complete)
	assert.Equal(t, "", suggestions[0].Suggestion)
	assert.Equal(t, "[UTC+07:00]", suggestions[0].Hint)
	assert.Equal(t, "Set timezone", suggestions[0].Description)

	suggestions = th.App.getSuggestions([]*model.AutocompleteData{jira}, "", "jira timezone bla", model.SYSTEM_ADMIN_ROLE_ID)
	assert.Len(t, suggestions, 0)

}

// createJiraAutocompleteData will create autocomplete data for jira plugin. For testing purposes only.
func createJiraAutocompleteData() *model.AutocompleteData {
	jira := model.NewAutocompleteData("jira", "[command]", "Available commands: connect, assign, disconnect, create, transition, view, subscribe, settings, install cloud/server, uninstall cloud/server, help")

	connect := model.NewAutocompleteData("connect", "[url]", "Connect your Mattermost account to your Jira account")
	jira.AddCommand(connect)

	disconnect := model.NewAutocompleteData("disconnect", "", "Disconnect your Mattermost account from your Jira account")
	jira.AddCommand(disconnect)

	assign := model.NewAutocompleteData("assign", "[issue]", "Change the assignee of a Jira issue")
	assign.AddDynamicListArgument("List of issues is downloading from your Jira account", "/url/issue-key")
	assign.AddDynamicListArgument("List of assignees is downloading from your Jira account", "/url/assignee")
	jira.AddCommand(assign)

	create := model.NewAutocompleteData("create", "[issue text]", "Create a new Issue")
	create.AddTextArgument("This text is optional, will be inserted into the description field", "[text]", "")
	jira.AddCommand(create)

	transition := model.NewAutocompleteData("transition", "[issue]", "Change the state of a Jira issue")
	assign.AddDynamicListArgument("List of issues is downloading from your Jira account", "/url/issue-key")
	assign.AddDynamicListArgument("List of states is downloading from your Jira account", "/url/states")
	jira.AddCommand(transition)

	subscribe := model.NewAutocompleteData("subscribe", "", "Configure the Jira notifications sent to this channel")
	jira.AddCommand(subscribe)

	view := model.NewAutocompleteData("view", "[issue]", "View the details of a specific Jira issue")
	assign.AddDynamicListArgument("List of issues is downloading from your Jira account", "/url/issue-key")
	jira.AddCommand(view)

	settings := model.NewAutocompleteData("settings", "", "Update your user settings")
	notifications := model.NewAutocompleteData("notifications", "[on/off]", "Turn notifications on or off")

	items := []model.AutocompleteListItem{
		{
			Hint: "Turn notifications on",
			Item: "on",
		},
		{
			Hint: "Turn notifications off",
			Item: "off",
		},
	}
	notifications.AddStaticListArgument("Turn notifications on or off", items)
	settings.AddCommand(notifications)
	jira.AddCommand(settings)

	timezone := model.NewAutocompleteData("timezone", "", "Update your timezone")
	timezone.AddNamedTextArgument("zone", "Set timezone", "[UTC+07:00]", "")
	jira.AddCommand(timezone)

	install := model.NewAutocompleteData("install", "", "Connect Mattermost to a Jira instance")
	install.RoleID = model.SYSTEM_ADMIN_ROLE_ID
	cloud := model.NewAutocompleteData("cloud", "", "Connect to a Jira Cloud instance")
	urlPattern := "https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"
	cloud.AddTextArgument("input URL of the Jira Cloud instance", "[URL]", urlPattern)
	install.AddCommand(cloud)
	server := model.NewAutocompleteData("server", "", "Connect to a Jira Server or Data Center instance")
	server.AddTextArgument("input URL of the Jira Server or Data Center instance", "[URL]", urlPattern)
	install.AddCommand(server)
	jira.AddCommand(install)

	uninstall := model.NewAutocompleteData("uninstall", "", "Disconnect Mattermost from a Jira instance")
	uninstall.RoleID = model.SYSTEM_ADMIN_ROLE_ID
	cloud = model.NewAutocompleteData("cloud", "", "Disconnect from a Jira Cloud instance")
	cloud.AddTextArgument("input URL of the Jira Cloud instance", "[URL]", urlPattern)
	uninstall.AddCommand(cloud)
	server = model.NewAutocompleteData("server", "", "Disconnect from a Jira Server or Data Center instance")
	server.AddTextArgument("input URL of the Jira Server or Data Center instance", "[URL]", urlPattern)
	uninstall.AddCommand(server)
	jira.AddCommand(uninstall)

	return jira
}
