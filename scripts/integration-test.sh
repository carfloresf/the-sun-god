#!/usr/bin/env bash
export GIN_MODE=release

# verify if venom is already installed
if ! command -v venom &> /dev/null
then
    # download venom if not installed (only for MacOS ARM64)
    if [[ ! -f ./venom ]]
    then
        wget https://github.com/ovh/venom/releases/download/v1.1.0-beta.4/venom.darwin-arm64
        mv venom.darwin-arm64 venom
        chmod +x venom
    else
        echo "venom already downloaded."
    fi
else
  echo "venom already installed."
fi

# run binary
echo "--------------------- starting application ---------------------------"
$1 -configFile tests/config.yml &
sleep 5

echo "--------------------- integration tests starting ---------------------"
# run integration venom tests
./venom run tests
echo "--------------------- integration tests finished ---------------------"

# kill server
echo "killing $2"
while pkill $2; do
    sleep 1
done

# delete database file
echo "Deleting files."
rm -rf tests/feed-test.db
rm -rf *.log



