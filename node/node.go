package node

import (
	"context"
	"fmt"

	"github.com/djeday123/blockchain2/proto"
	"google.golang.org/grpc/peer"
)

// type Server struct {
// 	listenAddr string
// 	ln         net.Listener
// }

// func New(listenAddr string) (*Server, error) {
// 	ln, err := net.Listen("tpc", listenAddr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Server{
// 		listenAddr: listenAddr,
// 		ln:         ln,
// 	}, nil
// }

type Node struct {
	//peers map[net.Addr]*grpc.ClientConn
	version string
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		version: "bl-0.1",
	}
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	ourVersion := &proto.Version{
		Version: n.version,
		Height:  100,
	}

	p, _ := peer.FromContext(ctx)

	fmt.Printf("received version from %s: %+v\n", v, p.Addr)

	return ourVersion, nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Println("received tx from:", peer)
	return &proto.Ack{}, nil
}
