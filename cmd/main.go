package main

import (
	"context"

	"github.com/MovieStoreGuy/signalview/pkg/client"
	"github.com/MovieStoreGuy/signalview/pkg/config"
	"github.com/MovieStoreGuy/signalview/pkg/endpoint/detector"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	v = viper.New()
)

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
	payload, err := detector.GetMatching(context.Background(), log, c, conf, nil)
	if err != nil {
		log.WithError(err).Fatal("unable to process request")
		return
	}
	rep := &detector.Reporter{
		Payload: payload,
	}
	val, err := rep.Generate()
	if err != nil {
		log.WithError(err).Fatal("failed to generate report")
	}
	log.Info(val)
}
