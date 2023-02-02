/*
Copyright © 2023 Carlos Marcelo crodriguezde@microsoft.com
*/
package cmd

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	"github.com/crodriguezde/rtdashs/pkg/consumer"
	kafkaPlayloads "github.com/crodriguezde/rtdashs/pkg/kafkaPayloads"
	"github.com/crodriguezde/rtdashs/pkg/server"
	"github.com/crodriguezde/rtdashs/pkg/sink"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rtdashs",
	Short: "real time dashboards demo",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, err := cmd.Flags().GetString("addr")
		if err != nil {
			return err
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}

		broker, err := cmd.Flags().GetString("broker")
		if err != nil {
			return err
		}

		topic, err := cmd.Flags().GetString("topic")
		if err != nil {
			return err
		}

		version, err := cmd.Flags().GetString("version")
		if err != nil {
			return err
		}

		log.Printf("listening on http://%s", addr)

		// Initialize the notify channel
		send := make(chan *kafkaPlayloads.Cpu)
		ctx := context.Background()

		// Start listening to kafka topic
		if verbose {
			sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
		}

		consumer, err := consumer.StartSync(broker, topic, version, send)
		if err != nil {
			return err
		}

		defer consumer.Close()

		collector := sink.NewServer(ctx, send)
		collector.Start()

		// disable generator
		// generator.RunCpuGenerator(ctx, send)

		handler := server.NewHandler(ctx, collector.GetChannel())

		err = http.ListenAndServe(addr, handler)
		if err != nil {
			return err
		}
		<-ctx.Done()

		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("addr", "a", "localhost:3000", "bind address")
	rootCmd.Flags().StringP("broker", "b", "kafka.default.svc.cluster.local:9092", "Kafka bootstrap brokers, separated list")
	rootCmd.Flags().StringP("topic", "t", "telemetry-synthetics-cpu", "Topic to listen for cpu metrics")
	rootCmd.Flags().StringP("version", "s", "2.1.1", "Kafka cluster version")
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose logging")
}
