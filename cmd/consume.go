/*
Copyright © 2023 weixiaodong

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/weixiaodong/more/common/rabbitmq"
)

// consumeCmd represents the client command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rabbitmq.StartRabbitMQConsumer(handler)
		select {}
	},
}

func handler(msg []byte) error {
	log.Println(string(msg))
	return nil
}
func init() {
	rootCmd.AddCommand(consumeCmd)
}
