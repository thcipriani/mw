/*
Copyright © 2020 Kosta Harlan <kosta@kostaharlan.net>
Copyright © 2020 Tyler Cipriani <tcipriani@wikimedia.org>
Copyright © 2020 Brennen Bearnes <bbearnes@wikimedia.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"github.com/kostajh/mw/setup"
)

func runStage(cmd *cobra.Command, args []string) {
}

var startCmd = &cobra.Command{
	Use:   "stage",
	Short: "Stage a gerrit patch",
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = "Starting the development environment "
		s.Start()
		setup.Linux()
		command := exec.Command("docker-compose", "up", "-d")
		if isLinuxHost() {
			command.Env = os.Environ()
			command.Env = append(
				command.Env,
				fmt.Sprintf("MW_DOCKER_UID=%s", string(os.Getuid())),
				fmt.Sprintf("MW_DOCKER_GID=%s", string(os.Getgid())))
		}
		stdoutStderr, _ := command.CombinedOutput()
		fmt.Print(string(stdoutStderr))
		s.Stop()
		handlePortError(stdoutStderr)

		if composerDependenciesNeedInstallation() {
			promptToInstallComposerDependencies()
		}

		printSuccess()
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if !isInCoreDirectory() {
			os.Exit(1)
		}
		if isLinuxHost() {
			setup.Linux()
		}
	},
}


