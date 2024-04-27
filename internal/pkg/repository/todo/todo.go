package todo

import (
	"fmt"
	"os"
	"path/filepath"

	task_entity "github.com/dgshulgin/go_final_project/internal/pkg/entity"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db  *sqlx.DB
	log logrus.FieldLogger
}

const (
	sqlCreateTable = "CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY, date DATETIME DEFAULT CURRENT_DATE, title TEXT NOT NULL, comment TEXT, repeat VARCHAR(45));"
	sqlCreateIndex = "CREATE INDEX `fk_scheduler_date` ON `scheduler` (`date` ASC);"
)

func (repo Repository) Close() {
	if repo.db != nil {
		repo.db.Close()
	}
}

func NewRepository(storage string, log logrus.FieldLogger) (*Repository, error) {
	ok := checkStorageExist(storage)
	if !ok {
		log.Printf("файл хранилища %s не найден, создаю новое хранилище...\n", storage)
		db, err := sqlx.Connect("sqlite3", storage)
		if err != nil {
			return nil, fmt.Errorf("не удалось создать хранилище, %w", err)
		}
		defer db.Close()
		_, err = db.Exec(sqlCreateTable)
		if err != nil {
			return nil, fmt.Errorf("не удалось создать таблицу, %w", err)
		}
		_, err = db.Exec(sqlCreateIndex)
		if err != nil {
			return nil, fmt.Errorf("не удалось создать табличный индекс, %w", err)
		}
	}

	log.Printf("открытие хранилища %s\n", storage)
	db, err := sqlx.Connect("sqlite3", storage)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть хранилище, %w", err)
	}
	return &Repository{db, log}, nil
}

// Варианты:
// scheduler.db
// DBFILE=scheduler.db
// DBFILE=path/to/scheduler.db
func checkStorageExist(name string) bool {
	dbPath := name
	dir, file := filepath.Split(name)
	if len(dir) == 0 {
		//Делаем предположение, что файл хранилища находится рядом с исполняемым файлом
		//Нарастить имя файла хранилища до полного пути
		appPath, err := os.Executable()
		if err != nil {
			return false
		}
		dbPath = filepath.Join(filepath.Dir(appPath), file)
		fmt.Printf("new dbPath: %s\n", dbPath)
	}
	_, err := os.Stat(dbPath)
	return err == nil //true, если файл существует
}

const (
	sqlInsertTask = "insert into scheduler (date, title, comment, repeat) values ($1, $2, $3, $4);"
)

func (r Repository) Save(task *task_entity.Task) (uint, error) {
	//tx := r.db.MustBegin()
	res, err := r.db.Exec(sqlInsertTask,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("не смог вставить, %w", err)
	}
	//tx.Commit()

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("не смог получить id, %w", err)
	}

	return uint(id), nil
}
