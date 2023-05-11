package brain

import (
	"CLITool/internal/models"
	"CLITool/internal/repository"
	"fmt"
	"reflect"
	"time"
)

type brain struct {
	repo repository.Repository
}

func New(repo repository.Repository) *brain {
	return &brain{
		repo: repo,
	}
}

func (b *brain) Start(val int, args []string) error {
	switch val {
	case 1:
		return b.CreateTable()
	case 2:
		name := args[0] + " " + args[1] + " " + args[2]

		var sex bool
		if args[4] == "Male" {
			sex = true
		} else if args[4] == "Female" {
			sex = false
		} else {
			return fmt.Errorf("Gender's first letter should be capitalized")
		}

		layout := "01/02/2006"
		date, err := time.Parse(layout, args[3])
		if err != nil {
			fmt.Println(date)
			fmt.Println(reflect.TypeOf(date))
			return fmt.Errorf("Wrong date of birth")
		}

		user := &models.User{
			Name: &name,
			DOB:  date,
			Sex:  sex,
		}
		return b.repo.CreateRecord(*user)
	case 3:
		// Вывод всех строк с уникальным значением ФИО+дата,
		// отсортированным по ФИО , вывести ФИО, Дату рождения, пол, кол-во полных лет.
		return b.repo.GetUniqueRecords()
	case 4:
		return b.repo.AutoFill()
	case 5:
		return b.repo.GetRecords()
	case 6:
		return b.repo.GetIndexedRecords()
	default:
		return fmt.Errorf("error")
	}
}
