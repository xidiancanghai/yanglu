package config

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestInitEnvConf(t *testing.T) {
	//InitEnvConf("./env.toml")
	InitEnvConf(GetEnvPath() + "/env.json")

	res := GetEnv()
	assert.Equal(t, res, EnvLocal)
	assert.Equal(t, GetLogPath(), "./logs")

}
