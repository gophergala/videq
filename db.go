package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	conn *sql.DB // global variable to share it between main and short lived functions (and eg. the HTTP handler)
}

type DbConfig struct {
	DbHost string
	DbName string
	DbUser string
	DbPass string
	Debug  bool
}

//func OpenDB(host, name, user, pass string) *Database {
func (d *Database) OpenDB(cfg DbConfig) (error, *Database) {

	dba, err := sql.Open("mysql", cfg.DbUser+":"+cfg.DbPass+"@tcp("+cfg.DbHost+":3306)/"+cfg.DbName+"?charset=utf8mb4,utf8")
	if err != nil {
		// VAZNO: NIKAD SE NE OKINE!!!! (al svejedno treba provjeravat ... XXX TODO doh)
		//log.Debug(err)
		//		log.Debug("debug %s", Password("secret"))
		return err, nil
	}
	// prebacio u fju iznad koji poziva OpenDB. zasto? da mi se ne zatvori kad izadjem iz ove fje?
	//	defer db.Close()

	err = dba.Ping() // zato se cesto koristi ping
	if err != nil {
		//log.Fatal(err)
		return err, nil
	}
	dba.SetMaxIdleConns(100)
	dba.SetMaxOpenConns(200)

	return nil, &Database{conn: dba}
}

func (d *Database) CloseDB() {
	d.conn.Close()
}
