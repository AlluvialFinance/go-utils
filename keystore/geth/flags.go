//nolint:revive // package name intentionally reflects domain, not directory name
package gethkeystore

import (
	cmdutils "github.com/kilnfi/go-utils/cmd/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Flags(v *viper.Viper, f *pflag.FlagSet) {
	KeystorePathFlag(v, f)
	KeystorePasswordFlag(v, f)
}

func ConfigFromViper(v *viper.Viper) *Config {
	return &Config{
		Path:     GetKeystorePath(v),
		Password: GetKeystorePassword(v),
	}
}

const (
	keyStorePathFlag     = "keystore-path"
	keyStorePathViperKey = "keystore.path"
	keyStorePathEnv      = "KEYSTORE_PATH"
)

// KeystorePathFlag register flag for the path to the file keystore
func KeystorePathFlag(v *viper.Viper, f *pflag.FlagSet) {
	f.String(keyStorePathFlag, "", cmdutils.FlagDesc("Directory to store keys", keyStorePathEnv))

	if err := v.BindPFlag(keyStorePathViperKey, f.Lookup(keyStorePathFlag)); err != nil {
		panic(err)
	}
	if err := v.BindEnv(keyStorePathViperKey, keyStorePathEnv); err != nil {
		panic(err)
	}
}

func GetKeystorePath(v *viper.Viper) string {
	return v.GetString(keyStorePathViperKey)
}

const (
	keyStorePasswordFlag     = "keystore-password"
	keyStorePasswordViperKey = "keystore.password"
	keyStorePasswordEnv      = "KEYSTORE_PASSWORD"
)

// KeystorePasswordFlag register flag for the password used to encrypt keys in keystore
func KeystorePasswordFlag(v *viper.Viper, f *pflag.FlagSet) {
	f.String(keyStorePasswordFlag, "", cmdutils.FlagDesc("Password used to encrypt key files", keyStorePasswordEnv))

	if err := v.BindPFlag(keyStorePasswordViperKey, f.Lookup(keyStorePasswordFlag)); err != nil {
		panic(err)
	}
	if err := v.BindEnv(keyStorePasswordViperKey, keyStorePasswordEnv); err != nil {
		panic(err)
	}
}

func GetKeystorePassword(v *viper.Viper) string {
	return v.GetString(keyStorePasswordViperKey)
}
