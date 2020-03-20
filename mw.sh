#!/bin/bash

set -euo pipefail

haz() {
    command -v "$1" &> /dev/null
}

SRC='./src'

[ ! -d "$SRC" ] || exit 1
haz 'docker' || exit 1
haz 'docker-compose' || exit 1

MW_DOCKER_UID=$(id -u)
MW_DOCKER_GID=$(id -g)

mkdir "$SRC"
pushd "$SRC" || exit
git init
git fetch --depth 1 "https://gerrit.wikimedia.org/r/mediawiki/core" refs/changes/04/577004/11 && git cherry-pick FETCH_HEAD

cat << END > docker-compose.override.yml
version: '3.7'
services:
  mediawiki:
    # On Linux, these lines ensure file ownership is set to your host user/group
    user: "${MW_DOCKER_UID}:${MW_DOCKER_GID}"
END

docker-compose up -d
docker-compose exec mediawiki composer update
docker-compose exec mediawiki bash -c 'php maintenance/install.php \
--server $MW_SERVER \
--scriptpath=$MW_SCRIPTPATH \
--dbtype $MW_DBTYPE \
--dbpath $MW_DBPATH \
--lang $MW_LANG \
--pass $MW_PASS \
$MW_SITENAME $MW_USER'

git clone --depth=1 https://github.com/wikimedia/mediawiki-skins-vector "$SRC/skins/Vector"
echo "wfLoadSkin( 'Vector' );" >> "$SRC/LocalSettings.php"
