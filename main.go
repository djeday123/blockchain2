package main

import (
	"context"
	"log"
	"time"

	"github.com/djeday123/blockchain2/node"
	"github.com/djeday123/blockchain2/proto"
	"google.golang.org/grpc"
)

func main() {
	//node := node.NewNode()
	makeNode(":3000", []string{})
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"})
	time.Sleep(4 * time.Second)
	makeNode(":5000", []string{":4000"})

	// go func() {
	// 	for {
	// 		time.Sleep(4 * time.Second)
	// 	}
	// }()
	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listenAddr, bootstrapNodes)
	// if len(bootstrapNodes) > 0 {
	// 	if err := n.BootstrapNetwork(bootstrapNodes); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	return n
}

func makeTransaction() {

	client, err := grpc.Dial(":3000", grpc.WithInsecure())
	//client, err := grpc.Dial(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	version := &proto.Version{
		Version:    "bl-0.1",
		Height:     1,
		ListenAddr: ":4000",
	}

	_, err = c.Handshake(context.TODO(), version)
	if err != nil {
		log.Fatal(err)
	}
}
