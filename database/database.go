// Package database Database controller provide database information's
package database

import (
	"api-gateway/ent"
	"context"
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Database struct {
	Def         *ent.Client
	RedisClient *redis.Client
	Ctx         context.Context
}

var Db Database

func setUpDatabase() {
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_NAME"), os.Getenv("POSTGRES_PASSWORD")))
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(`
		with
			_service (name) as (
			VALUES
				('google-drive'),
				('github'),
				('twitter'),
				('miro'),
				('slack'),
				('gmail'),
				('google-forms'),
				('google-sheet'),
				('google-calendar'),
				('time'),
				('dropbox'),
				('spotify'),
				('google-slide'),
				('gitlab'),
				('typeform'),
				('sql')
			),
			_service_to_add as (
				select v.name
				from _service v
				where not exists (
					select 1
					from services s
					where s.name = v.name
				)
			)
			insert into services (name)
			select name
			from _service_to_add
			on conflict do nothing;
		`)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)
	db.Close()
}

// NewDatabase Initialize the databases connections (postgres and redis) with the environment variables
func NewDatabase() (Database, error) {
	client, err := ent.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_NAME"), os.Getenv("POSTGRES_PASSWORD")))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	} else {
		fmt.Println("Postgres database is ready")
	}
	RedisDb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	fmt.Println("Redis database is ready")
	fmt.Println(os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_NAME"), os.Getenv("POSTGRES_PASSWORD"))
	Db = Database{Def: client, RedisClient: RedisDb, Ctx: context.Background()}
	if err := Db.Def.Schema.Create(Db.Ctx); err != nil {
		fmt.Println(err)
		log.Fatalf("failed creating schema resources: %v", err)
	}
	setUpDatabase()
	return Db, err
}
