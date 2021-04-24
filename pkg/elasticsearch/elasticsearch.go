package main

import (
	"fmt"
	"regexp"

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

func getDomainStatus(domainName *string) (*elasticsearchservice.ElasticsearchDomainStatus, error) {
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

	return esDomainStatus.DomainStatus, nil
}

func getInstanceType(domainName *string) (*string, error) {
	esDomainStatus, err := getDomainStatus(domainName)
	if err != nil {
		return nil, err
	}

	return esDomainStatus.ElasticsearchClusterConfig.InstanceType, nil
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
	if err != nil {
		return err
	}

	return nil
}

func checkDedicatedMasterNodes(domainName *string) error {
	// https://docs.aws.amazon.com/elasticsearch-service/latest/developerguide/es-managedomains-dedicatedmasternodes.html
	isProduction, err := isProduction(domainName)
	if err != nil {
		return err
	}
	if !isProduction {
		return nil
	}

	esDomainStatus, err := getDomainStatus(domainName)
	if err != nil {
		return err
	}

	dedicatedMasterCount := esDomainStatus.ElasticsearchClusterConfig.DedicatedMasterCount
	if dedicatedMasterCount == nil {
		fmt.Println(*domainName, "- has no dedicated master node")
		return nil
	}

	if *dedicatedMasterCount < 3 {
		fmt.Println(*domainName, "- has less than 3 dedicated master nodes")
	}
	if *dedicatedMasterCount > 0 && *dedicatedMasterCount%2 == 0 {
		fmt.Println(*domainName, "- has an even number of dedicated master nodes")
	}

	return nil
}

func main() {
	elasticsearchDomains, _ := getAllElasticsearchDomains()
	for _, domainName := range elasticsearchDomains {
		err := CheckInstanceType(&domainName)
		err = checkDedicatedMasterNodes(&domainName)

		if err != nil {
			fmt.Println(err)
		}
	}
}
