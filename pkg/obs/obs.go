package obs

import (
	"fmt"
	"strconv"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/typedefs"
)

type OBS struct {
	Client *goobs.Client
}

type Task struct {
	Text       string
	PosX       float64
	PosY       float64
	Width      float64
	Height     float64
	Color      Color
	Background *Background
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type Background struct {
	Color Color
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
	input_name := "test"
	background_name := "background"
	_, err := obs.RemoveSceneItem(obs.GetSceneItemId("Main", input_name), "Main")
	if err != nil {
		fmt.Println("Error Removing scene item: " + err.Error())
	}
	_, err = obs.RemoveSceneItem(obs.GetSceneItemId("Main", background_name), "Main")
	if err != nil {
		fmt.Println("Error Removing scene item: " + err.Error())
	}
	color, err := convertColor(task.Color)
	if err != nil {
		fmt.Println("Error converting color: " + err.Error())
	}
	background_color, err := convertColor(task.Background.Color)
	if err != nil {
		fmt.Println("Error converting background color: " + err.Error())
	}
	_, err = obs.CreateInput("color_source_v3", "Main", "background", true, map[string]interface{}{ // color_source_v3
		"text":   task.Text,
		"color":  background_color,
		"width":  task.Width + 4,
		"height": task.Height + 4,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = obs.CreateInput("text_ft2_source_v2", "Main", "test", true, map[string]interface{}{ // color_source_v3
		"text":   task.Text,
		"color1": color,
		"color2": color,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	trans2, _ := obs.GetSceneItemTransform(obs.GetSceneItemId("Main", background_name), "Main")
	if err != nil {
		fmt.Println("Error getting scene item transform: " + err.Error())
	}
	_, err = obs.SetSceneItemTransform(obs.GetSceneItemId("Main", background_name), "Main", task.PosX-2, task.PosY-2, task.Width+4, task.Height+4, trans2.SceneItemTransform.SourceWidth, trans2.SceneItemTransform.SourceHeight)
	if err != nil {
		fmt.Println(err.Error())
	}

	trans1, _ := obs.GetSceneItemTransform(obs.GetSceneItemId("Main", input_name), "Main")
	if err != nil {
		fmt.Println("Error getting scene item transform: " + err.Error())
	}
	_, err = obs.SetSceneItemTransform(obs.GetSceneItemId("Main", input_name), "Main", task.PosX, task.PosY, task.Width, task.Height, trans1.SceneItemTransform.SourceWidth, trans1.SceneItemTransform.SourceHeight)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (obs *OBS) GetInput(name string) (*inputs.GetInputSettingsResponse, error) {
	return obs.Client.Inputs.GetInputSettings(&inputs.GetInputSettingsParams{
		InputName: name,
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
	return obs.Client.Inputs.CreateInput(&inputs.CreateInputParams{
		InputKind:        kind,
		SceneName:        scene,
		InputName:        name,
		SceneItemEnabled: &enabled,
		InputSettings:    settings,
	})
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
