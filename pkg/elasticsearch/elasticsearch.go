package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"regexp"
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

func getInstanceType(domainName *string) (*string, error) {
	// TODO: can we extract this session logic out?
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := elasticsearchservice.New(sess)

	esDomainStatus, err := svc.DescribeElasticsearchDomain(&elasticsearchservice.DescribeElasticsearchDomainInput{
		DomainName: domainName,
	})
	if err != nil {
		return nil, err
	}

	return esDomainStatus.DomainStatus.ElasticsearchClusterConfig.InstanceType, nil
}

func isProduction(domainName *string) (bool, error) {
	// TODO: make this list configurable
	matchedProd, err := regexp.MatchString(`(?i)prod`, *domainName)
	matchedDemo, err := regexp.MatchString(`(?i)demo`, *domainName)

	if err != nil {
		return false, err
	}

	return matchedProd || matchedDemo, nil
}

func CheckInstanceType(domainName *string) error {
	// https://docs.aws.amazon.com/elasticsearch-service/latest/developerguide/aes-bp.html
	// Don't use T2 or t3.small instances for production domains; they can become
	// unstable under sustained heavy load. t3.medium instances are an option for
	// small production workloads (both as data nodes and dedicated master nodes)
	instanceType, err := getInstanceType(domainName)
	if err != nil {
		return err
	}

	isProduction, err := isProduction(domainName)
	if err != nil {
		return err
	}

	usingT2, err := regexp.MatchString(`t2`, *instanceType)
	if isProduction && usingT2 {
		fmt.Println(*domainName, "- You should not use `t2` instances for production")
	}

	usingT3Small, err := regexp.MatchString(`t3\.small`, *instanceType)
	if isProduction && usingT3Small {
		fmt.Println(*domainName, "- You should not use `t3.small` instances for production")
	}

	return nil
}

func main() {
	elasticsearchDomains, _ := getAllElasticsearchDomains()
	for _, domainName := range elasticsearchDomains {
		err := CheckInstanceType(&domainName)
		if err != nil {
			fmt.Println(err)
		}
	}
}
