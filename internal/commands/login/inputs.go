package login

import (
	"github.com/10gen/realm-cli/internal/cli/user"
	"github.com/10gen/realm-cli/internal/terminal"
)

type inputs struct {
	PublicAPIKey  string
	PrivateAPIKey string
}

func (i *inputs) Resolve(profile *user.Profile, ui terminal.UI) error {
	user := profile.Credentials()

	if i.PublicAPIKey == "" {
		if user.PublicAPIKey == "" {
			if err := ui.Input(&i.PublicAPIKey, terminal.AskOptions{Message: "Public API Key", Default: user.PublicAPIKey}); err != nil {
				return err
			}
		} else {
			i.PublicAPIKey = user.PublicAPIKey
		}
	}

	if i.PrivateAPIKey == "" {
		if user.PrivateAPIKey == "" {
			if err := ui.Password(&i.PrivateAPIKey, terminal.AskOptions{Message: "Private API Key"}); err != nil {
				return err
			}
		} else {
			i.PrivateAPIKey = user.PrivateAPIKey
		}
	}

	return nil
}
