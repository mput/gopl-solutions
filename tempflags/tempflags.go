package tempflags

import (
	"flag"
	"fmt"

	"mput.me/gopl/tempconv"
)

type celsiusFlag struct{ tempconv.Celsius }

func (c *celsiusFlag) Set(s string) error {
	var val float64
	var unit string
	fmt.Sscanf(s, "%f%s", &val, &unit)
	switch unit {
	case "C", "°C":
		c.Celsius = tempconv.Celsius(val)
	case "K", "°K":
		c.Celsius = tempconv.KelvinToCelsius(val)
	case "F", "°F":
		c.Celsius = tempconv.FahrenheitToCelsius(val)
	default:
		return fmt.Errorf("invalid format of temperature %s", s)
	}
	return nil
}

func CelsiusFlag(name string, value tempconv.Celsius, usage string) *tempconv.Celsius {
	v := celsiusFlag{value}
	flag.Var(&v, name, usage)
	return &v.Celsius
}
