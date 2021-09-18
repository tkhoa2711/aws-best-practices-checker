package main

import (
	"flag"
	"fmt"

	"github.com/tkhoa2711/aws-best-practices-checker/pkg/elasticsearch"
	"github.com/tkhoa2711/aws-best-practices-checker/pkg/s3"
)

var (
	ignoreNonProdEnv bool
)

func main() {
	flag.BoolVar(&ignoreNonProdEnv, "ignore-non-prod", true, "Whether to ignore non-production environments, such as dev and staging")

	var err error

	err = elasticsearch.Check()
	if err != nil {
		fmt.Println(err)
	}

	err = s3.Check()
	if err != nil {
		fmt.Println(err)
	}
}
