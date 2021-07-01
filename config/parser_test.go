package config

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
)

func TestInitEnvConf(t *testing.T) {
	//InitEnvConf("./env.toml")
	InitEnvConf(GetEnvPath() + "/env.toml")

	res := GetEnv()
	fmt.Println("res = ", res)
	assert.Equal(t, res, EnvLocal)
	assert.Equal(t, GetLogPath(), "./logs")

}
