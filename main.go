package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
)

var wg sync.WaitGroup
type fn func()

func initDb() {
	//db, err := sql.Open("mysql", 
}

func createUser() {
	rand.Seed(time.Now().UnixNano())
	birthYear := rand.Intn(50) + 1960
	location := rand.Intn(30)
	gender := rand.Intn(2)
	fmt.Printf("NEW USER: Birth year: %v, location: %v, gender: %v\n", birthYear, location, gender)
}

func deleteUser() {
	fmt.Println("DELETED USER !!!\n")
}

func createActivity() {
	rand.Seed(time.Now().UnixNano())
	//get user from database
	activityType := rand.Intn(20)
	intensity := rand.Intn(5)
	duration := rand.Intn(120) + 15
	fmt.Printf("NEW EVENT: Activity: %v, Intensity: %v, Duration: %v\n", activityType, intensity, duration)
}

func postMessage(event fn, randomRange int) {
	defer wg.Done()
	for {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(randomRange)
		go event()
		time.Sleep(time.Duration(r) * time.Microsecond)
	}
}

func main() {
	wg.Add(111)
	for i := 1; i <= 100; i++ {
		go postMessage(createActivity, 100000000)
	}
	for j := 1; j <= 10; j++ {
		go postMessage(createUser, 1000000000)
	}
	go postMessage(deleteUser, 1000000000)
	wg.Wait()
}
