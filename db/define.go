package huomao_db

import (
	"database/sql"
)

type iDb interface {
	Connect(string)
	Close()
	ExecSql(string) error
	Insert(string) (int, error)
	Delete(string) (int, error)
	Update(string) (int, error)
	Select(string) (error, *sql.Rows)
}

func GetDbClass(which string) iDb {
	switch which {
	case "mysql":
		return &mysql{}
	case "sqlite3":
		return &sqlite3{}
	}
	return nil
}
