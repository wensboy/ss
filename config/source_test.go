package config

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func MockLoadAll() {
	_global_config.SetDir("../spec")
	_global_config.Load(
		"config.json",
	)
	err := godotenv.Load("../spec/.env")
	if err != nil {
		panic(err)
	}
}

func Test_EnvSource(t *testing.T) {
	MockLoadAll()
	input_1, ok := Lookup[string](GEnvSource("TEST_ENV_SOURCE_1"))
	if !ok {
		t.Errorf("Expected to find env variable TEST_ENV_SOURCE_1, but it was not found")
	}
	if input_1 != "test_value_1" {
		t.Errorf("Expected value 'test_value_1' for env variable TEST_ENV_SOURCE_1, but got '%s'", input_1)
	}
	input_2, ok := Lookup[int](GEnvSource("test", "env", "source", "2"))
	if !ok {
		t.Errorf("Expected to find env variable TEST_ENV_SOURCE_2, but it was not found")
	}
	if input_2 != 8080 {
		t.Errorf("Expected value 8080 for env variable TEST_ENV_SOURCE_2, but got %d", input_2)
	}
	input_3, ok := Lookup[int64](GEnvSource("test.env.source.3"))
	if !ok {
		t.Errorf("Expected to find env variable TEST_ENV_SOURCE_3, but it was not found")
	}
	if input_3 != 123456789 {
		t.Errorf("Expected value 123456789 for env variable TEST_ENV_SOURCE_3, but got %d", input_3)
	}
	input_4, ok := Lookup[float64](GEnvSource("test.env.source.4"))
	if !ok {
		t.Errorf("Expected to find env variable TEST_ENV_SOURCE_4, but it was not found")
	}
	if input_4 != 3.14 {
		t.Errorf("Expected value 3.14 for env variable TEST_ENV_SOURCE_4, but got %f", input_4)
	}
	input_5, ok := Lookup[bool](GEnvSource("test.env.source.5"))
	if !ok {
		t.Errorf("Expected to find env variable TEST_ENV_SOURCE_5, but it was not found")
	}
	if input_5 != true {
		t.Errorf("Expected value true for env variable TEST_ENV_SOURCE_5, but got %v", input_5)
	}
}

func Test_ConfigSource(t *testing.T) {
	MockLoadAll()
	input_1, ok := Lookup[string](GConfigSource("config.test.config.1"))
	if !ok {
		t.Errorf("Expected to find config variable config.test.config.1, but it was not found")
	}
	if input_1 != "test-config-1" {
		t.Errorf("Expected value 'test-config-1' for config variable config.test.config.1, but got '%s'", input_1)
	}
	input_2, ok := Lookup[int](GConfigSource("config.test.config.2"))
	if !ok {
		t.Errorf("Expected to find config variable config.test.config.2, but it was not found")
	}
	if input_2 != 123 {
		t.Errorf("Expected value 123 for config variable config.test.config.2, but got %d", input_2)
	}
	input_3, ok := Lookup[int64](GConfigSource("config.test.config.3.1"))
	if !ok {
		t.Errorf("Expected to find config variable config.test.config.3.1, but it was not found")
	}
	if input_3 != 987654321 {
		t.Errorf("Expected value 987654321 for config variable config.test.config.3.1, but got %d", input_3)
	}
	input_4, ok := Lookup[float64](GConfigSource("config.test.config.4"))
	if !ok {
		t.Errorf("Expected to find config variable config.test.config.4, but it was not found")
	}
	if input_4 != 3.14 {
		t.Errorf("Expected value 3.14 for config variable config.test.config.4, but got %f", input_4)
	}
	input_5, ok := Lookup[bool](GConfigSource("config.test.config.5"))
	if !ok {
		t.Errorf("Expected to find config variable config.test.config.5, but it was not found")
	}
	if input_5 != false {
		t.Errorf("Expected value false for config variable config.test.config.5, but got %v", input_5)
	}
}

func Test_FlagSource(t *testing.T) {
	MockLoadAll()
	rootCmd := InitCommand("../spec/command.json")
	configCmd := GetCommand("ss.config")
	configCmd.Action = func(ctx context.Context, cmd *cli.Command) error {
		input_1, ok := Lookup[string](GFlagSource("ss.config.dir"))
		if !ok {
			t.Errorf("Expected to find flag variable ss.config.dir, but it was not found")
		}
		if input_1 != "/tmp/ss/config" {
			t.Errorf("Expected value '/tmp/ss/config' for flag variable ss.config.dir, but got '%s'", input_1)
		}
		input_2, ok := Lookup[string](GFlagSource("ss.config.ext"))
		if !ok {
			t.Errorf("Expected to find flag variable ss.config.ext, but it was not found")
		}
		if input_2 != "yaml" {
			t.Errorf("Expected value 'flag-value-2' for flag variable ss.config.ext, but got '%s'", input_2)
		}
		t.Logf("input1: %s, input2: %s", input_1, input_2)
		return nil
	}
	if err := rootCmd.Run(context.TODO(), []string{"ss", "config", "--ext", "yaml"}); err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
}
