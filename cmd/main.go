package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MovieStoreGuy/signalview/pkg/client"
	"github.com/MovieStoreGuy/signalview/pkg/config"
	"github.com/MovieStoreGuy/signalview/pkg/endpoint/detector"
	"github.com/MovieStoreGuy/signalview/pkg/endpoint/teams"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	v = viper.New()
)

type DerpStruct struct {
	Name      string
	Detectors []detector.DetailResultsPayload
}

func main() {
	conf, err := config.NewDefault(v)
	switch err {
	case pflag.ErrHelp:
		return
	}
	c := client.NewConfiguredClient(
		conf.ConfigureClient,
	)
	log := logrus.New()
	log.SetOutput(os.Stdout)
	detectors, err := detector.GetMatching(context.Background(), log, c, conf, nil)
	if err != nil {
		log.WithError(err).Fatal("unable to process request")
		return
	}
	team, err := teams.GetMatching(context.Background(), log, c, conf, nil)
	if err != nil {
		log.WithError(err).Fatal("unable to process request")
		return
	}
	filter := map[string]*DerpStruct{}
	for _, t := range team.Results {
		filter[t.ID] = &DerpStruct{
			Name: t.Name,
		}
	}
	for _, det := range detectors.Results {
		for _, t := range det.Teams {
			if info, exist := filter[t]; exist {
				info.Detectors = append(info.Detectors, det)
			}
		}
	}
	for _, t := range filter {
		fmt.Println("Team, id, detector name, created, modified, by")
		for _, det := range t.Detectors {
			fmt.Println(t.Name, ", ", det.ID, ", ", det.Name, ", ", det.Created, ", ", det.LastUpdated, ",", det.Creator)
		}
	}
}
