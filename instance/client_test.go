package instance

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestQueryInstances(t *testing.T) {
	account := os.Getenv("ACCOUNT")
	password := os.Getenv("PASSWORD")
	endpoint := os.Getenv("ENDPOINT")
	if endpoint == "" || account == "" || password == "" {
		t.Skip("Environment not set.")
	}
	client := NewInstanceClient(account, password, endpoint)
	instances, err := client.QueryInstances()
	if err != nil {
		t.Fatal("query error")
	}
	for _, instance := range instances {
		Json(instance)
	}
}

func Json(a any) {
	bytes, _ := json.Marshal(a)
	fmt.Println(string(bytes))
}
