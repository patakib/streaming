package main

import (
	//"fmt"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var wg sync.WaitGroup

type fn func()

func initDb() {
	db, err := sql.Open("mysql", "streamuser:streampass@tcp(127.0.0.1:3306)/streaming")
	defer db.Close()
	if err != nil {
		log.Fatal()
	}

	cleanupUser := "DROP TABLE IF EXISTS user"
	if _, err := db.Exec(cleanupUser); err != nil {
		log.Fatal()
	}
	cleanupActivity := "DROP TABLE IF EXISTS activity"
	if _, err := db.Exec(cleanupActivity); err != nil {
		log.Fatal()
	}

	userCreationSql := "CREATE OR REPLACE TABLE user (id INT NOT NULL AUTO_INCREMENT, birth_year INT, location VARCHAR(100), gender VARCHAR(10), PRIMARY KEY(id))"
	if _, err := db.Exec(userCreationSql); err != nil {
		log.Fatal()
	}
	for i := 1; i <= 1000; i++ {
		createUser()
	}

	activityCreationSql := "CREATE OR REPLACE TABLE activity (id INT NOT NULL AUTO_INCREMENT, user_id INT, activity_type VARCHAR(50), intensity INT, duration INT, PRIMARY KEY(id), CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE)"
	if _, err := db.Exec(activityCreationSql); err != nil {
		log.Fatal()
	}
}

func createUser() {
	db, err := sql.Open("mysql", "streamuser:streampass@tcp(127.0.0.1:3306)/streaming")
	defer db.Close()
	if err != nil {
		log.Fatal()
	}
	rand.Seed(time.Now().UnixNano())
	birthYear := rand.Intn(50) + 1960

	rand.Seed(time.Now().UnixNano())
	locationInt := rand.Intn(30)
	locations := []string{
		"München", "Stuttgart", "Berlin", "Frankfurt am Main", "Freiburg am Breisgau", "Karlsruhe", "Köln",
		"Ulm", "Passau", "Leipzig", "Hamburg", "Düsseldorf", "Dortmund", "Essen", "Bremen", "Dresden", "Hannover",
		"Nürnberg", "Duisburg", "Bochum", "Wuppertal", "Bielefeld", "Bonn", "Münster", "Mannheim", "Augsburg",
		"Wiesbaden", "Mönchengladbach", "Gelsenkirchen", "Aachen",
	}
	selectedLocation := locations[locationInt]

	rand.Seed(time.Now().UnixNano())
	genderInt := rand.Intn(2)
	genders := []string{"male", "female"}
	selectedGender := genders[genderInt]
	if _, err := db.Exec("INSERT INTO user(birth_year, location, gender) VALUES (?, ?, ?)", birthYear, selectedLocation, selectedGender); err != nil {
		log.Fatal()
	}
	log.Println("New user added.")
}

func deleteUser() {
	db, err := sql.Open("mysql", "streamuser:streampass@tcp(127.0.0.1:3306)/streaming")
	defer db.Close()
	if err != nil {
		log.Fatal()
	}
	var count int
	rows, err := db.Query("SELECT count(id) FROM user")
	if err != nil {
		log.Fatal()
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal()
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal()
	}
	rand.Seed(time.Now().UnixNano())
	userId := rand.Intn(count)
	if _, err := db.Exec("DELETE FROM user WHERE id=?", userId+1); err != nil {
		log.Fatal()
	}
	log.Printf("Deleted user: %v", userId+1)
}

func createActivity() {
	db, err := sql.Open("mysql", "streamuser:streampass@tcp(127.0.0.1:3306)/streaming")
	defer db.Close()
	if err != nil {
		log.Fatal()
	}
	var count int
	rows, err := db.Query("SELECT count(id) FROM user")
	if err != nil {
		log.Fatal()
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal()
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal()
	}

	rand.Seed(time.Now().UnixNano())
	userId := rand.Intn(count)
	activityTypes := []string{
		"running", "trail running", "hiking", "mountain biking", "downhill", "road biking", "swimming", "walking",
		"martial art", "football", "aerobic", "yoga", "pilates", "stretching", "crossfit", "rowing", "circuit training",
		"trx", "indoor climbing", "outdoor climbing",
	}
	rand.Seed(time.Now().UnixNano())
	activityType := rand.Intn(20)
	selectedActivityType := activityTypes[activityType]
	rand.Seed(time.Now().UnixNano())
	intensity := rand.Intn(5)
	rand.Seed(time.Now().UnixNano())
	duration := rand.Intn(120) + 15

	if _, err := db.Exec("INSERT INTO activity(user_id, activity_type, intensity, duration) VALUES(?, ?, ?, ?)", userId+1, selectedActivityType, intensity, duration); err != nil {
		log.Fatal()
	}
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
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	initDb()
	wg.Add(111)
	for i := 1; i <= 100; i++ {
		go postMessage(createActivity, 100000000)
	}
	for j := 1; j <= 10; j++ {
		go postMessage(createUser, 1000000000)
	}
	wg.Wait()
}
