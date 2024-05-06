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
	//sqlCreateTable = "CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY, date DATETIME DEFAULT CURRENT_DATE, title TEXT NOT NULL, comment TEXT, repeat VARCHAR(45));"
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
	sqlInsertTask = "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4);"
	//sqlUpdateTask = "UPDATE scheduler SET date=$2, title=$3, comment=$4, repeat=$5 WHERE id=$1;"
	sqlUpdateTask = "UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat WHERE id=:id;"
)

func (r Repository) Save(task *task_entity.Task) (uint, error) {

	if task.TaskId == 0 {
		res, err := r.db.Exec(sqlInsertTask,
			task.Date,
			task.Title,
			task.Comment,
			task.Repeat)
		if err != nil {
			return 0, fmt.Errorf("не смог вставить, %w", err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("не смог получить id, %w", err)
		}

		return uint(id), nil
	} else {
		r.log.Printf("task for update = %v", task)
		_, err := r.db.NamedExec(sqlUpdateTask, task)
		// _, err := r.db.Exec(sqlUpdateTask,
		// 	task.TaskId,
		// 	task.Date,
		// 	task.Title,
		// 	task.Comment,
		// 	task.Repeat)
		r.log.Printf("update err=%w", err)
		if err != nil {
			return 0, fmt.Errorf("не смог обновить, %w", err)
		}
		return task.TaskId, nil
	}
}

const (
	sqlSelectAll  = "SELECT * FROM scheduler ORDER BY date ASC;"
	sqlSelectById = "SELECT * FROM scheduler WHERE id=$1;"
)

func (r Repository) Get(ids []uint) (map[uint]task_entity.Task, error) {

	r.log.Debugf("поиск идентификаторов %v", ids)

	if len(ids) > 0 {

		var count int

		err := r.db.QueryRow("SELECT COUNT(*) FROM scheduler WHERE id=$1;", ids[0]).Scan(&count)
		if err != nil {
			r.log.Errorf("%w", err)
			return nil, err
		}

		if count == 0 {
			r.log.Printf("идентификатор %d не найден", ids[0])
			return nil, fmt.Errorf("идентификатор %d не найден", ids[0])
		}

		rows, err := r.db.Queryx(sqlSelectById, ids[0])
		if err != nil {
			//кричать в лог
			r.log.Errorf("Queryx не сработал, %w", err)
			return nil, err
		}
		m := make(map[uint]task_entity.Task)
		for rows.Next() {
			var tt task_entity.Task
			err = rows.StructScan(&tt)
			if err != nil {
				//кричать в лог
				r.log.Errorf("StructScan не сработал, %w", err)
				return nil, err
			}
			m[tt.TaskId] = tt
		}
		return m, nil

	} else {
		tasks := []task_entity.Task{}
		err := r.db.Select(&tasks, sqlSelectAll)
		if err != nil {
			//кричать в лог
			r.log.Errorf("Select не сработал, %w", err)
			return nil, err
		}

		m := make(map[uint]task_entity.Task, len(tasks))
		for _, ret := range tasks {
			m[ret.TaskId] = ret
		}
		return m, nil
	}
}
