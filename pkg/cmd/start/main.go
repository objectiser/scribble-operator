package start

import (
	"runtime"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/objectiser/scribble-operator/pkg/apis"
	"github.com/objectiser/scribble-operator/pkg/controller"
	"github.com/objectiser/scribble-operator/pkg/version"
)

// NewStartCommand starts the Scribble Operator
func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new Scribble Operator",
		Long:  "Starts a new Scribble Operator",
		Run: func(cmd *cobra.Command, args []string) {
			start(cmd, args)
		},
	}

	cmd.Flags().String("scribble-version", version.DefaultScribble(), "The Scribble version to use")
	viper.BindPFlag("scribble-version", cmd.Flags().Lookup("scribble-version"))

	cmd.Flags().String("scribble-monitor-image", "objectiser/scribble-monitor", "The Docker image for the Scribble monitor")
	viper.BindPFlag("scribble-monitor-image", cmd.Flags().Lookup("scribble-monitor-image"))

	return cmd
}

func start(cmd *cobra.Command, args []string) {
	log.WithFields(log.Fields{
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"version": runtime.Version(),
	}).Print("Go")

	log.WithField("version", version.Get().OperatorSdk).Print("operator-sdk")

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Fatalf("failed to get watch namespace: %v", err)
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{Namespace: namespace})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Fatal(err)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Fatal(err)
	}

	log.Print("Starting the Cmd.")

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
