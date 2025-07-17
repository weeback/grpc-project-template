package jwt

import (
	"github.com/google/uuid"
	"time"
)

func NewOption() *Option {
	return &Option{
		sessionId: uuid.NewString(),
		liveTime:  90 * time.Second,
	}
}

type Option struct {
	sessionId, userId string
	liveTime          time.Duration
}

func (src *Option) SetLiveTime(d time.Duration) *Option {
	dst := *src
	dst.liveTime = d
	return &dst
}

func (src *Option) LiveTime() time.Duration {
	return src.liveTime
}

func (src *Option) SetSessionId(sessionId string) *Option {
	dst := *src
	dst.sessionId = sessionId
	return &dst
}

func (src *Option) SessionId() string {
	return src.sessionId
}

func (src *Option) SetUserId(userId string) *Option {
	dst := *src
	dst.userId = userId
	return &dst
}

func (src *Option) UserId() string {
	return src.userId
}
