package commands

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/config"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldocker"
	"github.com/maliceio/malice/malice/persist"
	"github.com/maliceio/malice/plugins"
	"github.com/maliceio/malice/utils"
)

func cmdScan(path string, logs bool) {
	if len(path) > 0 {
		// Check that file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Fatal(path + ": no such file or directory")
		}

		docker := maldocker.NewDockerClient()

		// Check for existance of malice volume
		if _, exists, _ := docker.VolumeExists("malice"); !exists {
			log.Debug("Volume malice not found.")
			_, err := docker.CreateVolume("malice")
			er.CheckError(err)
		}
		log.Debug("Volume malice found.")

		if plugins.InstalledPluginsCheck(docker) {
			log.Debug("All enabled plugins are installed.")
		} else {
			// Prompt user to install all plugins?
			fmt.Println("All enabled plugins not installed would you like to install them now? (yes/no)\n[Warning] This can take a while if it is the first time you have ran Malice.")
			if util.AskForConfirmation() {
				plugins.UpdateAllPlugins(docker)
			}
		}

		file := persist.File{
			Path: path,
		}

		file.Init()
		// Output File Hashes
		file.ToMarkdownTable()
		// fmt.Println(string(file.ToJSON()))

		plugins.RunIntelPlugins(docker, file.MD5, true)

		log.Debug("Looking for plugins that will run on: ", file.Mime)
		// Iterate over all applicable installed plugins
		plugins := plugins.GetPluginsForMime(file.Mime, true)
		log.Debug("Found these plugins: ")
		for _, plugin := range plugins {
			log.Debugf(" - %v", plugin.Name)
		}

		for _, plugin := range plugins {
			log.Debugf(">>>>> RUNNING Plugin: %s >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", plugin.Name)
			// go func() {
			// Start Plugin Container
			// TODO: don't use the default of true for --logs
			cont, err := plugin.StartPlugin(docker, file.SHA256, true)
			er.CheckError(err)

			log.WithFields(log.Fields{
				"id": cont.ID,
				"ip": docker.GetIP(),
				// "url":      "http://" + maldocker.GetIP(),
				"name": cont.Name,
				"env":  config.Conf.Environment.Run,
			}).Debug("Plugin Container Started")

			err = docker.RemoveContainer(cont, false, false, false)
			er.CheckError(err)
			// }()
			// Clean up the Plugin Container
			// TODO: I want to reuse these containers for speed eventually.

			// time.Sleep(10 * time.Millisecond)
		}
		// time.Sleep(60 * time.Second)
		log.Debug("Done with plugins.")
	} else {
		log.Error("Please supply a valid file to scan.")
	}

}
