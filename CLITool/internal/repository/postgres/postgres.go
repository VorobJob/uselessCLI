package postgres

import (
	"CLITool/internal/models"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

type postgres struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *postgres {
	return &postgres{
		DB: DB,
	}
}

// Creates table only once
func (p *postgres) CreateTable() error {
	// fmt.Println(1)
	if err := p.DB.AutoMigrate(models.Records{}); err != nil {
		return fmt.Errorf("migrate db: %w", err)
	}
	return nil
}

func (p *postgres) CreateRecord(user models.User) error {
	record := models.Records{User: user}

	if err := p.DB.Create(&record).Error; err != nil {
		return fmt.Errorf("create: %w", err)
	}
	if err := p.CreateIndex(); err != nil {
		return fmt.Errorf("indexing error: %w", err)
	}
	// fmt.Println("Index created successfully")
	return nil
}

// не читает первый аргумент после флага, приходится писать запрос типа
// go run main.go -2 create Fname Putin Biden 03/02/1956 Male
// в целом работает

func (p *postgres) CreateIndex() error {
	var indexExists bool

	// Check if the index already exists
	err := p.DB.Raw("SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'male_f_index')").Row().Scan(&indexExists)
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}

	if indexExists {
		fmt.Println("Index already exists")
		return nil
	}

	// Create the index
	if err := p.DB.Exec(`CREATE INDEX male_f_index ON records ("sex", "name") WHERE "sex" = TRUE AND "name" ILIKE 'F%';`).Error; err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	fmt.Println("Index created successfully")
	return nil
}

// func (p *postgres) CreateIndex() error {
// 	if err := p.DB.Exec(`DROP INDEX IF EXISTS male_f_index;`).Error; err != nil {
// 		return fmt.Errorf("error dropping index: %w", err)
// 	}
// 	if err := p.DB.Exec(`CREATE INDEX male_f_index ON records ("sex", "name") WHERE "sex" = TRUE AND "name" ILIKE 'F%';`).Error; err != nil {
// 		return fmt.Errorf("error creating index: %w", err)
// 	}
// 	return nil
// }

func (p *postgres) GetUniqueRecords() error {
	var uniqueRecords []models.Records
	if err := p.DB.Select("DISTINCT ON (name, dob) name, dob, sex").Order("name ASC").Find(&uniqueRecords).Error; err != nil {
		return err
	}
	for _, record := range uniqueRecords {
		age := time.Since(record.User.DOB).Hours() / 24 / 365
		date := record.User.DOB.Format("2006-01-02")
		if record.User.Sex {
			fmt.Println(*record.User.Name, date, "Male", int(age))
		} else {
			fmt.Println(*record.User.Name, date, "Female", int(age))
		}
	}
	return nil
}

func (p *postgres) AutoFill() error {

	firstNames, err := readNamesFromFile("internal/repository/postgres/first-names.txt")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	midNames, err := readNamesFromFile("internal/repository/postgres/middle-names.txt")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	lastNames, err := readNamesFromFile("internal/repository/postgres/last-names.txt")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	fNames, err := readNamesFromFile("internal/repository/postgres/male-F-names.txt")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	// Отсюда не исполняется, ошибку не выдает
	var alphabet [26]string = [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := 0; i < 1000000; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		midName := midNames[rand.Intn(len(midNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]

		lastName = alphabet[rand.Intn(len(alphabet))] + strings.ToLower(lastName)

		// fullName := fmt.Sprintf("%s %s %s", lastName, midName, firstName)

		fullName := lastName + " " + midName + " " + firstName

		gender := rand.Intn(2) == 1

		user := &models.User{
			Name: &fullName,
			DOB:  randomDate(),
			Sex:  gender,
		}
		record := models.Records{User: *user}

		if err := p.DB.Create(&record).Error; err != nil {
			return fmt.Errorf("create: %w", err)
		}
	}

	for i := 0; i < 100; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		midName := midNames[rand.Intn(len(midNames))]
		fName := lastNames[rand.Intn(len(fNames))]
		fullName := fName + " " + midName + " " + firstName
		// fullName := fmt.Sprintf("%s %s %s", fName, midName, firstName)

		user := &models.User{
			Name: &fullName,
			DOB:  randomDate(),
			Sex:  true,
		}
		record := models.Records{User: *user}

		if err := p.DB.Create(&record).Error; err != nil {
			return err
		}
	}
	if err := p.CreateIndex(); err != nil {
		return fmt.Errorf("indexing error: %w", err)
	}
	return nil
}

func (p *postgres) GetRecords() error {
	var fRecords []models.Records
	start := time.Now()

	rows, err := p.DB.Raw("SELECT * FROM records WHERE name LIKE 'F%' AND sex = ?", true).Rows()

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var record models.Records

		err := rows.Scan(&record.ID, &record.User.Name, &record.User.DOB, &record.User.Sex)

		if err != nil {
			return err
		}

		fRecords = append(fRecords, record)
	}
	for _, record := range fRecords {
		fmt.Println(*record.Name, record.DOB.Format("2006/01/02"), "Male")
	}
	fmt.Println(time.Since(start))
	return nil
}

func (p *postgres) GetIndexedRecords() error {
	var records []models.Records
	start := time.Now()

	// The index should be used in the following query
	err := p.DB.Where("sex = ? AND name ILIKE ?", true, "F%").Find(&records).Error
	if err != nil {
		return err
	}

	for _, record := range records {
		fmt.Println(*record.Name, record.DOB.Format("2006/01/02"), "Male")
	}
	fmt.Println(time.Since(start))
	return nil
}

func readNamesFromFile(filename string) ([]string, error) {

	var names []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

func randomDate() time.Time {
	year := rand.Intn(100) + 1920

	month := rand.Intn(12) + 1

	day := rand.Intn(30) + 1

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
