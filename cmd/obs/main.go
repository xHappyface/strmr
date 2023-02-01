package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/jnrprgmr/dog/pkg/obs"
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
	cli.obs.SetTask(cli.task)
}

func main() {
	client, err := goobs.New("localhost:4455", goobs.WithPassword("test123"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	obs := obs.New(client)
	obsCli := New(obs)
	obsCli.Execute()
}
