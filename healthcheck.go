package redis

import (
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

const HealthyMessage = "redis is OK"

func (c *Client) Checker(state *health.CheckState) error {
	_, err := c.client.Ping().Result()
	if err != nil {
		// Generic error
		_ = state.Update(health.StatusCritical, err.Error(), 0)
		return err
	}
	// Success
	_ = state.Update(health.StatusOK, HealthyMessage, 0)
	return nil
}
