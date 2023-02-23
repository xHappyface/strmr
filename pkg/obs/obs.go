package obs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/typedefs"
)

type OBS struct {
	Client *goobs.Client
}

type Task struct {
	Text       string      `json:"text"`
	PosX       float64     `json:"pos_x"`
	PosY       float64     `json:"pos_y"`
	Width      float64     `json"width"`
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

func New(client *goobs.Client) *OBS {
	return &OBS{
		Client: client,
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

func (obs *OBS) SetTask(task Task) error {
	background_name := "background"
	input_name := "test"
	background_exists := true
	task_exists := true
	_, err := obs.GetInput(background_name)
	if err != nil {
		background_exists = false
	}
	_, err = obs.GetInput(input_name)
	if err != nil {
		task_exists = false
	}
	background_color, err := convertColor(task.Background.Color)
	if err != nil {
		return errors.New("Error converting background color: " + err.Error())
	}
	background_settings := map[string]interface{}{
		"color":  background_color,
		"width":  task.Width + 4,
		"height": task.Height + 4,
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
	if background_exists {
		if task.Background != nil {
			_, err := obs.SetInputSettings(background_name, background_settings)
			if err != nil {
				return errors.New(err.Error())
			}
			trans2, err := obs.GetSceneItemTransform(obs.GetSceneItemId("Main", background_name), "Main")
			if err != nil {
				return errors.New("Error getting scene item transform: " + err.Error())
			}
			_, err = obs.SetSceneItemTransform(obs.GetSceneItemId("Main", background_name), "Main", task.PosX-2, task.PosY-2, task.Width+4, task.Height+4, trans2.SceneItemTransform.SourceWidth, trans2.SceneItemTransform.SourceHeight)
			if err != nil {
				return errors.New(err.Error())
			}
		} else {
			_, err = obs.RemoveSceneItem(obs.GetSceneItemId("Main", background_name), "Main")
			if err != nil {
				return errors.New("Error Removing scene item: " + err.Error())
			}
		}
	} else {
		if task.Background != nil {
			_, err = obs.CreateInput("color_source_v3", "Main", background_name, true, background_settings)
			if err != nil {
				return errors.New(err.Error())
			}
			trans2, err := obs.GetSceneItemTransform(obs.GetSceneItemId("Main", background_name), "Main")
			if err != nil {
				return errors.New("Error getting scene item transform: " + err.Error())
			}
			_, err = obs.SetSceneItemTransform(obs.GetSceneItemId("Main", background_name), "Main", task.PosX-2, task.PosY-2, task.Width+4, task.Height+4, trans2.SceneItemTransform.SourceWidth, trans2.SceneItemTransform.SourceHeight)
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	if len(task.Text) > 0 {
		if task_exists {
			_, err := obs.SetInputSettings(input_name, task_settings)
			if err != nil {
				return errors.New(err.Error())
			}

		} else {
			_, err = obs.CreateInput("text_ft2_source_v2", "Main", input_name, true, task_settings)
			if err != nil {
				return errors.New(err.Error())
			}
		}
		trans, err := obs.GetSceneItemTransform(obs.GetSceneItemId("Main", input_name), "Main")
		if err != nil {
			return errors.New("Error getting scene item transform: " + err.Error())
		}
		_, err = obs.SetSceneItemTransform(obs.GetSceneItemId("Main", input_name), "Main", task.PosX, task.PosY, task.Width, task.Height, trans.SceneItemTransform.SourceWidth, trans.SceneItemTransform.SourceHeight)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
	// _, err := obs.RemoveSceneItem(obs.GetSceneItemId("Main", input_name), "Main")
	// if err != nil {
	// 	fmt.Println("Error Removing scene item: " + err.Error())
	// }
	// if task.Background != nil {
	// 	_, err = obs.RemoveSceneItem(obs.GetSceneItemId("Main", background_name), "Main")
	// 	if err != nil {
	// 		fmt.Println("Error Removing scene item: " + err.Error())
	// 	}
	// }
	// if len(task.Text) > 0 {
	// 	color, err := convertColor(task.Color)
	// 	if err != nil {
	// 		fmt.Println("Error converting color: " + err.Error())
	// 	}
	// 	if task.Background != nil {
	// 		background_color, err := convertColor(task.Background.Color)
	// 		if err != nil {
	// 			fmt.Println("Error converting background color: " + err.Error())
	// 		}
	// 		_, err = obs.CreateInput("color_source_v3", "Main", "background", true, map[string]interface{}{ // color_source_v3
	// 			"text":   task.Text,
	// 			"color":  background_color,
	// 			"width":  task.Width + 4,
	// 			"height": task.Height + 4,
	// 		})
	// 		if err != nil {
	// 			fmt.Println(err.Error())
	// 		}
	// 	}
	// 	_, err = obs.CreateInput("text_ft2_source_v2", "Main", "test", true, map[string]interface{}{ // color_source_v3
	// 		"text":   task.Text,
	// 		"color1": color,
	// 		"color2": color,
	// 	})
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	if task.Background != nil {
	// 		trans2, _ := obs.GetSceneItemTransform(obs.GetSceneItemId("Main", background_name), "Main")
	// 		if err != nil {
	// 			fmt.Println("Error getting scene item transform: " + err.Error())
	// 		}
	// 		_, err = obs.SetSceneItemTransform(obs.GetSceneItemId("Main", background_name), "Main", task.PosX-2, task.PosY-2, task.Width+4, task.Height+4, trans2.SceneItemTransform.SourceWidth, trans2.SceneItemTransform.SourceHeight)
	// 		if err != nil {
	// 			fmt.Println(err.Error())
	// 		}
	// 	}
	// 	trans1, _ := obs.GetSceneItemTransform(obs.GetSceneItemId("Main", input_name), "Main")
	// 	if err != nil {
	// 		fmt.Println("Error getting scene item transform: " + err.Error())
	// 	}
	// 	_, err = obs.SetSceneItemTransform(obs.GetSceneItemId("Main", input_name), "Main", task.PosX, task.PosY, task.Width, task.Height, trans1.SceneItemTransform.SourceWidth, trans1.SceneItemTransform.SourceHeight)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	return err
	// }
	// return nil
}

func (obs *OBS) GetInput(name string) (*inputs.GetInputSettingsResponse, error) {
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

func (obs *OBS) SetSceneItemTransform(item_id float64, name string, posX, posY, width, height, sourceWidth, sourceHeight float64) (*sceneitems.SetSceneItemTransformResponse, error) {
	return obs.Client.SceneItems.SetSceneItemTransform(&sceneitems.SetSceneItemTransformParams{
		SceneItemId: item_id,
		SceneName:   name,
		SceneItemTransform: &typedefs.SceneItemTransform{
			PositionX:    posX,
			PositionY:    posY,
			Alignment:    5,
			BoundsWidth:  1,
			BoundsHeight: 1,
			ScaleX:       width / sourceWidth,
			ScaleY:       height / sourceHeight,
			BoundsType:   "OBS_BOUNDS_NONE",
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
