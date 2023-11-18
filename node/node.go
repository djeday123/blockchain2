package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/djeday123/blockchain2/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	//"google.golang.org/grpc/credentials/insecure"
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
	version    string
	listenAddr string
	logger     *zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.Version

	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Development = true
	loggerConfig.EncoderConfig.TimeKey = ""
	//loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger, _ := loggerConfig.Build()
	return &Node{
		peers:   make(map[proto.NodeClient]*proto.Version),
		version: "bl-0.1",
		logger:  logger.Sugar(),
	}
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	// fmt.Printf("[%s] new peer connected: (%s) - height (%d) \n", n.listenAddr, v.ListenAddr, v.Height)

	n.logger.Debugw("new peer connected", "addr", v.ListenAddr, "height", v.Height)
	n.peers[c] = v
}

func (n *Node) deletePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) BootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			return err
		}

		v, err := c.Handshake(context.Background(), n.getVersion())
		if err != nil {
			n.logger.Error("handshake error:", err)
			continue
		}

		n.addPeer(c, v)
	}

	return nil
}

func (n *Node) Start(listenAddr string) error {
	n.listenAddr = listenAddr

	var (
		opts       = []grpc.ServerOption{}
		grpcServer = grpc.NewServer(opts...)
	)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	proto.RegisterNodeServer(grpcServer, n)

	n.logger.Infow("node started...", "port", n.listenAddr)
	return grpcServer.Serve(ln)
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	c, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.addPeer(c, v)

	return n.getVersion(), nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Println("received tx from:", peer)
	return &proto.Ack{}, nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "bl-0.1",
		Height:     0,
		ListenAddr: n.listenAddr,
	}
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	c, err := grpc.Dial(listenAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return proto.NewNodeClient(c), nil
}
