package config

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	EnvPrefix = "SIGNALVIEW"

	paramVerbose = "verbose"
	paramToken   = "token"
	paramRelm    = "sfx-relm"
	paramTimeout = "http-timeout"
	paramFormat  = "log-format"

	defaultRelm    = "us0"
	defaultTimeout = 16 * time.Second
	defaultFormat  = "text"
)

// Runtime is the config object to be passed around through different
// functions and packages
type Runtime struct {
	Relm  string
	Token string

	Verbose bool
	Output  string
	Timeout time.Duration
}

// NewDefault returns the default runtime config for the application
// and configures the global state
func NewDefault(v *viper.Viper) (*Runtime, error) {
	if v == nil {
		return nil, errors.New("no viper object configured")
	}
	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.SetTypeByDefaultValue(true)
	v.AutomaticEnv()

	rt := &Runtime{}
	cmd := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	cmd.BoolVar(&rt.Verbose, paramVerbose, false, "enable to show debub information within application")
	cmd.StringVar(&rt.Token, paramToken, "", "SignalFx token in order to query their endpoints")
	cmd.StringVar(&rt.Relm, paramRelm, defaultRelm, "SignalFx relm to query")
	cmd.StringVar(&rt.Output, paramFormat, defaultFormat, "defines what format to output the logs (alternative is json)")
	cmd.DurationVar(&rt.Timeout, paramTimeout, defaultTimeout, "default timeout for the internal http client")

	cmd.VisitAll(func(flag *pflag.Flag) {
		if err := v.BindPFlag(flag.Name, flag); err != nil {
			panic(err) // Should never happen
		}
	})

	err := cmd.Parse(os.Args[1:])
	// Ensure global state is configured
	defer configureLogging(rt)

	return rt, err
}

func configureLogging(rt *Runtime) {
	if rt.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if rt.Output == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func (rt *Runtime) ConfigureClient(c *http.Client) {
	c.Timeout = rt.Timeout
}
