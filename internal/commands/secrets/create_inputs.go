package secrets

import (
	"github.com/10gen/realm-cli/internal/cli"
	"github.com/10gen/realm-cli/internal/cli/user"
	"github.com/10gen/realm-cli/internal/terminal"
)

type createInputs struct {
	cli.ProjectInputs
	Name  string
	Value string
}

func (i *createInputs) Resolve(profile *user.Profile, ui terminal.UI) error {
	if err := i.ProjectInputs.Resolve(ui, profile.WorkingDirectory, true); err != nil {
		return err
	}

	if i.Name == "" {
		if err := ui.Input(&i.Name, terminal.AskOptions{Message: "Secret Name"}); err != nil {
			return err
		}
	}

	if i.Value == "" {
		if err := ui.Password(&i.Value, terminal.AskOptions{Message: "Secret Value"}); err != nil {
			return err
		}
	}

	return nil
}
