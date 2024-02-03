package domain

import "time"

// User 领域对象，是DDD中的entity
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Nickname string
	AboutMe  string
	Birthday string
	Ctime    time.Time
}
