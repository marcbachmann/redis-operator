package v1

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	maxNameLength = 48
)

// Validate set the values by default if not defined and checks if the values given are valid
func (r *RedisFailover) Validate() error {
	if len(r.Name) > maxNameLength {
		return fmt.Errorf("name length can't be higher than %d", maxNameLength)
	}

	initialRedisCustomConfig := defaultRedisCustomConfig
	if r.Bootstrapping() {
		if r.Spec.BootstrapNode.Host == "" {
			return errors.New("BootstrapNode must include a host when provided")
		}

		if r.Spec.BootstrapNode.Port == "" {
			r.Spec.BootstrapNode.Port = strconv.Itoa(defaultRedisPort)
		}
		initialRedisCustomConfig = bootstrappingRedisCustomConfig
	}
	redisCustomConfigExists := false
	for customConfig := range r.Spec.Redis.CustomConfig {
		for initialConfig := range initialRedisCustomConfig {
			if customConfig == initialConfig {
				redisCustomConfigExists = true
			}
		}
	}
	if !redisCustomConfigExists {
		r.Spec.Redis.CustomConfig = append(initialRedisCustomConfig, r.Spec.Redis.CustomConfig...)
	}

	if r.Spec.Redis.Image == "" {
		r.Spec.Redis.Image = defaultImage
	}

	if r.Spec.Sentinel.Image == "" {
		r.Spec.Sentinel.Image = defaultImage
	}

	if r.Spec.Redis.Replicas <= 0 {
		r.Spec.Redis.Replicas = defaultRedisNumber
	}

	if r.Spec.Redis.Port <= 0 {
		r.Spec.Redis.Port = defaultRedisPort
	}

	if r.Spec.Sentinel.Replicas <= 0 {
		r.Spec.Sentinel.Replicas = defaultSentinelNumber
	}

	if r.Spec.Redis.Exporter.Image == "" {
		r.Spec.Redis.Exporter.Image = defaultExporterImage
	}

	if r.Spec.Sentinel.Exporter.Image == "" {
		r.Spec.Sentinel.Exporter.Image = defaultSentinelExporterImage
	}

	if len(r.Spec.Sentinel.CustomConfig) == 0 {
		r.Spec.Sentinel.CustomConfig = defaultSentinelCustomConfig
	}

	return nil
}
