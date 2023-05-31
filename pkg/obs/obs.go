package obs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/config"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/typedefs"
	"github.com/jnrprgmr/strmr/pkg/database"
)

type Config struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	RecordingDir string `yaml:"recording_dir"`
}

type OBS struct {
	Client                      *goobs.Client
	ScreenSourceName            string
	TaskSourceName              string
	BackgroundSourceName        string
	AvatarSourceName            string
	OverlayTextSourceName       string
	OverlayBackgroundSourceName string
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
	SourceScreenType               string = "xshm_input"
	SourceTextType                 string = "text_ft2_source_v2"
	SourceColorBlockType           string = "color_source_v3"
	SourceBrowser                  string = "browser_source"
	ProfileParameterOutputFileName string = "FilenameFormatting"
)

func New(client *goobs.Client, screen_name, task_name, background_name, avatar_name, overlay_text_name, overlay_background_name string) *OBS {
	return &OBS{
		Client:                      client,
		ScreenSourceName:            screen_name,
		TaskSourceName:              task_name,
		BackgroundSourceName:        background_name,
		AvatarSourceName:            avatar_name,
		OverlayTextSourceName:       overlay_text_name,
		OverlayBackgroundSourceName: overlay_background_name,
	}
}

func ConvertColor(c Color) (int64, error) {
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

func (obs *OBS) RefreshSources(background_config database.Metadata, task_text string, task_config database.Metadata) error {
	screen_exists := true
	task_exists := true
	background_exists := true
	avatar_exists := true
	overlay_text_exists := true
	overlay_background_exists := true
	_, err := obs.GetInputSettings(obs.ScreenSourceName)
	if err != nil {
		task_exists = false
	}
	_, err = obs.GetInputSettings(obs.TaskSourceName)
	if err != nil {
		screen_exists = false
	}
	_, err = obs.GetInputSettings(obs.BackgroundSourceName)
	if err != nil {
		background_exists = false
	}
	_, err = obs.GetInputSettings(obs.AvatarSourceName)
	if err != nil {
		avatar_exists = false
	}
	_, err = obs.GetInputSettings(obs.OverlayTextSourceName)
	if err != nil {
		overlay_text_exists = false
	}
	_, err = obs.GetInputSettings(obs.OverlayBackgroundSourceName)
	if err != nil {
		overlay_background_exists = false
	}
	current_scene, err := obs.GetCurrentScene()
	if err != nil {
		return err
	}
	// Screen
	screen_settings := map[string]interface{}{
		"advanced": false,
		"screen":   1, // hard coded
	}
	if !screen_exists {
		_, err = obs.CreateInput(SourceScreenType, current_scene, obs.ScreenSourceName, true, screen_settings)
		if err != nil {
			return errors.New(err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	// Background
	vals := strings.Split(background_config.MetadataValue, ",")
	if len(vals) != 5 {
		return errors.New("background metadata config values not as expected")
	}
	color, err := strconv.Atoi(vals[0])
	if err != nil {
		return errors.New("background metadata config color not as expected: " + err.Error())
	}
	width, err := strconv.ParseFloat(vals[1], 64)
	if err != nil {
		return errors.New("background metadata config width not as expected: " + err.Error())
	}
	height, err := strconv.ParseFloat(vals[2], 64)
	if err != nil {
		return errors.New("background metadata config height not as expected: " + err.Error())
	}
	posx, err := strconv.ParseFloat(vals[3], 64)
	if err != nil {
		return errors.New("background metadata config posx not as expected: " + err.Error())
	}
	posy, err := strconv.ParseFloat(vals[4], 64)
	if err != nil {
		return errors.New("background metadata config posy not as expected: " + err.Error())
	}
	background_settings := map[string]interface{}{
		"color":  color,
		"width":  width + 4,
		"height": height + 4,
	}
	if !background_exists {
		_, err = obs.CreateInput(SourceColorBlockType, current_scene, obs.BackgroundSourceName, true, background_settings)
		if err != nil {
			return errors.New(err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	_, err = obs.SetSceneItemTransform(obs.GetSceneItemId(current_scene, obs.BackgroundSourceName), current_scene, posx-2, posy-2, width+4, height+4)
	if err != nil {
		return errors.New("Background: " + err.Error())
	}
	time.Sleep(2 * time.Second)
	// Task
	vals = strings.Split(task_config.MetadataValue, ",")
	if len(vals) != 5 {
		return errors.New("task metadata config values not as expected")
	}
	color, err = strconv.Atoi(vals[0])
	if err != nil {
		return errors.New("task metadata config color not as expected: " + err.Error())
	}
	width, err = strconv.ParseFloat(vals[1], 64)
	if err != nil {
		return errors.New("task metadata config width not as expected: " + err.Error())
	}
	height, err = strconv.ParseFloat(vals[2], 64)
	if err != nil {
		return errors.New("task metadata config height not as expected: " + err.Error())
	}
	posx, err = strconv.ParseFloat(vals[3], 64)
	if err != nil {
		return errors.New("task metadata config posx not as expected: " + err.Error())
	}
	posy, err = strconv.ParseFloat(vals[4], 64)
	if err != nil {
		return errors.New("task metadata config posy not as expected: " + err.Error())
	}
	task_settings := map[string]interface{}{
		"text":   task_text,
		"color1": color,
		"color2": color,
	}
	if !task_exists {
		_, err = obs.CreateInput(SourceTextType, current_scene, obs.TaskSourceName, true, task_settings)
		if err != nil {
			return errors.New(err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	_, err = obs.SetSceneItemTransform(obs.GetSceneItemId(current_scene, obs.TaskSourceName), current_scene, posx, posy, width, height)
	if err != nil {
		return errors.New("Task " + err.Error())
	}
	time.Sleep(2 * time.Second)
	// Avatar
	avatar_settings := map[string]interface{}{
		"url":                 "http://localhost:8080/avatar",
		"width":               400,
		"height":              400,
		"restart_when_active": true,
	}
	if !avatar_exists {
		_, err = obs.CreateInput(SourceBrowser, current_scene, obs.AvatarSourceName, true, avatar_settings)
		if err != nil {
			return errors.New(err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	_, err = obs.SetSceneItemTransform(obs.GetSceneItemId(current_scene, obs.AvatarSourceName), current_scene, 1160, 475, 400, 400)
	if err != nil {
		return errors.New("Avatar: " + err.Error())
	}
	time.Sleep(2 * time.Second)
	// Overlay background
	color, err = strconv.Atoi("4291297280")
	if err != nil {
		return errors.New("overlay background metadata config color not as expected: " + err.Error())
	}
	width, err = strconv.ParseFloat("1600", 64)
	if err != nil {
		return errors.New("overlay background metadata config width not as expected: " + err.Error())
	}
	height, err = strconv.ParseFloat("900", 64)
	if err != nil {
		return errors.New("overlay background metadata config height not as expected: " + err.Error())
	}
	posx, err = strconv.ParseFloat("0", 64)
	if err != nil {
		return errors.New("overlay background metadata config posx not as expected: " + err.Error())
	}
	posy, err = strconv.ParseFloat("0", 64)
	if err != nil {
		return errors.New("overlay background metadata config posy not as expected: " + err.Error())
	}
	overlay_background_settings := map[string]interface{}{
		"color":  color,
		"width":  width,
		"height": height,
	}
	if !overlay_background_exists {
		_, err = obs.CreateInput(SourceColorBlockType, current_scene, obs.OverlayBackgroundSourceName, true, overlay_background_settings)
		if err != nil {
			return errors.New(err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	_, err = obs.SetSceneItemTransform(obs.GetSceneItemId(current_scene, obs.OverlayBackgroundSourceName), current_scene, posx, posy, width, height)
	if err != nil {
		return errors.New("Overlay Background: " + err.Error())
	}
	time.Sleep(2 * time.Second)
	// Overlay Text
	color, err = strconv.Atoi("4291297280")
	if err != nil {
		return errors.New("overlay text metadata config color not as expected: " + err.Error())
	}
	width, err = strconv.ParseFloat("1600", 64)
	if err != nil {
		return errors.New("overlay text metadata config width not as expected: " + err.Error())
	}
	height, err = strconv.ParseFloat("150", 64)
	if err != nil {
		return errors.New("overlay text metadata config height not as expected: " + err.Error())
	}
	posx, err = strconv.ParseFloat("0", 64)
	if err != nil {
		return errors.New("overlay text metadata config posx not as expected: " + err.Error())
	}
	posy, err = strconv.ParseFloat("0", 64)
	if err != nil {
		return errors.New("overlay text metadata config posy not as expected: " + err.Error())
	}
	overlay_text_settings := map[string]interface{}{
		"text":   "test",
		"color1": color,
		"color2": color,
	}
	if !overlay_text_exists {
		_, err = obs.CreateInput(SourceTextType, current_scene, obs.OverlayTextSourceName, true, overlay_text_settings)
		if err != nil {
			return errors.New(err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	_, err = obs.SetSceneItemTransform(obs.GetSceneItemId(current_scene, obs.OverlayTextSourceName), current_scene, posx, posy, width, height)
	if err != nil {
		return errors.New("Overlay Text: " + err.Error())
	}
	time.Sleep(2 * time.Second)
	err = obs.PressInputPropertiesButton(&inputs.PressInputPropertiesButtonParams{
		InputName:    obs.AvatarSourceName,
		PropertyName: "refreshnocache",
	})
	if err != nil {
		return errors.New(err.Error())
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
	color, err := ConvertColor(task.Color)
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
		background_color, err := ConvertColor(task.Background.Color)
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
	} else {
		if task_exists {
			_, err = obs.RemoveSceneItem(obs.GetSceneItemId(current_scene, obs.TaskSourceName), current_scene)
			if err != nil {
				return errors.New("Error Removing task scene item: " + err.Error())
			}
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

func (obs *OBS) PressInputPropertiesButton(params *inputs.PressInputPropertiesButtonParams) error {
	_, err := obs.Client.Inputs.PressInputPropertiesButton(params)
	return err
}

func (obs *OBS) SetSceneItemEnabled(item_id float64, scene_name string, enabled bool) error {
	_, err := obs.Client.SceneItems.SetSceneItemEnabled(&sceneitems.SetSceneItemEnabledParams{
		SceneItemEnabled: &enabled,
		SceneItemId:      item_id,
		SceneName:        scene_name,
	})
	return err
}
