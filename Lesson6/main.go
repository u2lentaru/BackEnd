package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

//Shard struct
type Shard struct {
	Address string
	Number  int
}

//Manager struct - shard manager
type Manager struct {
	size int
	ss   *sync.Map
}

//ErrorShardNotFound
var (
	ErrorShardNotFound = errors.New("shard not found")
)

//NewManager func
func NewManager(size int) *Manager {
	return &Manager{
		size: size,
		ss:   &sync.Map{},
	}
}

//Add func
func (m *Manager) Add(s *Shard) {
	m.ss.Store(s.Number, s)
}

//ShardByID func
func (m *Manager) ShardByID(entityID int) (*Shard, error) {
	if entityID < 0 {
		return nil, ErrorShardNotFound
	}
	n := entityID % m.size
	if s, ok := m.ss.Load(n); ok {
		return s.(*Shard), nil
	}
	return nil, ErrorShardNotFound
}

//Pool type
type Pool struct {
	sync.RWMutex
	cc       map[string]*sql.DB
	srSwitch int
}

//NewPool func
func NewPool() *Pool {
	return &Pool{
		cc:       map[string]*sql.DB{},
		srSwitch: 1,
	}
}

//Connection func
func (p *Pool) Connection(addr string) (*sql.DB, error) {
	p.RLock()
	if c, ok := p.cc[addr]; ok {
		defer p.RUnlock()
		return c, nil
	}
	p.RUnlock()

	p.Lock()
	defer p.Unlock()
	if c, ok := p.cc[addr]; ok {
		return c, nil
	}
	var err error
	p.cc[addr], err = sql.Open("postgres", addr)
	return p.cc[addr], err
}

//User type
type User struct {
	UserID int
	Name   string
	Age    int
	Spouse int
}

func (u *User) connection() (*sql.DB, error) {
	s, err := m.ShardByID(u.UserID)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

//Create func
func (u *User) Create() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO "users" VALUES ($1, $2, $3, $4)`, u.UserID,
		u.Name, u.Age, u.Spouse)
	return err
}

func (u *User) Read() error {
	var (
		c   *sql.DB
		err error
	)

	//task3
	p.srSwitch *= -1

	if p.srSwitch > 0 {
		c, err = u.rconnection()
	} else {
		c, err = u.connection()
	}

	if err != nil {
		return err
	}

	r := c.QueryRow(`SELECT "name", "age", "spouse" FROM "users" WHERE 
	"user_id" = $1`, u.UserID)
	return r.Scan(
		&u.Name,
		&u.Age,
		&u.Spouse,
	)
}

//Update func
func (u *User) Update() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "users" SET "name" = $2, "age" = $3, "spouse" = 
	$4 WHERE "user_id" = $1`, u.UserID,
		u.Name, u.Age, u.Spouse)
	return err
}

//Delete func
func (u *User) Delete() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "users" WHERE "user_id" = $1`, u.UserID)
	return err
}

////////////////////////task2/////////////////////

//RManager struct - replica manager
type RManager struct {
	size int
	ss   *sync.Map
}

//NewRManager func
func NewRManager(size int) *RManager {
	return &RManager{
		size: size,
		ss:   &sync.Map{},
	}
}

//Add func
func (rm *RManager) Add(s *Shard) {
	rm.ss.Store(s.Number, s)
}

//ShardByID func
func (rm *RManager) ShardByID(entityID int) (*Shard, error) {
	if entityID < 0 {
		return nil, ErrorShardNotFound
	}
	n := entityID % rm.size
	if s, ok := rm.ss.Load(n); ok {
		return s.(*Shard), nil
	}
	return nil, ErrorShardNotFound
}

func (u *User) rconnection() (*sql.DB, error) {
	s, err := rm.ShardByID(u.UserID)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

var (
	m  = NewManager(3)
	rm = NewRManager(3)
	p  = NewPool()
)

////////////////////////main////////////////////

func main() {
	m.Add(&Shard{"port=8100 user=test password=test dbname=test sslmode=disable", 0})
	m.Add(&Shard{"port=8110 user=test password=test dbname=test sslmode=disable", 1})
	m.Add(&Shard{"port=8120 user=test password=test dbname=test sslmode=disable", 2})

	rm.Add(&Shard{"port=8101 user=test password=test dbname=test sslmode=disable", 0})
	rm.Add(&Shard{"port=8111 user=test password=test dbname=test sslmode=disable", 1})
	rm.Add(&Shard{"port=8121 user=test password=test dbname=test sslmode=disable", 2})

	uu := []*User{
		{3, "Joe Biden", 78, 10},
		{10, "Jill Biden", 69, 3},
		{14, "Donald Trump", 74, 17},
		{17, "Melania Trump", 78, 14},
	}
	for _, u := range uu {
		err := u.Create()
		// err := u.Read()
		if err != nil {
			fmt.Println(fmt.Errorf("error on create user %v: %w", u, err))
		}
	}
}
