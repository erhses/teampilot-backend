package dba

import (
	"fmt"
	"log"
	"os"
	"teampilot/structs/dts"
	"time"

	"github.com/joho/godotenv"
	"github.com/joomcode/errorx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dbConfig := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASS"),
	)
	db, err := gorm.Open(postgres.Open(dbConfig), &gorm.Config{
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(
		&dts.Base{},
		&dts.User{},
		&dts.Role{},
		&dts.Permission{},
		&dts.Group{},
		&dts.Feedback{},
		&dts.QuizResponse{},
		&dts.Team{},
		&dts.StudentProfile{},
		&dts.Class{},

	// &databases.Passport{},
	// &databases.User{},
	// &databases.UserBankInfo{},
	// &databases.UserCertificateInfo{},
	// &databases.UserCitizenshipInfo{},

	// // Role permission
	// &databases.Role{},
	// &databases.Permission{},
	// &databases.Group{},

	// // translation and order
	// &databases.Order{},
	// &databases.Translation{},
	// &databases.TranslationRequest{},
	// &databases.Message{},
	// &databases.FligthMetadata{},
	// &databases.TrainMetadata{},
	// &databases.HotelMetadata{},
	// &databases.RoomData{},

	// // Agent infos
	// &databases.Certificate{},
	// &databases.CertificateLevel{},

	// // Payment
	// &databases.Payment{},
	// &databases.PaymentItem{},
	// &databases.Currency{},

	// // public
	// &databases.Content{},
	// &databases.Language{},
	// &databases.Banner{},
	// &databases.Bank{},

	// // notification
	// &databases.Notification{},
	// &databases.NotificationMap{},
	)

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(10 * time.Second)
	sqlDB.SetMaxIdleConns(50)
	DB = db

	go func(dbConfig string) {
		var intervals = []time.Duration{3 * time.Second, 3 * time.Second, 15 * time.Second, 30 * time.Second, 60 * time.Second}
		for {
			time.Sleep(60 * time.Second)
			sqlDB, _ := DB.DB()
			if e := sqlDB.Ping(); e != nil {
			L:
				for i := 0; i < len(intervals); i++ {
					e2 := RetryHandler(3, func() (bool, error) {
						var err error
						DB, err = gorm.Open(postgres.Open(dbConfig), &gorm.Config{
							PrepareStmt:                              true,
							SkipDefaultTransaction:                   true,
							DisableForeignKeyConstraintWhenMigrating: true,
						})
						if err != nil {
							return false, errorx.ConcurrentUpdate.Wrap(err, "database error")
						}
						return true, nil
					})
					if e2 != nil {
						fmt.Println(e.Error())
						time.Sleep(intervals[i])
						if i == len(intervals)-1 {
							i--
						}
						continue
					}
					break L
				}

			}
		}
	}(dbConfig)
}

func RetryHandler(n int, f func() (bool, error)) error {
	ok, er := f()
	if ok && er == nil {
		return nil
	}
	if n-1 > 0 {
		return RetryHandler(n-1, f)
	}
	return er
}
