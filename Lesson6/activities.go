package main

import "database/sql"

//Activity type
type Activity struct {
	UserID int
	Date   int64
	Name   string
}

func (a *Activity) connection() (*sql.DB, error) {
	s, err := m.ShardByID(a.UserID)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

func (a *Activity) rconnection() (*sql.DB, error) {
	s, err := rm.ShardByID(a.UserID)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

//Create func
func (a *Activity) Create() error {
	c, err := a.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO "activities" VALUES ($1, $2, $3)`, a.UserID,
		a.Date, a.Name)
	return err
}

func (a *Activity) Read() error {
	var (
		c   *sql.DB
		err error
	)

	//task3
	p.srSwitch *= -1

	if p.srSwitch > 0 {
		c, err = a.rconnection()
	} else {
		c, err = a.connection()
	}

	// c, err := a.rconnection()
	if err != nil {
		return err
	}
	r := c.QueryRow(`SELECT "date", "name" FROM "activities" WHERE 
	"user_id" = $1`, a.UserID)
	return r.Scan(
		&a.Date,
		&a.Name,
	)
}

//Update func
func (a *Activity) Update() error {
	c, err := a.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "activities" SET "date" = $2, "name" = $3 
	WHERE "user_id" = $1`, a.UserID,
		a.Date, a.Name)
	return err
}

//Delete func
func (a *Activity) Delete() error {
	c, err := a.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "activities" WHERE "user_id" = $1`, a.UserID)
	return err
}
