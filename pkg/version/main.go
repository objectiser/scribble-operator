package version

import (
	"fmt"
	"runtime"

	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/spf13/viper"
)

var (
	version         string
	buildDate       string
	defaultScribble string
)

// Version holds this Operator's version as well as the version of some of the components it uses
type Version struct {
	Operator    string `json:"scribble-operator"`
	BuildDate   string `json:"build-date"`
	Scribble    string `json:"scribble-version"`
	Go          string `json:"go-version"`
	OperatorSdk string `json:"operator-sdk-version"`
}

// Get returns the Version object with the relevant information
func Get() Version {
	var scribble string
	if viper.IsSet("scribble-version") {
		scribble = viper.GetString("scribble-version")
	} else {
		scribble = defaultScribble
	}

	return Version{
		Operator:    version,
		BuildDate:   buildDate,
		Scribble:    scribble,
		Go:          runtime.Version(),
		OperatorSdk: sdkVersion.Version,
	}
}

func (v Version) String() string {
	return fmt.Sprintf(
		"Version(Operator='%v', BuildDate='%v', Scribble='%v', Go='%v', OperatorSDK='%v')",
		v.Operator,
		v.BuildDate,
		v.Scribble,
		v.Go,
		v.OperatorSdk,
	)
}

// DefaultScribble returns the default Scribble to use when no versions are specified via CLI or configuration
func DefaultScribble() string {
	return defaultScribble
}
