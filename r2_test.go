package r2_test

import (
	"fmt"
	"testing"

	"github.com/pro200/go-r2"
)

func TestR2(t *testing.T) {
	// read only token: R2 Dev - 4fe41a62e8
	r2.Init(r2.Config{
		AccountId:       "882b4dd326c8b297dd43e5df32f08b6d",
		AccessKeyID:     "4fe41a62e8e71d9475c6ad46a3b8bb3c",
		SecretAccessKey: "33c2eb7deac4d187e1d3d2a51cf0b6620a66107cec95a6d9abdcaed059e2dacd",
	})

	r, _, err := r2.List("dev", "", 10)
	fmt.Println(r)

	if err != nil {
		t.Error(err)
	}
}
