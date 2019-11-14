package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var busRoute = flag.String("bus-route", "", "mention the route of bus")
var busStopName = flag.String("bus-stop-name", "", "Mention the name of the bus stop")
var busDirection = flag.String("bus-direction", "", "can take values like north, south, east, west")

func flagChecker() {
	flag.Parse()

	if len(*busRoute) < 1 {
		fmt.Println("Missing bus route which is mandtory")
		flag.Usage()
		os.Exit(1)
	}
}

type RouteDescription struct {
	Description string `json:"Description"`
	ProviderID  string `json:"ProviderID"`
	Route       string `json:"Route"`
}

type BusDirectionCheck struct {
	Text  string `json:"Text"`
	Value int    `json:"Value,string"`
}

type StopValuesCheck struct {
	Text  string `json:"Text"`
	Value string `json:"Value"`
}

type NextTrip struct {
	Actual           bool   `json:"Actual"`
	BlockNumber      int    `json:"BlockNumber"`
	DepartureText    string `json:"DepartureText"`
	DepartureTime    string `json:"DepartureTime"`
	Description      string `json:"Description"`
	Gate             string `json:"Gate"`
	Route            string `json:"Route"`
	RouteDirection   string `json:"RouteDirection"`
	Terminal         string `json:"Terminal"`
	VehicleHeading   int    `json:"VehicleHeading"`
	VehicleLatitude  int    `json:"VehicleLatitude"`
	VehicleLongitude int    `json:"VehicleLongitude"`
}

// finding theroute number from the busRoute
func routeNumber() string {
	url := "http://svc.metrotransit.org/NexTrip/Routes?format=json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("response code: %d", err)
	}
	body, err := ioutil.ReadAll(resp.Body) // using ioutil funtin to read the body of resp.
	if err != nil {
		panic(err.Error())
	}

	var routeDescription []RouteDescription

	// routeDescription := []RouteDescription{}
	jsonErr := json.Unmarshal(body, &routeDescription)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for i := range routeDescription {
		if *busRoute == routeDescription[i].Description {
			// i := RouteDescription{}
			routeNo := routeDescription[i].Route
			return routeNo
		}
	}
	// fmt.Println("Route not detected")
	return "1"
}

// This func returns the value of route depending on direction
func busDirectionRetriver() int {
	switch *busDirection {
	case "south":
		return 1
	case "east":
		return 2
	case "west":
		return 3
	case "north":
		return 4
	default:
		panic("unrecognized bus direction, Available options south,east,north,west")
	}
}

// this function verifies the direction of a specific route
func busDirectionChecker() int {
	url := "http://svc.metrotransit.org/NexTrip/Directions/"
	urlContextPath := url + routeNumber() + "?format=json"
	// fmt.Println(urlContextPath)
	resp, err := http.Get(urlContextPath)
	if err != nil {
		fmt.Printf("response code: %d", err)
	}
	body, err := ioutil.ReadAll(resp.Body) // using ioutil funtin to read the body of resp.
	if err != nil {
		panic(err.Error())
	}

	var busDirectionCheck []BusDirectionCheck

	// routeDescription := []RouteDescription{}
	jsonErr := json.Unmarshal(body, &busDirectionCheck)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	for i := range busDirectionCheck {
		if busDirectionRetriver() == busDirectionCheck[i].Value {
			return busDirectionCheck[i].Value
		}
	}
	panic("Direction doesnt match the route")
}

// returns value of the stop
func stopValues() string {
	url := "http://svc.metrotransit.org/NexTrip/Stops/"
	busDirectionToString := strconv.Itoa(busDirectionRetriver())
	contextPath := url + routeNumber() + "/" + busDirectionToString + "?format=json"
	// fmt.Println(contextPath)
	resp, err := http.Get(contextPath)
	if err != nil {
		fmt.Printf("response code: %d", err)
	}
	body, err := ioutil.ReadAll(resp.Body) // using ioutil funtin to read the body of resp.
	if err != nil {
		panic(err.Error())
	}

	var stopValuesCheck []StopValuesCheck

	// routeDescription := []RouteDescription{}
	jsonErr := json.Unmarshal(body, &stopValuesCheck)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	for i := range stopValuesCheck {
		if *busStopName == stopValuesCheck[i].Text {
			return stopValuesCheck[i].Value
		}
	}
	panic("Stop doesnt belong to the route")
}

// returns the bus time
func nextBusTrip() string {
	url := "http://svc.metrotransit.org/NexTrip/"
	busDirectionToString := strconv.Itoa(busDirectionRetriver())
	contextUrl := url + routeNumber() + "/" + busDirectionToString + "/" + stopValues() + "?format=json"
	// fmt.Println(contextUrl)

	resp, err := http.Get(contextUrl)
	if err != nil {
		fmt.Printf("response code: %d", err)
	}
	body, err := ioutil.ReadAll(resp.Body) // using ioutil funtin to read the body of resp.
	if err != nil {
		panic(err.Error())
	}

	var nextTrip []NextTrip

	jsonErr := json.Unmarshal(body, &nextTrip)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return nextTrip[0].DepartureText

}

// http://svc.metrotransit.org/NexTrip/Providers
// List all provides

// http://svc.metrotransit.org/NexTrip/Routes
// List all routes and route numbers

// http://svc.metrotransit.org/NexTrip/Directions/901?format=json
// http://svc.metrotransit.org/NexTrip/Directions/{ROUTE}?format=json
// gives northbound or eastbound

// http://svc.metrotransit.org/NexTrip/Stops/901/4?format=json
// http://svc.metrotransit.org/NexTrip/Stops/{ROUTE}/{DIRECTION}
// gives stops and stop values depending on the direction like south or north
// need to use stop value from here to get time of departure

// http://svc.metrotransit.org/NexTrip/56330?format=json
// http://svc.metrotransit.org/NexTrip/{StopID}
// gives all routes of a stop

// http://svc.metrotransit.org/NexTrip/901/4/TF21?format=json
// http://svc.metrotransit.org/NexTrip/{ROUTE}/{DIRECTION}/{STOP}

func main() {

	flagChecker()
	routeNo := routeNumber()
	if routeNo == "1" {
		fmt.Println("Route name not detected")
		fmt.Println("Please use a valid name")
	}

	busDirectionRetriver()
	busDirectionChecker()
	stopValues()

	k := nextBusTrip()
	fmt.Println(k)
}
