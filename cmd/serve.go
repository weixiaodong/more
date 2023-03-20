/*
Copyright © 2023 weixiaodong

*/
package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/weixiaodong/more/common/config"
	"github.com/weixiaodong/more/common/etcdv3"
	"github.com/weixiaodong/more/middleware"
	"github.com/weixiaodong/more/middleware/ratelimit"
	"github.com/weixiaodong/more/middleware/recovery"
	"github.com/weixiaodong/more/protos/pb"
	"github.com/weixiaodong/more/service"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		address := config.GeGrpcServiceAddr()
		startGrpcServer(address)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().String("address", ":8080", "server address")
}

func startGrpcServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		// grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		// 	grpc_ctxtags.StreamServerInterceptor(),
		// 	grpc_opentracing.StreamServerInterceptor(),
		// 	grpc_prometheus.StreamServerInterceptor,
		// 	grpc_zap.StreamServerInterceptor(zapLogger),
		// 	grpc_auth.StreamServerInterceptor(myAuthFunction),
		// 	grpc_recovery.StreamServerInterceptor(),
		// )),
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(
			// grpc_ctxtags.UnaryServerInterceptor(),
			// grpc_opentracing.UnaryServerInterceptor(),
			// grpc_prometheus.UnaryServerInterceptor,
			// grpc_zap.UnaryServerInterceptor(zapLogger),
			// grpc_auth.UnaryServerInterceptor(myAuthFunction),
			recovery.UnaryServerInterceptor(),
			ratelimit.UnaryServerInterceptor(),
		)))
	pb.RegisterGreeterServer(s, &service.HelloService{})

	//把服务注册到etcd
	ser, err := etcdv3.NewServiceRegister(
		config.GetDiscoveryEndpoints(),
		config.GetDiscoveryServiceNamePrefix(),
		address,
		config.GetDiscoveryTimeout(),
	)
	if err != nil {
		log.Fatalf("register service err: %v", err)
	}
	defer ser.Close()
	go ser.ListenLeaseRespChan()

	log.Println("启动服务", address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
