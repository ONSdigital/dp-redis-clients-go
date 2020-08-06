package redis

import (
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

const HealthyMessage = "redis is OK"

func (c *Client) Checker(state *health.CheckState) error {
	_, err := c.client.Ping().Result()
	if err != nil {
		// Generic error
		return state.Update(health.StatusCritical, err.Error(), 0)
	}
	// Success
	return state.Update(health.StatusOK, HealthyMessage, 0)
}
