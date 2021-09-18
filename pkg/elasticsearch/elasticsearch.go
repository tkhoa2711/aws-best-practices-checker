package elasticsearch

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice/types"
)

func GetAllElasticsearchDomains() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		return nil, err
	}

	svc := elasticsearchservice.NewFromConfig(cfg)

	output, err := svc.ListDomainNames(context.TODO(), &elasticsearchservice.ListDomainNamesInput{})
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

func getDomainStatus(domainName *string) (*types.ElasticsearchDomainStatus, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		return nil, err
	}

	svc := elasticsearchservice.NewFromConfig(cfg)

	out, err := svc.DescribeElasticsearchDomain(
		context.TODO(),
		&elasticsearchservice.DescribeElasticsearchDomainInput{
			DomainName: domainName,
		},
	)
	if err != nil {
		return nil, err
	}

	return out.DomainStatus, nil
}

func getInstanceType(domainName *string) (*string, error) {
	esDomainStatus, err := getDomainStatus(domainName)
	if err != nil {
		return nil, err
	}

	return (*string)(&esDomainStatus.ElasticsearchClusterConfig.InstanceType), nil
}

func isProduction(domainName *string) (bool, error) {
	// TODO: make this list configurable
	patterns := []string{
		`(?i)prod`,
		`(?i)demo`,
	}

	for _, p := range patterns {
		matched, err := regexp.MatchString(p, *domainName)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
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
	if err != nil {
		return err
	}
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

func CheckDedicatedMasterNodes(domainName *string) error {
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