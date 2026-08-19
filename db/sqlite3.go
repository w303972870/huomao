package huomao_db

import (
	//"fmt"
	"database/sql"
	"database/sql/driver"
	_ "github.com/mattn/go-sqlite3"
	"github.com/w303972870/huomao/tool"
	"path/filepath"
)

type sqlite3 struct {
	db *sql.DB
}

func (s *sqlite3) Connect(dsn string) {
	//dbdir = fmt.Sprint( dbdir , "/db/" )
	//driver_tool.HM_Tools.MkDir( dbdir )
	//s.db, _ = sql.Open( "sqlite3" , fmt.Sprint( dbdir , "sso.db" ) )
	dbdir, _ := filepath.Split(dsn)
	huomao_tool.GetToolClass().MkDir(dbdir)
	s.db, _ = sql.Open("sqlite3", dsn)
}

func (s *sqlite3) Close() {
	s.db.Close()
}

func (s *sqlite3) ExecSql(sql string) error {
	_, err := s.db.Exec(sql)
	return err
}

func (s *sqlite3) exec(sql string) (driver.Result, error) {
	return s.db.Exec(sql)
}

func (s *sqlite3) Insert(sql string) (int, error) {

	res, err := s.exec(sql)
	if err != nil {
		return 0, err
	}
	a, err := res.LastInsertId()
	return int(a), err
}

func (s *sqlite3) Delete(sql string) (int, error) {
	res, err := s.exec(sql)
	if err != nil {
		return 0, err
	}
	a, err := res.RowsAffected()
	return int(a), err
}

func (s *sqlite3) Update(sql string) (int, error) {
	res, err := s.exec(sql)
	if err != nil {
		return 0, err
	}
	a, err := res.RowsAffected()
	return int(a), err
}

func (s *sqlite3) Select(sql string) (error, *sql.Rows) {
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
