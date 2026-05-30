package bot

import (
	"sync"

	"aislaydbot/internal/i18n"
	"aislaydbot/internal/pptx"
)

type Step int

const (
	StepIdle Step = iota
	StepAwaitingTopic
	StepAwaitingCount
	StepAwaitingTheme
	StepGenerating
)

type Session struct {
	Lang  i18n.Lang
	Step  Step
	Topic string
	Count int
	Theme pptx.ThemeID
}

type Store struct {
	mu   sync.RWMutex
	data map[int64]*Session
}

func NewStore() *Store {
	return &Store{data: make(map[int64]*Session)}
}

func (s *Store) Get(chatID int64) *Session {
	s.mu.RLock()
	sess, ok := s.data[chatID]
	s.mu.RUnlock()
	if ok {
		return sess
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.data[chatID]; ok {
		return sess
	}
	sess = &Session{Lang: i18n.UZ, Step: StepIdle}
	s.data[chatID] = sess
	return sess
}

func (s *Store) Update(chatID int64, fn func(*Session)) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.data[chatID]
	if !ok {
		sess = &Session{Lang: i18n.UZ, Step: StepIdle}
		s.data[chatID] = sess
	}
	fn(sess)
	return sess
}

func (s *Store) Reset(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.data[chatID]; ok {
		sess.Step = StepIdle
		sess.Topic = ""
		sess.Count = 0
		sess.Theme = ""
	}
}
