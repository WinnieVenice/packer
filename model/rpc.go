package model

import (
	"fmt"
	"packer/pb"

	"google.golang.org/grpc"
)

type RpcClient struct {
	CrawlServiceClient *pb.CrawlServiceClient
	CrawlServiceConn   *grpc.ClientConn
}

func (cli *RpcClient) NewCrawlServiceConn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", "localhost", 9851),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	cli.CrawlServiceConn = conn
	return conn, err
}
func (cli *RpcClient) GetCrawlServiceConn() (*grpc.ClientConn, error) {
	if cli.CrawlServiceConn == nil {
		_, err := cli.NewCrawlServiceConn()
		if err != nil {
			return nil, err
		}
	}
	return cli.CrawlServiceConn, nil
}
func (cli *RpcClient) NewCrawlServiceClient() (*pb.CrawlServiceClient, error) {
	_, err := cli.NewCrawlServiceConn()
	if err != nil {
		return nil, err
	}
	client := pb.NewCrawlServiceClient(cli.CrawlServiceConn)
	return &client, nil
}

func (cli *RpcClient) GetCrawlServiceClient() (*pb.CrawlServiceClient, error) {
	var err error
	if cli.CrawlServiceClient == nil {
		cli.CrawlServiceClient, err = cli.NewCrawlServiceClient()
	}
	if err != nil {
		return nil, err
	}

	return cli.CrawlServiceClient, nil
}
