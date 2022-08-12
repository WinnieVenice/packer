package crawl

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/WinnieVenice/packer/conf"
	"github.com/WinnieVenice/packer/idl"
)

var (
	once   sync.Once
	client idl.CrawlServiceClient
)

func Client() idl.CrawlServiceClient {
	once.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", viper.GetString("client.crawler.host"), viper.GetInt("client.crawler.port")),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(err)
		}
		client = idl.NewCrawlServiceClient(conn)
	})
	return client
}
