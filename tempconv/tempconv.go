package tempconv

import (
	"fmt"
)

type Celsius float64


func (c Celsius) String() string {
	return fmt.Sprintf("%g°C", c)
}

func FahrenheitToCelsius(f float64) Celsius {
	return Celsius((f - 32) * 5.0 / 9.0)
}

func KelvinToCelsius(f float64) Celsius {
	return Celsius(f - 273.15)
}
