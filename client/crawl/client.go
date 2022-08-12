package crawl

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/WinnieVenice/packer/conf"
	"github.com/WinnieVenice/packer/idl"
)

var (
	once   sync.Once
	client idl.CrawlServiceClient
)

func Client() idl.CrawlServiceClient {
	once.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", conf.V.GetString("client.crawler.host"), conf.V.GetInt("client.crawler.port")),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(err)
		}
		client = idl.NewCrawlServiceClient(conn)
	})
	return client
}
