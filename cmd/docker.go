/*
Copyright © 2020 Kosta Harlan <kosta@kostaharlan.net>

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

func composeOverride() {
}

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Provides subcommands for interacting with MediaWiki's docker development environment",
	RunE:  nil,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = "Starting the development environment "
		s.Start()
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
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop development environment",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !isInCoreDirectory() {
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = "Stopping development environment "
		s.Start()
		command := exec.Command("docker-compose", "stop")
		stdoutStderr, err := command.CombinedOutput()
		s.Stop()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", stdoutStderr)
	},
}

func printSuccess() {
	portCommand := exec.Command("docker-compose", "port", "mediawiki", "8080")
	portCommandOutput, _ := portCommand.CombinedOutput()
	// Replace 0.0.0.0 with localhost
	fmt.Printf("Success! View MediaWiki-Docker at http://%s",
		strings.Replace(string(portCommandOutput), "0.0.0.0", "localhost", 1))
}

func handlePortError(stdoutStderr []byte) {
	portError := strings.Index(string(stdoutStderr), " failed: port is already allocated")
	if portError > 0 {
		// TODO: This breaks if someone set port 80 for example.
		fmt.Println(string(stdoutStderr))
		fmt.Printf("Port %s is already allocated! \n\nPlease override the port via a .env file, see https://www.mediawiki.org/wiki/MediaWiki-Docker for instructions\n",
			string(stdoutStderr[portError-4:])[0:4])
		os.Exit(1)
	}
}

func promptToInstallComposerDependencies() {
	fmt.Println("MediaWiki has some external dependencies that need to be installed")
	prompt := promptui.Prompt{
		IsConfirm: true,
		Label:     "Install dependencies now",
	}
	_, err := prompt.Run()
	if err == nil {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = "Installing Composer dependencies (this may take a few minutes) "
		s.Start()
		os.Mkdir("cache", 0700)
		depsCommand := exec.Command(
			"docker-compose",
			"exec",
			"-T",
			"mediawiki",
			"composer",
			"update",
		)
		out, err := depsCommand.CombinedOutput()
		if err != nil {
			fmt.Print(string(out))
			log.Fatal(err)
			os.Exit(1)
		}
		s.Stop()
	}
}

func composerDependenciesNeedInstallation() bool {
	// Detect if composer dependencies are not installed and prompt user to install
	dependenciesCheck := exec.Command(
		"docker-compose",
		"exec",
		"-T",
		"mediawiki",
		"php",
		"maintenance/install.php",
		"--help",
	)
	stdoutStderr, _ := dependenciesCheck.CombinedOutput()
	return strings.Index(string(stdoutStderr), " dependencies that need to be installed") > 0
}

func isInCoreDirectory() bool {
	if _, err := os.Stat("README.mediawiki"); err == nil {
		return true
	}
	fmt.Println("❌ Please run this command within the root of the MediaWiki core repository.")
	return false
}

func isLinuxHost() bool {
	unameCommand := exec.Command("uname")
	stdout, err := unameCommand.CombinedOutput()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return string(stdout) == "Linux\n"
}

func init() {
	rootCmd.AddCommand(dockerCmd)

	dockerCmd.AddCommand(startCmd)
	dockerCmd.AddCommand(stopCmd)
}
