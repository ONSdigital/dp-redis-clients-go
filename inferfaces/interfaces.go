package inferfaces

//go:generate moq -out mock/mock_clientmanager.go -pkg mock . ClientManager

import "github.com/ONSdigital/dp-sessions-api/session"

//ClientManager - interface for redis client
type ClientManager interface {
	Set(s *session.Session) error
}
