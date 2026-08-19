package huomao_db

import (
	//"fmt"
	//"identity/driver/tool"
	"database/sql"
	"database/sql/driver"
	_ "github.com/go-sql-driver/mysql"
)

type mysql struct {
	db *sql.DB
}

func (s *mysql) Connect(dsn string) {
	s.db, _ = sql.Open("mysql", dsn)
}

func (s *mysql) Close() {
	s.db.Close()
}

func (s *mysql) ExecSql(sql string) error {
	_, err := s.db.Exec(sql)
	return err
}

func (s *mysql) exec(sql string) (driver.Result, error) {
	return s.db.Exec(sql)
}

func (s *mysql) Insert(sql string) (int, error) {

	res, err := s.exec(sql)
	if err != nil {
		return 0, err
	}
	a, err := res.LastInsertId()
	return int(a), err
}

func (s *mysql) Delete(sql string) (int, error) {
	res, err := s.exec(sql)
	if err != nil {
		return 0, err
	}
	a, err := res.RowsAffected()
	return int(a), err
}

func (s *mysql) Update(sql string) (int, error) {
	res, err := s.exec(sql)
	if err != nil {
		return 0, err
	}
	a, err := res.RowsAffected()
	return int(a), err
}

func (s *mysql) Select(sql string) (error, *sql.Rows) {
	stmt, err := s.db.Prepare(sql)
	if err != nil {
		return err, nil
	}
	rows, err := stmt.Query()
	if err != nil {
		return err, nil
	}
	return nil, rows
}
