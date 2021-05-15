# chia-bouncer

![really](https://media.giphy.com/media/5fBH6zf7l8bxukYh74Q/giphy.gif)

Tiny CLI tool to remove unwanted connections from your Chia Node based on the Geo IP Location (Country). 
The Tool is written in golang and acts more like a command line wrapper around `chia show...`
and the cli tool `geoiplookup`

## Requirements

- `chia` is installed
- `geoiplookup` is installed (on Ubuntu 18 and Ubuntu 20 the package is called `geoip-bin`)
- update the `geoiplookup` IPv4 database (more details see below)

### Install `geoiplookup` and update the database 

The tool `geoiplookup` is required to perform the Geo IP Location lookup.
Ubuntu 18 and Ubuntu 20 users can simply install it via:

```shell
sudo apt-get install geoip-bin
```

The installed package contains a database which is pretty old. 
To update it, download the newest from [here](https://dl.miyuru.lk/geoip/dbip/country/dbip4.dat.gz)
and unpack to `/usr/share/GeoIP`

```shell
wget https://dl.miyuru.lk/geoip/dbip/country/dbip4.dat.gz
gunzip dbip4.dat.gz
sudo mv dbip4.dat /usr/share/GeoIP/GeoIP.dat
```

You can check the age of the Geo IP Location database by running
```shell
geoiplookup localhost -v
   GeoIP Country Edition: DBIPLite-Country-CSV_20210401 converted to legacy DB with sherpya/geolite2legacy by miyuru.lk
   ...
```

## Usage

If your chia executable is located in `$HOME/chia-blockchain/venv/bin/chia`, you can simply run:
```bash
# assumes chia executable is located in $HOME/chia-blockchain/venv/bin/chia
# removes all connections from mars
chia-bouncer mars
```
To specify a custom path to your chia executable use `--chiaexec` or `-e`
```bash
# custom defined chia executable
# removes all connections from "elon on mars"
chia-bouncer -e /home/steffen/chia-blockchain/venv/bin/chia elon on mars
```
Call with `--help` or `-h` to see the help page 
```bash
# custom defined chia executable
chia-bouncer -c
 NAME:
    chia-bouncer - remove nodes by given geo ip location from your connections

 USAGE:
    chia-bouncer -ce /home/steffen/chia-blockchain/venv/bin/chia mars

 VERSION:
    unknown

 DESCRIPTION:
    Tool will lookup connections via 'chia show -c', get ip locations via geoiplookup and remove nodes from specified location via 'chia show -r'

 AUTHOR:
    st3ffn <funk.up.up@gmail.com>

 COMMANDS:
    help, h  Shows a list of commands or help for one command

 GLOBAL OPTIONS:
    --chiaexec CHIA-EXECUTABLE, -e CHIA-EXECUTABLE  CHIA-EXECUTABLE. normally located inside the bin folder of your venv directory (default: $HOME/chia-blockchain/venv/bin/chia)
    --version, -v                                   print the version (default: false)
    --help, -h                                      show help (default: false)

 COPYRIGHT:
    GNU GPLv3
```

## Building from source

To build from the sources go 1.16 is required.
```shell
go build
```
If you want to build it for your target linux target system (amd64) you can use the `build-linux.sh` shell script

## Integration

The script can easily be used with cron. Simply open the users crontab via `crontab -e` and add the following line.

```shell
# run chia-bouncer every 2 hours and remove all connections from mars
0 */2 * * * /PATH/TO/chia-bouncer mars

```