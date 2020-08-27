package dpredis

import (
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

const HealthyMessage = "redis is OK"

func (c *Client) Checker(state *health.CheckState) error {
	err := c.Ping()
	if err != nil {
		// Generic error
		return state.Update(health.StatusCritical, err.Error(), 0)
	}
	// Success
	return state.Update(health.StatusOK, HealthyMessage, 0)
}
