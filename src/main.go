package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

type config struct {
	API struct {
		Key string `yaml:"key"`
	} `yaml:"api"`
	Location struct {
		Latitude  string `yaml:"latitude"`
		Longitude string `yaml:"longitude"`
	} `yaml:"location"`
	Preferences struct {
		Unit string `yaml:"unit"`
	} `yaml:"preferences"`
}

type weather struct {
	Lat      float32
	Lon      float32
	Timezone string
	Current  struct {
		Temp       float32
		Feels_like float32
		Humidity   float32
		Weather    []struct {
			Description string
		}
	}
}

var usrCFG string
var cfg config

// check if config file is present
func configCheck() bool {
	usrCFG, _ = os.UserConfigDir()
	usrCFG = usrCFG + "/wthr/config.yml"

	if _, err := os.Stat(usrCFG); err == nil {
		return true
	}

	return false
}

// parses config.yml to a config struct (cfg)
func configFetch() {
	cfgFile, err := os.Open(usrCFG)
	if err != nil {
		panic(err)
	}
	defer cfgFile.Close()

	decoder := yaml.NewDecoder(cfgFile)

	err = decoder.Decode(&cfg)

	// check if config has empty options
	if cfg.API.Key == "" || cfg.Location.Latitude == "" || cfg.Location.Longitude == "" || cfg.Preferences.Unit == "" {
		fmt.Println(aurora.Red("config file options empty?"))
		os.Exit(1)
	}

	if err != nil {
		if err.Error() == "EOF" {
			fmt.Println(aurora.Red("malformed config"))
			os.Exit(1)
		}
	}

}

func getRes(url string) *http.Response {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	return res
}

func jsonRes(res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	res.Body.Close()

	var x json.RawMessage

	json.Unmarshal(body, &x)

	return x
}

func getWeather(url string) weather {
	var j weather
	x := getRes(url)
	y := jsonRes(x)
	json.Unmarshal(y, &j)

	return j

}

// only god knows
func wtf(x []struct{ Description string }) string {
	str := fmt.Sprintf("%v", x)
	str = str[2 : len(str)-2]
	return str
}

func main() {
	if !configCheck() {
		fmt.Println(aurora.Red("Couldn't locate config file"))
		os.Exit(1)
	}

	configFetch()

	url := "https://api.openweathermap.org/data/2.5/onecall?lat=" + cfg.Location.Latitude + "&lon=" + cfg.Location.Longitude + "&exclude=minutely,hourly,daily" + "&units=" + cfg.Preferences.Unit + "&appid=" + cfg.API.Key

	j := getWeather(url)

	fmt.Println(j.Lat, "/", j.Lon)
	fmt.Println(j.Timezone)
	fmt.Println("it feels like", j.Current.Feels_like)
	fmt.Println(wtf(j.Current.Weather))

}
