package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
)

// application holds the dependencies for our web application, such as loggers and the user repository.
type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	userRepo    UserRepository
	postRepo    PostRepository
	templateDir string
	publicPath  string
	tp          *TemplateRenderer
	session     *sessions.Session
}

func main() {
	db, err := connectToDatabase("users_database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	session := sessions.New([]byte("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4"))
	session.Lifetime = 24 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteLaxMode

	app := &application{
		errorLog:    log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		infoLog:     log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		userRepo:    NewSQLUserRepository(db),
		postRepo:    NewSQLPostRepository(db),
		templateDir: "./templates",
		publicPath:  "./public",
		session:     session,
	}
	app.tp = NewTemplateRenderer(app.templateDir, false) // 2nd parameter isDev is for running in localdev

	log.Println("Listening on :8080")
	if err := app.serve(); err != nil {
		log.Fatal(err)
	}

	// TODO: Testing code to clean up later
	// posts, meta, err := app.postRepo.GetAll(Filter{
	// 	PageSize: 10,
	// 	Page:     1,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("posts : %+v\n", posts)
	// fmt.Printf("meta : %+v\n", meta)
}

// connectToDatabase establishes a connection to the SQLite database and returns the database handle.
func connectToDatabase(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
