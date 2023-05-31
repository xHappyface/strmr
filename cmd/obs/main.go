package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/jnrprgmr/strmr/pkg/obs"
	"github.com/spf13/cobra"
)

type OBSCli struct {
	task    string
	rootCmd *cobra.Command
	obs     *obs.OBS
}

func New(obs *obs.OBS) *OBSCli {
	return &OBSCli{
		task: "",
		obs:  obs,
		rootCmd: &cobra.Command{
			Use:   "obs",
			Short: "OBStudios CLI",
			Long:  `Gives easy access to update obs settings by interacting with the webserver`,
		},
	}
}

func (cli *OBSCli) Execute() {
	cli.rootCmd.PersistentFlags().StringVar(&cli.task, "task", "", "set the on-screen task")
	cli.rootCmd.Run = cli.Run
	if err := cli.rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (cli *OBSCli) Run(cmd *cobra.Command, args []string) {
	t := obs.Task{
		Text:   cli.task,
		PosX:   33,
		PosY:   811,
		Width:  300,
		Height: 90,
		Color: obs.Color{
			R: 255,
			G: 0,
			B: 20,
			A: 255,
		},
		Background: &obs.Background{
			Color: obs.Color{
				R: 200,
				G: 200,
				B: 200,
				A: 255,
			},
		},
	}
	cli.obs.SetTask(t)
}

func main() {
	client, err := goobs.New("localhost:4455", goobs.WithPassword("test123"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	obs := obs.New(client, "strmr-screen", "strmr-task-text", "strmr-task-background", "strmr-avatar", "strmr-overlay-text", "strmr-overlay-background")
	obsCli := New(obs)
	obsCli.Execute()
}
