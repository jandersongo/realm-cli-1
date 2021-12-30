package app

import (
	"github.com/10gen/realm-cli/internal/cli/user"
	"github.com/10gen/realm-cli/internal/cloud/realm"
	"github.com/10gen/realm-cli/internal/local"
	"github.com/10gen/realm-cli/internal/terminal"
)

type initInputs struct {
	newAppInputs
}

func (i *initInputs) Resolve(profile *user.Profile, ui terminal.UI) error {
	app, appErr := local.LoadAppConfig(profile.WorkingDirectory)
	if appErr != nil {
		return appErr
	}
	if app.RootDir != "" {
		return errProjectExists{}
	}

	if i.RemoteApp == "" {
		if i.Name == "" {
			if err := ui.Input(&i.Name, terminal.AskOptions{Message: "App Name"}); err != nil {
				return err
			}
		}
		if i.DeploymentModel == realm.DeploymentModelEmpty {
			i.DeploymentModel = flagDeploymentModelDefault
		}
		if i.Location == realm.LocationEmpty {
			i.Location = flagLocationDefault
		}
		if i.ConfigVersion == realm.AppConfigVersionZero {
			i.ConfigVersion = realm.DefaultAppConfigVersion
		}
	}

	return nil
}
