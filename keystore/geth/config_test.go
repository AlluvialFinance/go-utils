package gethkeystore

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestKeystorePathFlag(t *testing.T) {
	v := viper.New()
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Register the keystore path flag
	KeystorePathFlag(v, f)

	// Check if the flag is registered
	flag := f.Lookup("keystore-path")
	assert.NotNil(t, flag, "keystore-path flag should be registered")

	// Simulate setting the flag through command-line args
	if err := f.Set("keystore-path", "/test/keystore"); err != nil {
		t.Fatal(err)
	}

	// Check if the value is correctly bound to Viper
	assert.Equal(t, "/test/keystore", v.GetString("keystore.path"), "keystore.path should be set correctly from flags")

	// Test getting the value using GetKeystorePath function
	assert.Equal(t, "/test/keystore", GetKeystorePath(v), "GetKeystorePath should return the correct value")
}

func TestKeystorePasswordFlag(t *testing.T) {
	v := viper.New()
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Register the keystore password flag
	KeystorePasswordFlag(v, f)

	// Check if the flag is registered
	flag := f.Lookup("keystore-password")
	assert.NotNil(t, flag, "keystore-password flag should be registered")

	// Simulate setting the flag through command-line args
	if err := f.Set("keystore-password", "supersecret"); err != nil {
		t.Fatal(err)
	}
	// Check if the value is correctly bound to Viper
	assert.Equal(t, "supersecret", v.GetString("keystore.password"), "keystore.password should be set correctly from flags")

	// Test getting the value using GetKeystorePassword function
	assert.Equal(t, "supersecret", GetKeystorePassword(v), "GetKeystorePassword should return the correct value")
}

func TestConfigFromViper(t *testing.T) {
	v := viper.New()
	v.Set("keystore.path", "/test/keystore")
	v.Set("keystore.password", "supersecret")

	// Test if ConfigFromViper retrieves the correct values
	config := ConfigFromViper(v)
	assert.Equal(t, "/test/keystore", config.Path, "Config.Path should match the Viper value")
	assert.Equal(t, "supersecret", config.Password, "Config.Password should match the Viper value")
}

func TestKeystoreEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("KEYSTORE_PATH", "/env/keystore")
	os.Setenv("KEYSTORE_PASSWORD", "envpassword")
	defer os.Unsetenv("KEYSTORE_PATH")
	defer os.Unsetenv("KEYSTORE_PASSWORD")

	v := viper.New()
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Register the flags
	Flags(v, f)

	// Verify if Viper picks up the environment variables
	assert.Equal(t, "/env/keystore", v.GetString("keystore.path"), "keystore.path should be set from environment variable")
	assert.Equal(t, "envpassword", v.GetString("keystore.password"), "keystore.password should be set from environment variable")
}
