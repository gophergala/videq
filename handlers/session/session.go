package session

import (
	"database/sql"
	//	"log"
	"math/rand"
	"net/http"
	"time"

	alog "github.com/cenkalti/log"
)

import _ "github.com/go-sql-driver/mysql"

func Sid(r *http.Request) (string, error) {
	sid, err := r.Cookie("sid")
	if err != nil {
		return "", err
	}

	return sid.Value, nil
}

type Handler struct {
	dsn        string
	passHandle http.Handler
	log        alog.Logger
	db         *sql.DB
}

func NewHandler(log alog.Logger, db *sql.DB, passHandler http.Handler) *Handler {
	h := new(Handler)
	h.passHandle = passHandler
	h.log = log
	h.db = db

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("sid")

	if err == http.ErrNoCookie {
		// db, err := sql.Open("mysql", h.dsn)
		// if err != nil {
		// 	h.log.Error(err)
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }
		// defer db.Close()

		var sid string

		for {
			getSid := func() (bool, error) {
				sid = UniqueKey()

				_, err := h.db.Exec("LOCK TABLE session WRITE")
				if err != nil {
					return false, err
				}
				defer func() {
					_, err := h.db.Exec("UNLOCK TABLES")
					if err != nil {
						h.log.Error(err)
					}
				}()

				var exists int
				err = h.db.QueryRow("SELECT 1 FROM session WHERE sid=? ", sid).Scan(&exists)
				switch {

				case err == sql.ErrNoRows:
					_, err := h.db.Exec("INSERT INTO session (sid, ts) VALUES (?, UNIX_TIMESTAMP())", sid)
					if err != nil {
						return false, err
					}
					return true, nil

				case err != nil:
					return false, err

				default:
					return false, nil

				}
			}

			generated, err := getSid()
			if err != nil {
				h.log.Error(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if generated {
				break
			}
		}

		cookie := &http.Cookie{}
		cookie.Name = "sid"
		cookie.Value = sid
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(7 * 24 * time.Hour)

		http.SetCookie(w, cookie)
	}

	h.passHandle.ServeHTTP(w, r)
}

func UniqueKey() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	newKey := make([]rune, 32)
	for i := range newKey {
		newKey[i] = letters[rand.Intn(len(letters))]
	}
	return string(newKey)
}
