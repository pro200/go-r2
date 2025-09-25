package r2_test

import (
	"fmt"
	"testing"

	"github.com/pro200/go-env"
	"github.com/pro200/go-r2"
)

/* .config.env
ACCOUNT_ID:        string
ACCESS_KEY_ID:     string
SECRET_ACCESS_KEY: string
*/

func TestR2(t *testing.T) {
	config, err := env.New()
	if err != nil {
		t.Error(err)
	}

	storage, err := r2.New(r2.Config{
		AccountId:       config.Get("ACCOUNT_ID"),
		AccessKeyID:     config.Get("ACCESS_KEY_ID"),
		SecretAccessKey: config.Get("SECRET_ACCESS_KEY"),
	})
	if err != nil {
		t.Error("R2 연결 실패:", err)
	}

	result, _, err := storage.List("dev", "", 10)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)
}
