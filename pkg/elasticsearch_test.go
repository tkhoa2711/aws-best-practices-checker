package elasticsearch

import (
	"testing"
)

func TestFetchAllElasticsearchDomains(t *testing.T) {
	_, err := GetAllElasticsearchDomains()
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
