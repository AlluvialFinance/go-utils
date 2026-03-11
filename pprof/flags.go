package pprof

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	enabledFlag     = "pprof-enabled"
	enabledViperKey = "pprof.enabled"
	enabledEnv      = "PPROF_ENABLED"

	addressFlag     = "pprof-address"
	addressViperKey = "pprof.address"
	addressEnv      = "PPROF_ADDRESS"
)

// Flags registers pprof-specific flags.
func Flags(v *viper.Viper, f *pflag.FlagSet) {
	f.Bool(enabledFlag, false, "Enable pprof debug server on separate port (dev only)")
	_ = v.BindPFlag(enabledViperKey, f.Lookup(enabledFlag))
	_ = v.BindEnv(enabledViperKey, enabledEnv)

	f.String(addressFlag, DefaultAddress, "pprof server address (default: localhost-only)")
	_ = v.BindPFlag(addressViperKey, f.Lookup(addressFlag))
	_ = v.BindEnv(addressViperKey, addressEnv)
}

// ConfigFromViper constructs Config from viper.
func ConfigFromViper(v *viper.Viper) *Config {
	cfg := &Config{
		Enabled: v.GetBool(enabledViperKey),
		Address: v.GetString(addressViperKey),
	}
	cfg.SetDefaults()
	return cfg
}
