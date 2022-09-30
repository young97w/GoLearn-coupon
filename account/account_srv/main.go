package main

import (
	"account/account_srv/biz"
	"account/internal"
	"account/log"
	"account/proto/pb"
	"account/util"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	port := util.GenRandomPort2()
	//if err != nil {
	//	log.Logger.Panic(err.Error())
	//}
	host := internal.AppConf.AccountSrvConfig.Host
	addr := fmt.Sprintf("%s:%d", host, port)
	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &biz.AccountServer{})
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Logger.Panic(err.Error())
	}
	method := "grpc"
	srvName := internal.AppConf.AccountSrvConfig.SrvName
	randUUID := uuid.New().String()
	tags := internal.AppConf.AccountSrvConfig.Tags
	cr := internal.NewConsulRegistry(host, method, srvName, randUUID, tags, port)
	err = cr.Register(server)
	if err != nil {
		log.Logger.Panic(err.Error())
	}
	errChan := make(chan error)
	go func() {
		err = server.Serve(listen)
		errChan <- err
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err = <-errChan:
		panic(err)
		return
	case <-c:
		log.Logger.Info("account_srv is exiting")
		err = cr.Deregister()
		if err != nil {
			log.Logger.Panic(err.Error())
		}
		fmt.Println("exiting gracefully")
	}
}
