package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database = connectDB()
var DocsC, DownloadsC int64
var FileChannel = make(chan *FileDB, 11)

type FileDB struct {
	gorm.Model
	Url        string `gorm:"index;unique;not null"`
	Type       string
	downloaded bool
}

func connectDB() *gorm.DB {
	dsn := "root:12345@tcp(localhost:3306)/your_database_name?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: silentLogger})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&FileDB{})
	return db
}

func GetFileFromDb() {
	DocumentsCount()
	limit, offset := 10, 0
	go func() {
		for {
			var docs []FileDB
			result := Database.Where("url = ?", false).Offset(offset).Limit(limit).Find(&docs)
			offset += limit
			for _, i := range docs {
				FileChannel <- &i
			}
			if result.Error != nil || len(docs) == 0 {
				close(FileChannel)
				break
			}
		}
	}()
}

func (f *FileDB) DoneAndSave() {
	f.downloaded = true
	Database.Save(f)
}

var silentLogger = logger.New(
	nil, // Use your preferred io.Writer if you want to log somewhere
	logger.Config{
		LogLevel: logger.Silent,
	},
)

func DocumentsCount() {
	Database.Model(&FileDB{}).Count(&DocsC)
	Database.Model(&FileDB{}).Where("downloaded = ?", true).Count(&DownloadsC)
}
