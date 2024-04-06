## AHT20 Client

A small client utility for reading and resetting the AHT20 or AM2301B sensor.

Tested on a Raspberry PI with Go 1.22

## Usage?

Basic usage information:

```nohighlight
$ aht --help
usage: aht20 [<flags>] <command> [<args> ...]

Client for the AHT20

Commands:
  read   Reads from the sensor
  reset  Resets the sensor

Global Flags:
  --help   Show context-sensitive help
  --bus=1  i2c bus number to use
  --debug  Enable debug logging
```

Reading the device:

```
$ aht read
Temperature: 24.75C
Relative Humidity: 49.05%
```

It supports JSON output:

```
$ aht read --json
{
  "humidity": 49.737072,
  "raw_humidity": 521531,
  "raw_temperature": 391857,
  "temperature": 24.740791
}
```

And also the format required by Choria Metric watchers:

```
$ aht20  read --choria --label location:home
{
  "labels": {
    "location": "home"
  },
  "metrics": {
    "humidity": 49.845886,
    "humidity_raw": 522672,
    "temperature": 24.757957,
    "temperature_raw": 391947
  }
}
```

## Contact?

R.I. Pienaar / rip@devco.net / [devco.net](https://www.devco.net/)