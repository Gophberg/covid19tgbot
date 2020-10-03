package bot

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PQ ...
)

// Nums ...
type Nums struct {
	ID                   int
	Cases                int
	Deaths               int
	Recovered            int
	Isolation            int
	Isolationathome      int
	Observation          int
	Quarantineathospital int
	Unquarantined        int
	Peopleathospital     int
	Ambulanced           int
}

// User ...
type User struct {
	ID                       int
	Tgid                     int
	Actived                  bool
	Subscases                bool
	Subsdeaths               bool
	Subsrecovered            bool
	Subsisolation            bool
	Subsisolationathome      bool
	Subsobservation          bool
	Subsquarantineathospital bool
	Subsunquarantined        bool
	Subspeopleathospital     bool
	Subsambulanced           bool
}

// Db ...
var Db *sql.DB
var nums = Nums{}

// DBInit ...
func DBInit(config *Config) {
	var err error
	Db, err = sql.Open("postgres", config.DbDataSourceName)
	if err != nil {
		log.Panic(err)
	}
}

// NumsAdd ...
func NumsAdd(n int) (err error) {
	nums := Nums{
		Cases:                n,
		Deaths:               0,
		Recovered:            0,
		Isolation:            0,
		Isolationathome:      0,
		Observation:          0,
		Quarantineathospital: 0,
		Unquarantined:        0,
		Peopleathospital:     0,
		Ambulanced:           0,
	}
	statement := "insert into nums (cases, deaths, recovered, isolation, isolationathome, observation, quarantineathospital, unquarantined, peopleathospital, ambulanced) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		nums.Cases,
		nums.Deaths,
		nums.Recovered,
		nums.Isolation,
		nums.Isolationathome,
		nums.Observation,
		nums.Quarantineathospital,
		nums.Unquarantined,
		nums.Peopleathospital,
		nums.Ambulanced,
	).Scan(
		&nums.ID,
	)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Added new row to db: %v", nums)
	return
}

// UserAdd ...
func UserAdd(tgid int) (err error) {
	u := User{
		Tgid:                     tgid,
		Actived:                  true,
		Subscases:                true,
		Subsdeaths:               false,
		Subsrecovered:            false,
		Subsisolation:            false,
		Subsisolationathome:      false,
		Subsobservation:          false,
		Subsquarantineathospital: false,
		Subsunquarantined:        false,
		Subspeopleathospital:     false,
		Subsambulanced:           false,
	}
	statement := "insert into users (tgid, actived, subscases, subsdeaths, subsrecovered, subsisolation, subsisolationathome, subsobservation, subsquarantineathospital, subsunquarantined, subspeopleathospital, subsambulanced) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		u.Tgid,
		u.Actived,
		u.Subscases,
		u.Subsdeaths,
		u.Subsrecovered,
		u.Subsisolation,
		u.Subsisolationathome,
		u.Subsobservation,
		u.Subsquarantineathospital,
		u.Subsunquarantined,
		u.Subspeopleathospital,
		u.Subsambulanced,
	).Scan(
		&u.ID,
	)
	log.Printf("%v", u.ID)
	if err != nil {
		log.Panic(err)
	}
	return
}

// Users ...
func Users() (users []User, err error) {
	rows, err := Db.Query("SELECT tgid FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&user.Tgid,
		)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

// GetLast ...
func GetLast() (nums Nums, err error) {
	// nums = Nums{}
	err = Db.QueryRow("SELECT id, cases, deaths, recovered, isolation, isolationathome, observation, quarantineathospital, unquarantined, peopleathospital, ambulanced FROM nums ORDER BY id DESC LIMIT 1").Scan(
		&nums.ID,
		&nums.Cases,
		&nums.Deaths,
		&nums.Recovered,
		&nums.Isolation,
		&nums.Isolationathome,
		&nums.Observation,
		&nums.Quarantineathospital,
		&nums.Unquarantined,
		&nums.Peopleathospital,
		&nums.Ambulanced,
	)
	return nums, err
}

// IsExist ...
func IsExist(tgID int) bool {
	u := User{}
	err := Db.QueryRow("SELECT tgid FROM users where tgid = $1", tgID).Scan(&u.Tgid)
	if err != nil {
		return false
	}
	return true
}
