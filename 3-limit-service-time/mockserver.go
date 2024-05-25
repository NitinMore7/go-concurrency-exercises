//////////////////////////////////////////////////////////////////////
//
// DO NOT EDIT THIS PART
// Your task is to edit `main.go`
//

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

// RunMockServer pretends to be a video processing service. It
// simulates user interacting with the Server.
func RunMockServer() {
	u1 := User{ID: 0, IsPremium: false}
	u2 := User{ID: 1, IsPremium: true}

	wg.Add(5)

	go createMockRequest(1, shortProcess, &u1) // Process 1should complete
	time.Sleep(1 * time.Second)

	go createMockRequest(2, longProcess, &u2) //Process 2  should complete
	time.Sleep(2 * time.Second)

	go createMockRequest(3, shortProcess, &u1) //Process 3 should be killed after 4 second
	time.Sleep(1 * time.Second)

	go createMockRequest(4, longProcess, &u1) // Process 4 should be killed
	go createMockRequest(5, shortProcess, &u2) // should complete

	wg.Wait()
}

func createMockRequest(pid int, fn func(ctx context.Context), u *User) {
	fmt.Println("UserID:", u.ID, "\tProcess", pid, "started.")
	res := HandleRequest(fn, u)

	if res {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "done.")
	} else {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "killed. (No quota left)")
	}

	wg.Done()
}

func shortProcess(ctx context.Context) {
	time.Sleep(6 * time.Second)
}

func longProcess(ctx context.Context) {
	time.Sleep(11 * time.Second)
}
