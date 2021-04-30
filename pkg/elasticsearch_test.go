package elasticsearch

import (
	"context"
	"testing"
	"time"
)

func TestFetAllElasticsearchDomains(t *testing.T) {
	_, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	// sess := integration.SessionWithDefaultRegion("ap-southeast-2")
	// svc := elasticsearchservice.New(sess)

	_, err := GetAllElasticsearchDomains()
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
