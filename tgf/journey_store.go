package tgf

import (
	"errors"
	"sync"
)

type Journey struct {
	ID              string
	TelegramChatID  int64
	ChatID          int64
	Command         string
	Next            int
	RawContext      []byte
	MessagesCleanup []int
}

var (
	JourneyNotFoundErr = errors.New("journey not found")
)

type JourneyStore interface {
	GetJourneyByChatID(chatID int64) (*Journey, error)
	CleanupChatJourney(chatID int64) error
	UpsertJourneyByTelegeramChatID(chatID int64, upsert Journey) (*Journey, error)
}

// TODO: maybe switch this for a rwmutex implementation?
// that way I can get type safety, not that it matters that much

type InMemJourneyStore struct {
	journeyMap sync.Map
}

func NewInMemJourneyStore() *InMemJourneyStore {
	return &InMemJourneyStore{}
}

func (s *InMemJourneyStore) GetJourneyByChat(chatID int64) (*Journey, error) {
	j, ok := s.journeyMap.Load(chatID)
	if !ok {
		return nil, JourneyNotFoundErr
	}

	journey, ok := j.(Journey)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	return &journey, nil
}

func (s *InMemJourneyStore) CleanupChatJourney(chatID int64) error {
	s.journeyMap.Delete(chatID)
	return nil
}

func (s *InMemJourneyStore) UpsertJourneyByTelegeramChatID(chatID int64, upsert Journey) (*Journey, error) {
	s.journeyMap.Store(chatID, upsert)
	j, ok := s.journeyMap.Load(chatID)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	journey, ok := j.(Journey)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	return &journey, nil
}
