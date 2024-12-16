package internal

import (
	"context"
	"encoding/json"
	"github.com/open-policy-agent/opa/config"
	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/plugins/bundle"
	"github.com/open-policy-agent/opa/plugins/discovery"
)

type ManualOverride struct {
}

func (m *ManualOverride) OnConfig(ctx context.Context, config *config.Config) (*config.Config, error) {
	config, err := discoveryPluginOverride(config)
	if err != nil {
		return config, err
	}
	return bundlePluginConfigOverride(config)
}

func (m *ManualOverride) OnConfigDiscovery(ctx context.Context, config *config.Config) (*config.Config, error) {
	return bundlePluginConfigOverride(config)
}

func discoveryPluginOverride(config *config.Config) (*config.Config, error) {
	var (
		discoveryConfig discovery.Config
		triggerManual   = plugins.TriggerManual
		message         []byte
	)

	if config.Discovery != nil {
		if err := json.Unmarshal(config.Discovery, &discoveryConfig); err == nil {
			discoveryConfig.Trigger = &triggerManual
			if message, err = json.Marshal(discoveryConfig); err == nil {
				config.Discovery = message
			} else {
				return config, err
			}
		} else {
			return config, err
		}
	}
	return config, nil
}

func bundlePluginConfigOverride(config *config.Config) (*config.Config, error) {
	var (
		bundlesConfig map[string]*bundle.Source
		manualTrigger = plugins.TriggerManual
		message       []byte
	)

	if config.Bundles != nil {
		if err := json.Unmarshal(config.Bundles, &bundlesConfig); err == nil {
			for _, bndlCfg := range bundlesConfig {
				bndlCfg.Trigger = &manualTrigger
			}
			if message, err = json.Marshal(bundlesConfig); err == nil {
				config.Bundles = message
			} else {
				return config, err
			}
		} else {
			return config, err
		}
	}
	return config, nil
}
