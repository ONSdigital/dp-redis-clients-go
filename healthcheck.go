package dp_redis

import (
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

const HealthyMessage = "redis is OK"

func (c *Redis) Checker(state *health.CheckState) error {
	_, err := c.client.Ping().Result()
	if err != nil {
		// Generic error
		state.Update(health.StatusCritical, err.Error(), 0)
		return err
	}
	// Success
	state.Update(health.StatusOK, HealthyMessage, 0)
	return nil
}
