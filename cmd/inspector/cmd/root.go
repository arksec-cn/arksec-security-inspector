package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/client"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/controller"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/logging"
	"os"
	"os/signal"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "cloud-native-security-inspector",
	Long:  `cloud-native-security-inspector`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		setup()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run()
		return nil
	},
}

func init() {

}

func setup() {
	client.K8sInit()
	logging.GetLogger().Info().Msgf(logging.Green("k8s client init success"))
}

func run() {

	stopCh := make(chan struct{}, 1)
	defer close(stopCh)

	kubeClient := client.K8sClient()
	c := controller.NewController(kubeClient, stopCh)

	logging.GetLogger().Info().Msgf(logging.Green("Starting the Controller"))
	go c.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logging.GetLogger().Info().Msgf("Server exiting")
	return
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
