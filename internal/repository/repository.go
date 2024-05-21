package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgshulgin/go_final_project/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db  *sqlx.DB
	log logrus.FieldLogger
}

const (
	sqlCreateTable = "CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY, date VARCHAR(8), title TEXT NOT NULL, comment TEXT, repeat VARCHAR(45));"
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
		log.Printf("файл хранилища %s не найден, создается новое хранилище", storage)
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

	log.Printf("открытие хранилища %s", storage)
	db, err := sqlx.Connect("sqlite3", storage)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть хранилище, %w", err)
	}
	return &Repository{db, log}, nil
}

// Поиск файла БД по заданному пути
// Возвращает true, если файл существует
func checkStorageExist(name string) bool {
	dbPath := name
	dir, file := filepath.Split(name)
	if len(dir) == 0 {
		//Предположение: файл хранилища находится рядом с исполняемым файлом
		appPath, err := os.Executable()
		if err != nil {
			return false
		}
		dbPath = filepath.Join(filepath.Dir(appPath), file)
	}
	_, err := os.Stat(dbPath)
	return err == nil
}

const (
	sqlInsertTask = "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);"
	sqlUpdateTask = "UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat WHERE id=:id;"
)

// Сохраняет информацию о задаче в БД
// Идентификатор задачи равный 0 означает, что поступила новая задача, которую необходимо
// создать (INSERT) БД. Идентификатор задачи со значением большим нуля означает, что
// такая задача уже существует в БД, и необходимо выполнить обновление (UPDATE).
func (r Repository) Save(task *entity.Task) (uint, error) {

	if task.TaskId == 0 {
		res, err := r.db.NamedExec(sqlInsertTask, task)
		if err != nil {
			return 0, fmt.Errorf("ошибка БД, %w", err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("ошибка БД, %w", err)
		}

		return uint(id), nil
	}

	_, err := r.db.NamedExec(sqlUpdateTask, task)
	if err != nil {
		return 0, fmt.Errorf("ошибка БД, %w", err)
	}
	return task.TaskId, nil
}

const (
	sqlSelectAll  = "SELECT * FROM scheduler ORDER BY date ASC;"
	sqlSelectById = "SELECT * FROM scheduler WHERE id=$1;"
	sqlCountById  = "SELECT COUNT(*) FROM scheduler WHERE id=$1;"
)

// Возвращает информацию о задачах по списку идентификаторов, либо все задачи
// если список идентификаторов пуст.
func (r Repository) Get(ids []uint) (map[uint]entity.Task, error) {

	r.log.Debugf("поиск идентификаторов %v", ids)

	if len(ids) > 0 {
		var count int
		err := r.db.QueryRow(sqlCountById, ids[0]).Scan(&count)
		if err != nil {
			return nil, err
		}

		if count == 0 {
			return nil, fmt.Errorf("идентификатор %d не найден", ids[0])
		}

		rows, err := r.db.Queryx(sqlSelectById, ids[0])
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		m := make(map[uint]entity.Task)
		for rows.Next() {
			var tt entity.Task
			err = rows.StructScan(&tt)
			if err != nil {
				return nil, err
			}
			m[tt.TaskId] = tt
		}
		return m, nil
	}

	tasks := []entity.Task{}
	err := r.db.Select(&tasks, sqlSelectAll)
	if err != nil {
		return nil, err
	}

	m := make(map[uint]entity.Task, len(tasks))
	for _, ret := range tasks {
		m[ret.TaskId] = ret
	}
	return m, nil
}

const (
	sqlDeleteById = "DELETE FROM scheduler WHERE id=$1;"
)

func (r Repository) Delete(ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("не указан идентификатор для удаления")
	}
	_, err := r.db.Exec(sqlDeleteById, ids[0])
	if err != nil {
		return err
	}

	return nil
}

const (
	sqlSelectByDate         = "SELECT * FROM scheduler WHERE date=$1;"
	sqlSelectByTitleComment = "SELECT * FROM scheduler WHERE title LIKE $1 OR comment LIKE $1;"
)

// Возвращает информацию о задачах соотв заданным параметрам.
// Поддерживается поиск только по полям Date, Title, Comment
func (r Repository) Lookup(task entity.Task) (map[uint]entity.Task, error) {

	r.log.Debugf("поиск задач по параметрам %v", task)

	if len(task.Date) > 0 {
		//поиск по дате
		rows, err := r.db.Queryx(sqlSelectByDate, task.Date)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		m := make(map[uint]entity.Task)
		for rows.Next() {
			var tt entity.Task
			err = rows.StructScan(&tt)
			if err != nil {
				return nil, err
			}
			m[tt.TaskId] = tt
		}
		return m, nil
	}

	// поиск по Title или Comment
	param := "%" + task.Title + "%"
	rows, err := r.db.Queryx(sqlSelectByTitleComment, param) //task.Title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[uint]entity.Task)
	for rows.Next() {
		var tt entity.Task
		err = rows.StructScan(&tt)
		if err != nil {
			return nil, err
		}
		m[tt.TaskId] = tt
	}
	return m, nil
}
