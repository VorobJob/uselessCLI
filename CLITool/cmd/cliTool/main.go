package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"CLITool/internal/repository/postgres"
	"CLITool/internal/service/brain"
	"CLITool/internal/storage"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Record struct {
	Name string
	DOB  string
	Sex  string
}

type MaleFNames struct {
	Record
	isFMale bool
}

type Repository struct {
	DB *gorm.DB
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.New(db)

	// fmt.Println("Usage: myapp [command] [flags]")

	brain := brain.New(repo)

	// считываешь флаг
	createTable := flag.Bool("1", false, "create a new table")
	var specialFlag string
	flag.StringVar(&specialFlag, "2", "", "a special flag")
	listUnique := flag.Bool("3", false, "list unique full name + date of birth records")
	autoFill := flag.Bool("4", false, "Autofill db with 1000000 records")
	getRecordsSimp := flag.Bool("5", false, "get all male users whose names start with F")
	getRecordsAdvanced := flag.Bool("6", false, "get all male users whose names start with F")

	flag.Parse()

	if *createTable {
		if err := brain.Start(1, nil); err != nil {
			log.Printf("error: %v", err)
		}
	}
	if specialFlag != "" {
		args := flag.Args()
		// fmt.Println(args)
		if len(args) != 5 {
			fmt.Println("Usage: command -2 {any word} lasName midName FirstName dateOfBirth Sex")
		}

		if err := brain.Start(2, args); err != nil {
			log.Printf("error: %v", err)
		}
	}
	if *listUnique {
		if err := brain.Start(3, nil); err != nil {
			log.Printf("error: %v", err)

			// fmt.Errorf("error: %w", err)
		}
	}
	if *autoFill {
		if err := brain.Start(4, nil); err != nil {
			log.Printf("error: %v", err)
		}
	}
	if *getRecordsSimp {
		if err := brain.Start(5, nil); err != nil {
			log.Printf("error: %v", err)
		}
	}
	if *getRecordsAdvanced {
		if err := brain.Start(6, nil); err != nil {
			log.Printf("error: %v", err)
		}
	}
}
