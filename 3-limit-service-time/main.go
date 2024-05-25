//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"context"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

const FreeLimit = 10 * time.Second
const PremiumLimit = 1<<63 - 1 // Maximum duration

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(ctx context.Context), u *User) bool {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()

	done := make(chan bool, 1)
	go func() {
		process(ctx)
		done <- true
	}()

	var limit time.Duration
	if u.IsPremium {
		limit = PremiumLimit
	} else {
		remaining := FreeLimit - time.Duration(u.TimeUsed)*time.Second
		if remaining > 0 {
			limit = remaining
		} else {
			return false
		}
	}

	select {
	case <-done:
		// The process finished before the time limit.
		u.TimeUsed += int64(time.Since(start).Seconds())
		return true
	case <-time.After(limit):
		// The process didn't finish before the time limit.
		u.TimeUsed += int64(limit.Seconds())
		cancel()
		return false
	}
}

func main() {
	RunMockServer()
}
