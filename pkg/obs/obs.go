package obs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/config"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/typedefs"
)

type Config struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	RecordingDir string `yaml:"recording_dir"`
}

type OBS struct {
	Client               *goobs.Client
	TaskSourceName       string
	BackgroundSourceName string
	AvatarSourceName     string
}

type Task struct {
	Text       string      `json:"text"`
	PosX       float64     `json:"pos_x"`
	PosY       float64     `json:"pos_y"`
	Width      float64     `json:"width"`
	Height     float64     `json:"height"`
	Color      Color       `json:"color"`
	Background *Background `json:"background,omitempty"`
}

type Color struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

type Background struct {
	Color Color `json:"color"`
}

const (
	SourceTextType                 string = "text_ft2_source_v2"
	SourceColorBlockType           string = "color_source_v3"
	ProfileParameterOutputFileName string = "FilenameFormatting"
)

func New(client *goobs.Client, task_name, background_name, avatar_name string) *OBS {
	return &OBS{
		Client:               client,
		TaskSourceName:       task_name,
		BackgroundSourceName: background_name,
		AvatarSourceName:     avatar_name,
	}
}

func convertColor(c Color) (int64, error) {
	r := c.R
	g := c.G
	b := c.B
	a := c.A
	hexNumber := fmt.Sprintf("%02x%02x%02x%02x", a, b, g, r)
	return strconv.ParseInt(hexNumber, 16, 64)
}

func (obs *OBS) ConvertIntToHex(c int64) (*string, error) {
	h := fmt.Sprintf("%x", c)
	if len(h) != 8 {
		return nil, errors.New("invalid color integer conversion")
	}
	return &h, nil
}

func (obs *OBS) RefreshSources() error {
	task_exists := true
	background_exists := true
	avatar_exists := true
	_, err := obs.GetInputSettings(obs.TaskSourceName)
	if err != nil {
		task_exists = false
	}
	_, err = obs.GetInputSettings(obs.BackgroundSourceName)
	if err != nil {
		background_exists = false
	}
	_, err = obs.GetInputSettings(obs.AvatarSourceName)
	if err != nil {
		avatar_exists = false
	}
	current_scene, err := obs.GetCurrentScene()
	if err != nil {
		return err
	}
	if task_exists {
		item_id := obs.GetSceneItemId(current_scene, obs.TaskSourceName)
		is_visible, err := obs.GetSceneSourceVisible(item_id, current_scene)
		if err != nil {
			return err
		}
		if *is_visible {
			err = obs.SetSceneSourceVisible(item_id, current_scene, false)
			if err != nil {
				return err
			}
			err = obs.SetSceneSourceVisible(item_id, current_scene, true)
			if err != nil {
				return err
			}
		}
	}
	if background_exists {
		item_id := obs.GetSceneItemId(current_scene, obs.BackgroundSourceName)
		is_visible, err := obs.GetSceneSourceVisible(item_id, current_scene)
		if err != nil {
			return err
		}
		if *is_visible {
			err = obs.SetSceneSourceVisible(item_id, current_scene, false)
			if err != nil {
				return err
			}
			err = obs.SetSceneSourceVisible(item_id, current_scene, true)
			if err != nil {
				return err
			}
		}
	}
	if avatar_exists {
		item_id := obs.GetSceneItemId(current_scene, obs.AvatarSourceName)
		is_visible, err := obs.GetSceneSourceVisible(item_id, current_scene)
		if err != nil {
			return err
		}
		if *is_visible {
			err = obs.SetSceneSourceVisible(item_id, current_scene, false)
			if err != nil {
				return err
			}
			err = obs.SetSceneSourceVisible(item_id, current_scene, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (obs *OBS) ConvertIntToColor(c int64) (*Color, error) {
	hex, err := obs.ConvertIntToHex(c)
	if err != nil {
		return nil, errors.New("cannot convert Color [a] value to uint")
	}
	h := *hex
	aHex := h[0:2]
	bHex := h[2:4]
	gHex := h[4:6]
	rHex := h[6:8]
	aU8, err := strconv.ParseUint(aHex, 16, 64)
	if err != nil {
		return nil, errors.New("cannot convert Color [a] value to uint")
	}
	bU8, err := strconv.ParseUint(bHex, 16, 64)
	if err != nil {
		return nil, errors.New("cannot convert Color [b] value to uint")
	}
	gU8, err := strconv.ParseUint(gHex, 16, 64)
	if err != nil {
		return nil, errors.New("cannot convert Color [g] value to uint")
	}
	rU8, err := strconv.ParseUint(rHex, 16, 64)
	if err != nil {
		return nil, errors.New("cannot convert Color [r] value to uint")
	}
	a := uint8(aU8)
	b := uint8(bU8)
	g := uint8(gU8)
	r := uint8(rU8)
	return &Color{
		A: a,
		B: b,
		G: g,
		R: r,
	}, nil
}

func (obs *OBS) SetTask(task Task) error {
	background_exists := true
	task_exists := true
	_, err := obs.GetInputSettings(obs.BackgroundSourceName)
	if err != nil {
		background_exists = false
	}
	_, err = obs.GetInputSettings(obs.TaskSourceName)
	if err != nil {
		task_exists = false
	}
	color, err := convertColor(task.Color)
	if err != nil {
		return errors.New("Error converting color: " + err.Error())
	}
	task_settings := map[string]interface{}{
		"text":   task.Text,
		"color1": color,
		"color2": color,
	}
	current_scene, err := obs.GetCurrentScene()
	if err != nil {
		return errors.New("Cannot get current scene in SetTask: " + err.Error())
	}
	if task.Background != nil {
		background_color, err := convertColor(task.Background.Color)
		if err != nil {
			return errors.New("Error converting background color: " + err.Error())
		}
		background_settings := map[string]interface{}{
			"color":  background_color,
			"width":  task.Width + 4,
			"height": task.Height + 4,
		}
		if background_exists {
			_, err = obs.SetInputSettings(obs.BackgroundSourceName, background_settings)
			if err != nil {
				return errors.New("Cannot set background settings in SetTask: " + err.Error())
			}
			_, err = obs.SetSceneItemTransform(obs.GetSceneItemId(current_scene, obs.BackgroundSourceName), current_scene, task.PosX-2, task.PosY-2, task.Width+4, task.Height+4)
			if err != nil {
				return errors.New(err.Error())
			}
		} else {
			_, err = obs.CreateInput(SourceColorBlockType, current_scene, obs.BackgroundSourceName, true, background_settings)
			if err != nil {
				return errors.New(err.Error())
			}
			task_id := obs.GetSceneItemId(current_scene, obs.TaskSourceName)
			if task_id > 0 {
				resp, err := obs.GetSceneItemIndex(obs.GetSceneItemId(current_scene, obs.BackgroundSourceName), current_scene)
				if err != nil {
					return errors.New(err.Error())
				}
				_, err = obs.SetSceneItemIndex(task_id, resp.SceneItemIndex, current_scene)
				if err != nil {
					return errors.New(err.Error())
				}
			}
			_, err = obs.SetSceneItemTransform(obs.GetSceneItemId(current_scene, obs.BackgroundSourceName), current_scene, task.PosX-2, task.PosY-2, task.Width+4, task.Height+4)
			if err != nil {
				return errors.New(err.Error())
			}
		}
	} else {
		if background_exists {
			_, err = obs.RemoveSceneItem(obs.GetSceneItemId(current_scene, obs.BackgroundSourceName), current_scene)
			if err != nil {
				return errors.New("Error Removing background scene item: " + err.Error())
			}
		}
	}

	if len(task.Text) > 0 {
		if task_exists {
			_, err := obs.SetInputSettings(obs.TaskSourceName, task_settings)
			if err != nil {
				return errors.New(err.Error())
			}

		} else {
			_, err = obs.CreateInput(SourceTextType, current_scene, obs.TaskSourceName, true, task_settings)
			if err != nil {
				return errors.New(err.Error())
			}
		}
		_, err = obs.SetSceneItemTransform(obs.GetSceneItemId(current_scene, obs.TaskSourceName), current_scene, task.PosX, task.PosY, task.Width, task.Height)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}

func (obs *OBS) GetInputSettings(name string) (*inputs.GetInputSettingsResponse, error) {
	return obs.Client.Inputs.GetInputSettings(&inputs.GetInputSettingsParams{
		InputName: name,
	})
}

func (obs *OBS) SetInputSettings(name string, settings map[string]interface{}) (*inputs.SetInputSettingsResponse, error) {
	return obs.Client.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{
		InputName:     name,
		InputSettings: settings,
	})
}

func (obs *OBS) GetSceneItemId(scene string, source string) float64 {
	resp, err := obs.Client.SceneItems.GetSceneItemId(&sceneitems.GetSceneItemIdParams{
		SceneName:  scene,
		SourceName: source,
	})
	if err != nil {
		return -1
	}
	return resp.SceneItemId
}

func (obs *OBS) GetSceneItemTransform(item_id float64, name string) (*sceneitems.GetSceneItemTransformResponse, error) {
	return obs.Client.SceneItems.GetSceneItemTransform(&sceneitems.GetSceneItemTransformParams{
		SceneItemId: item_id,
		SceneName:   name,
	})
}

func (obs *OBS) RemoveSceneItem(item_id float64, name string) (*sceneitems.RemoveSceneItemResponse, error) {
	return obs.Client.SceneItems.RemoveSceneItem(&sceneitems.RemoveSceneItemParams{
		SceneItemId: item_id,
		SceneName:   name,
	})
}

func (obs *OBS) CreateInput(kind string, scene string, name string, enabled bool, settings map[string]interface{}) (*inputs.CreateInputResponse, error) {
	resp, err := obs.Client.Inputs.CreateInput(&inputs.CreateInputParams{
		InputKind:        kind,
		SceneName:        scene,
		InputName:        name,
		SceneItemEnabled: &enabled,
		InputSettings:    settings,
	})
	if err != nil && strings.Contains(err.Error(), "601") { // resource already exists
		err = nil
	}
	return resp, err
}

func (obs *OBS) CreateSceneItem(name string) (*sceneitems.CreateSceneItemResponse, error) {
	enabled := true
	return obs.Client.SceneItems.CreateSceneItem(&sceneitems.CreateSceneItemParams{
		SceneItemEnabled: &enabled,
		SceneName:        "Main",
		SourceName:       name,
	})
}

func (obs *OBS) SetSceneSourceVisible(item_id float64, name string, visible bool) error {
	_, err := obs.Client.SceneItems.SetSceneItemEnabled(&sceneitems.SetSceneItemEnabledParams{
		SceneItemEnabled: &visible,
		SceneItemId:      item_id,
		SceneName:        name,
	})
	return err
}

func (obs *OBS) GetSceneSourceVisible(item_id float64, name string) (*bool, error) {
	settings, err := obs.Client.SceneItems.GetSceneItemEnabled(&sceneitems.GetSceneItemEnabledParams{
		SceneItemId: item_id,
		SceneName:   name,
	})
	if settings == nil {
		return nil, err
	}
	return &settings.SceneItemEnabled, err
}

func (obs *OBS) SetSceneItemTransform(item_id float64, name string, posX, posY, width, height float64) (*sceneitems.SetSceneItemTransformResponse, error) {
	return obs.Client.SceneItems.SetSceneItemTransform(&sceneitems.SetSceneItemTransformParams{
		SceneItemId: item_id,
		SceneName:   name,
		SceneItemTransform: &typedefs.SceneItemTransform{
			PositionX:    posX,
			PositionY:    posY,
			Alignment:    5, // Top Left
			BoundsWidth:  width,
			BoundsHeight: height,
			ScaleX:       1,
			ScaleY:       1,
			BoundsType:   "OBS_BOUNDS_MAX_ONLY",
		},
	})
}

func (obs *OBS) GetSceneItemIndex(item_id float64, name string) (*sceneitems.GetSceneItemIndexResponse, error) {
	return obs.Client.SceneItems.GetSceneItemIndex(&sceneitems.GetSceneItemIndexParams{
		SceneItemId: item_id,
		SceneName:   name,
	})
}

func (obs *OBS) SetSceneItemIndex(item_id float64, item_index float64, name string) (*sceneitems.SetSceneItemIndexResponse, error) {
	return obs.Client.SceneItems.SetSceneItemIndex(&sceneitems.SetSceneItemIndexParams{
		SceneItemId:    item_id,
		SceneName:      name,
		SceneItemIndex: item_index,
	})
}

func (obs *OBS) CreateScene(name string) (*scenes.CreateSceneResponse, error) {
	if len(name) == 0 {
		return nil, errors.New("Scene name length must be greater than 0")
	}
	return obs.Client.Scenes.CreateScene(&scenes.CreateSceneParams{
		SceneName: name,
	})
}

func (obs *OBS) GetSceneList() (*scenes.GetSceneListResponse, error) {
	return obs.Client.Scenes.GetSceneList()
}

func (obs *OBS) GetStreamStatus() (bool, error) {
	status, err := obs.Client.Stream.GetStreamStatus()
	return status.OutputActive, err
}

func (obs *OBS) ToggleStream() (bool, error) {
	status, err := obs.Client.Stream.ToggleStream()
	return status.OutputActive, err
}

func (obs *OBS) GetRecordStatus() (bool, error) {
	status, err := obs.Client.Record.GetRecordStatus()
	return status.OutputActive, err
}

func (obs *OBS) ToggleRecord() error {
	_, err := obs.Client.Record.ToggleRecord()
	return err
}

func (obs *OBS) GetCurrentScene() (string, error) {
	resp, err := obs.Client.Scenes.GetCurrentProgramScene()
	if err != nil {
		return "", err
	}
	if resp.CurrentProgramSceneName == "" {
		return "", errors.New("no current scene detected")
	}
	return resp.CurrentProgramSceneName, nil
}

func (obs *OBS) GetRecordDirectory() (string, error) {
	resp, err := obs.Client.Config.GetRecordDirectory()
	if err != nil {
		return "", errors.New("Cannot get current recording directory: " + err.Error())
	}
	if resp.RecordDirectory == "" {
		return "", errors.New("Recording directory was empty")
	}
	return resp.RecordDirectory, nil
}

func (obs *OBS) GetProfileParameter(parameter string) (string, error) {
	resp, err := obs.Client.Config.GetProfileParameter(&config.GetProfileParameterParams{
		ParameterCategory: "Output",
		ParameterName:     parameter,
	})
	if err != nil {
		return "", errors.New("Could not get profile parameter [" + parameter + "]: " + err.Error())
	}
	if resp.ParameterValue == "" {
		return "", errors.New("profile parameter [" + parameter + "] was empty")
	}
	return resp.ParameterValue, nil
}

func (obs *OBS) SetProfileParameter(parameter string, value string) error {
	_, err := obs.Client.Config.SetProfileParameter(&config.SetProfileParameterParams{
		ParameterCategory: "Output",
		ParameterName:     parameter,
		ParameterValue:    value,
	})
	if err != nil {
		return errors.New("Could not set profile parameter [" + parameter + "]: " + err.Error())
	}
	return nil
}
