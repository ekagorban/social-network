package mysql

import (
	"database/sql"
	"log"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

func New() *Store {
	cfg := mysqldriver.Config{
		User:   "root",
		Passwd: "pass",
		Net:    "tcp",
		Addr:   "db:3306",
		DBName: "db",
	}
	dst := cfg.FormatDSN()
	log.Println(dst)
	db, err := sql.Open("mysql", dst)
	if err != nil {
		panic(err)
	}

	for db.Ping() != nil {
		log.Println("here0")
		time.Sleep(1 * time.Second)
	}
	log.Println("here")
	_, err = db.Exec(`
create table user_data(
    id      varchar(36)   not null
        primary key,
    name    varchar(100)  null,
    surname varchar(100)  null,
    age     smallint      null,
    gender  char          null,
    hobbies varchar(1000) null,
    city    varchar(100)  null
);
`)
	log.Println(err)

	_, err = db.Exec(`
create table user_access
(
    login    varchar(20) not null
        primary key,
    password varchar(50) not null,
    user_id  varchar(36) null,
    constraint user_access_user_data_id_fk
        foreign key (user_id) references user_data (id)
            on update cascade on delete cascade
);
`)
	log.Println(err)

	_, err = db.Exec(`
create table friends
(
    user_id   varchar(36) null,
    friend_id varchar(36) null,
    constraint friend_id_fk
        foreign key (friend_id) references user_data (id)
            on update cascade on delete cascade,
    constraint user_id_fk
        foreign key (user_id) references user_data (id)
            on update cascade on delete cascade
);
`)
	log.Println(err)
	return &Store{db}
}

// docker run  --name my-mysql -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=test -p 33070:3306 -d mysql:8.0
// docker run --network todo-app --network-alias mysql --name my-mysql -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=test -p 3306:3306 -d mysql:8.0

//-v todo-mysql-data:/var/lib/mysql
