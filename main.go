package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

type fn func(user, pass, port, database string)

func initDb(user, pass, port, database string) {
	conn := fmt.Sprintf("%s:%s@tcp(db:%s)/%s", user, pass, port, database)
	db, err := sql.Open("mysql", conn)
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
	for i := 1; i <= 10000; i++ {
		createUser(user, pass, port, database)
	}

	activityCreationSql := "CREATE OR REPLACE TABLE activity (id INT NOT NULL AUTO_INCREMENT, user_id INT, activity_type VARCHAR(50), intensity INT, duration INT, datetime DATETIME, PRIMARY KEY(id), CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE)"
	if _, err := db.Exec(activityCreationSql); err != nil {
		log.Fatal()
	}
}

func createUser(user, pass, port, database string) {
	conn := fmt.Sprintf("%s:%s@tcp(db:%s)/%s", user, pass, port, database)
	db, err := sql.Open("mysql", conn)
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

func deleteUser(user, pass, port, database string) {
	conn := fmt.Sprintf("%s:%s@tcp(db:%s)/%s", user, pass, port, database)
	db, err := sql.Open("mysql", conn)
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

func createActivity(user, pass, port, database string) {
	conn := fmt.Sprintf("%s:%s@tcp(db:%s)/%s", user, pass, port, database)
	db, err := sql.Open("mysql", conn)
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

	datetime := time.Now().Format("2006-01-02 15:04:05")

	if _, err := db.Exec("INSERT INTO activity(user_id, activity_type, intensity, duration, datetime) VALUES(?, ?, ?, ?, ?)", userId+1, selectedActivityType, intensity, duration, datetime); err != nil {
		log.Fatal()
	}
}

func postMessage(event fn, randomRange int, user, pass, port, database string) {
	defer wg.Done()
	for {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(randomRange)
		go event(user, pass, port, database)
		time.Sleep(time.Duration(r) * time.Microsecond)
	}
}

func main() {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	err2 := godotenv.Load("local.env")
	if err2 != nil {
		log.Fatal()
	}
	db_user := os.Getenv("MYSQL_USER")
	db_pass := os.Getenv("MYSQL_PASS")
	db_port := os.Getenv("MYSQL_PORT")
	db_database := os.Getenv("MYSQL_DATABASE")

	initDb(db_user, db_pass, db_port, db_database)
	wg.Add(160)
	for i := 1; i <= 50; i++ {
		go postMessage(createActivity, 100000000, db_user, db_pass, db_port, db_database)
	}
	for j := 1; j <= 100; j++ {
		go postMessage(createUser, 100000000, db_user, db_pass, db_port, db_database)
	}
	wg.Wait()
}
