package main

import (
	"flag"
	"fmt"

	elasticsearch "github.com/tkhoa2711/aws-best-practices-checker/pkg/elasticsearch"
)

var (
	ignoreNonProdEnv bool
)

func main() {
	flag.BoolVar(&ignoreNonProdEnv, "ignore-non-prod", true, "Whether to ignore non-production environments, such as dev and staging")

	elasticsearchDomains, _ := elasticsearch.GetAllElasticsearchDomains()
	for _, domainName := range elasticsearchDomains {
		err := elasticsearch.CheckInstanceType(&domainName)
		if err != nil {
			fmt.Println(err)
		}

		err = elasticsearch.CheckDedicatedMasterNodes(&domainName)

		if err != nil {
			fmt.Println(err)
		}
	}
}
