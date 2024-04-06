package main

import (
	"encoding/json"
	"fmt"
	"github.com/choria-io/fisk"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/used255/go-aht20"
	"os"
)

var (
	debug        bool
	bus          int
	jsonFormat   bool
	choriaFormat bool
	labels       map[string]string
)

func main() {
	labels = make(map[string]string)

	app := fisk.New("aht20", "Client for the AHT20")
	app.Flag("bus", "i2c bus number to use").Default("1").IntVar(&bus)
	app.Flag("debug", "Enable debug logging").UnNegatableBoolVar(&debug)

	read := app.Command("read", "Reads from the sensor").Action(read)
	read.Flag("json", "Enables JSON output").UnNegatableBoolVar(&jsonFormat)
	read.Flag("choria", "Enables Choria Metric output").UnNegatableBoolVar(&choriaFormat)
	read.Flag("label", "Labels to apply to Choria Metric output").StringMapVar(&labels)

	app.Command("reset", "Resets the sensor").Action(reset)

	app.MustParseWithUsage(os.Args[1:])
}

func connect() (*i2c.I2C, *aht20.Device, error) {
	lvl := logger.InfoLevel
	if debug {
		lvl = logger.DebugLevel
	}
	logger.ChangePackageLogLevel("i2c", lvl)
	logger.ChangePackageLogLevel("aht20", lvl)

	bus, err := i2c.NewI2C(0x38, bus)
	if err != nil {
		return nil, nil, err
	}

	s := aht20.New(bus)
	s.Configure()

	return bus, &s, nil
}

func reset(_ *fisk.ParseContext) error {
	bus, s, err := connect()
	if err != nil {
		return err
	}
	defer bus.Close()

	s.Reset()

	return nil
}

func read(_ *fisk.ParseContext) error {
	bus, s, err := connect()
	if err != nil {
		return err
	}
	defer bus.Close()

	err = s.Read()
	if err != nil {
		return err
	}

	switch {
	case jsonFormat:
		data := map[string]any{
			"temperature":     s.Celsius(),
			"humidity":        s.RelHumidity(),
			"raw_temperature": s.RawTemp(),
			"raw_humidity":    s.RawHumidity(),
		}
		j, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(j))
	case choriaFormat:
		data := map[string]any{
			"labels": labels,
			"metrics": map[string]any{
				"temperature":     s.Celsius(),
				"humidity":        s.RelHumidity(),
				"temperature_raw": s.RawTemp(),
				"humidity_raw":    s.RawHumidity(),
			}}
		j, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(j))
	default:
		fmt.Printf("Temperature: %.2fC\n", s.Celsius())
		fmt.Printf("Relative Humidity: %.2f%%\n", s.RelHumidity())
	}

	return nil
}
