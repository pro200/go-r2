package r2_test

import (
	"fmt"
	"testing"

	"github.com/pro200/go-env"
	"github.com/pro200/go-r2"
)

func TestR2(t *testing.T) {
	err := env.Init()
	if err != nil {
		t.Fatal(err)
	}

	accountId, err := env.Get("ACCOUNT_ID")
	if err != nil {
		t.Fatal(err)
	}
	accessKeyId, err := env.Get("ACCESS_KEY_ID")
	if err != nil {
		t.Fatal(err)
	}
	secretAccessKey, err := env.Get("SECRET_ACCESS_KEY")
	if err != nil {
		t.Fatal(err)
	}

	r2.Init(r2.Config{
		AccountId:       accountId,
		AccessKeyID:     accessKeyId,
		SecretAccessKey: secretAccessKey,
	})

	r, _, err := r2.List("dev", "", 10)
	fmt.Println(r)

	if err != nil {
		t.Fatal(err)
	}

}
