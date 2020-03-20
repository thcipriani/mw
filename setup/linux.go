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
package setup

func dockerOverride() {
	_, err := os.Stat("docker-compose.override.yml")
	if err != nil {
		fmt.Println("Creating docker-compose.override.yml for correct user ID and group ID mapping from host to container")
		var data = `
version: '3.7'
services:
mediawiki:
user: "${MW_DOCKER_UID}:${MW_DOCKER_GID}"
`
		file, err := os.Create("docker-compose.override.yml")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.WriteString(data)
		if err != nil {
			log.Fatal(err)
		}
		file.Sync()
	}
}

func Linux() {
	// TODO: We should also check the contents for correctness, maybe
	// using docker-compose config and asserting that UID/GID mapping is present
	// and with correct values.
	dockerOverride()
}
