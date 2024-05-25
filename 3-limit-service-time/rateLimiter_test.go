package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	u1 := User{ID: 0, IsPremium: false}
	u2 := User{ID: 1, IsPremium: true}

	wg.Add(5)

	go testcreateMockRequest(t, 1, shortProcess, &u1) // Process 1should complete
	time.Sleep(1 * time.Second)

	go testcreateMockRequest(t, 2, longProcess, &u2) //Process 2  should complete
	time.Sleep(2 * time.Second)

	go testcreateMockRequest(t, 3, shortProcess, &u1) //Process 3 should be killed after 4 second
	time.Sleep(1 * time.Second)

	go testcreateMockRequest(t, 4, longProcess, &u1)  // Process 4 should be killed
	go testcreateMockRequest(t, 5, shortProcess, &u2) // should complete

	wg.Wait()
}

func testcreateMockRequest(t *testing.T, pid int, fn func(ctx context.Context), u *User) {
	//fmt.Println("UserID:", u.ID, "\tProcess", pid, "started.")
	res := HandleRequest(fn, u)

	if res {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "done.")
	} else {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "killed. (No quota left)")
	}

	// Add assertion
	if u.ID == 0 && pid == 4 {
		assert.False(t, res, "Expected process 4 for user 0 to be killed")
	} else {
		assert.True(t, res, "Expected process to be done")
	}

	wg.Done()
}
