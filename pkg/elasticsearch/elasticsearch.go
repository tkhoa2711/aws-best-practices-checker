package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
)

func getAllElasticsearchDomains() ([]string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := elasticsearchservice.New(sess)

	output, err := svc.ListDomainNames(&elasticsearchservice.ListDomainNamesInput{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	domainNames := []string{}
	for _, item := range output.DomainNames {
		domainNames = append(domainNames, *item.DomainName)
	}
	fmt.Println("All Elasticsearch domains:", domainNames)

	return domainNames, nil
}

func main() {
	getAllElasticsearchDomains()
}
