package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"math"
	"net/http"
	"os"
	"testtask/cmd/structure"
)

const (
	zoneCenterLatitude  = 57.988962
	zoneCenterLongitude = 56.204668
	zoneRadius          = float64(15)
)

func main() {
	app := cli.NewApp()
	app.Name = "Parent notificator simulator"
	app.Usage = ""

	app.Commands = []cli.Command{
		notify,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("exit")
		os.Exit(1)
	}
}

var flagApi = cli.StringFlag{
	Name:  "target",
	Usage: "full http url that will accept get request on /notify/parent",
	Value: "http://127.0.0.1:4080/notify/parent",
}

var flagAddr = cli.StringFlag{
	Name:  "addr",
	Usage: "",
	Value: "127.0.0.1:4079",
}

var notify = cli.Command{
	Name:      "start",
	Action:    notifyAction,
	Usage:     "Determine when the child has already left the zone and notify the parent",
	UsageText: ``,
	Flags: []cli.Flag{
		flagApi,
		flagAddr,
	},
}

func notifyAction(c *cli.Context) error {
	addr := c.String(flagAddr.Name)
	handler := http.NewServeMux()
	server := http.Server{Addr: addr, Handler: handler}

	handler.HandleFunc("/set-coord", func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var childCoords structure.Coordinate
		err := decoder.Decode(&childCoords)
		if err != nil {
			panic(err)
		}

		if !isChildInTheZone(&childCoords) {
			target := c.String(flagApi.Name)
			resp, err := http.Get(target)
			if err != nil {
				logrus.WithError(err).Fatal("something wrong")
			}
			if resp.StatusCode != http.StatusOK {
				logrus.Fatalf("bad status code, expected %d got %d", http.StatusOK, resp.StatusCode)
			}
			logrus.Info("successfully...")
			logrus.Infof("%+v\n", childCoords)
		}
	})

	logrus.Infof("listening: %s/set-coord", addr)
	return server.ListenAndServe()
}

/*
 	Определяю расстояние между двумя координатами по формуле гаверсинуса.
	Если расстояние между позицией ребенка и центром зоны превышает
	радиус зоны + погрешность, то делаю вывод, что ребенок покинул зону
*/
func isChildInTheZone(coordinate *structure.Coordinate) bool {
	// Надеюсь, я правильно понял)
	if coordinate.Reason != structure.ReasonStatusOk && coordinate.Source != structure.SourceFused {
		return true
	}

	const EarthRadius = 6378.137
	dLat := zoneCenterLatitude*math.Pi/180 - coordinate.Latitude*math.Pi/180
	dLon := zoneCenterLongitude*math.Pi/180 - coordinate.Longitude*math.Pi/180
	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(coordinate.Latitude*math.Pi/180)*math.Cos(zoneCenterLatitude*math.Pi/180)*
			math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := EarthRadius * c

	distanceBetweenPointsInMeters := d * 1000

	return zoneRadius+coordinate.Accuracy >= distanceBetweenPointsInMeters
}
