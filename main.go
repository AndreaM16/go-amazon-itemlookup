package main

import (
	"github.com/andream16/go-amazon-itemlookup/configuration"
	"github.com/andream16/go-amazon-itemlookup/crawler"
	"fmt"
)

func main() {
	fmt.Println("Getting configuration . . .")
	conf := configuration.InitConfiguration()
	fmt.Println(fmt.Sprintf("Successfully got configuration. Amazon data: AssociateTag=%s AWSAccessKeyId=%s AWSSecretKey=%s; Remote: API=%s PORT=%s",
		conf.Credentials.AssociateTag, conf.Credentials.AWSAccessKeyId, conf.Credentials.AWSSecretKey, conf.Api.Host, conf.Api.Port))
	fmt.Println("Now starting the crawler . . .")
	crawler.Crawl(conf)
}


