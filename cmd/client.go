/*
Copyright © 2023 weixiaodong

*/
package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/weixiaodong/more/common/config"
	"github.com/weixiaodong/more/common/etcdv3"
	"github.com/weixiaodong/more/protos/pb"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// host, _ := cmd.Flags().GetString("host")

		ser := etcdv3.NewServiceDiscovery(config.GetDiscoveryEndpoints())
		defer ser.Close()
		ser.WatchService(config.GetDiscoveryServiceNamePrefix())
		for {
			select {
			case <-time.Tick(10 * time.Second):
				log.Println(ser.GetServices())

			}
		}
		// callSayHello(host)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.Flags().String("host", ":8080", "server host")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func callSayHello(host string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s %v", host, err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	// Contact the server and print out its response.
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "www"})
	// log.Print(r)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	{
		r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "www"})
		// log.Print(r)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.Message)

	}
}
