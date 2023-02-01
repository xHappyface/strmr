package obs

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
)

type OBS struct {
	Client *goobs.Client
}

func New(client *goobs.Client) *OBS {
	return &OBS{
		Client: client,
	}
}

func (obs *OBS) SetTask(task string) error {
	_, err := obs.Client.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{
		InputName: "status",
		InputSettings: map[string]interface{}{
			"text": task,
		},
	})
	return err
}
