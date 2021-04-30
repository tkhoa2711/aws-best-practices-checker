package main

import (
	"fmt"

	elasticsearch "github.com/tkhoa2711/aws-best-practices-checker/pkg"
)

func main() {
	elasticsearchDomains, _ := elasticsearch.GetAllElasticsearchDomains()
	for _, domainName := range elasticsearchDomains {
		err := elasticsearch.CheckInstanceType(&domainName)
		err = elasticsearch.CheckDedicatedMasterNodes(&domainName)

		if err != nil {
			fmt.Println(err)
		}
	}
}
