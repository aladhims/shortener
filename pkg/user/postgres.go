package user

import (
	"database/sql"

	pb "github.com/aladhims/shortener/pkg/user/proto"
)

type postgres struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgres{
		db: db,
	}
}

func (p *postgres) Create(user *pb.User) (uint64, error) {
	var id int
	err := p.db.QueryRow("INSERT INTO users(fullname, email) VALUES($1, $2) RETURNING id", user.Fullname, user.Email).Scan(&id)
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (p *postgres) Get(id uint64) *pb.User {
	stmt := `SELECT * FROM users WHERE id=$1`

	var user *pb.User

	err := p.db.QueryRow(stmt, id).Scan(user)
	if err != nil {
		return nil
	}

	return user
}

func (p *postgres) GetByEmail(email string) *pb.User {
	stmt := `SELECT * FROM users WHERE email=$1`

	var user *pb.User

	err := p.db.QueryRow(stmt, email).Scan(user)
	if err != nil {
		return nil
	}

	return user
}

func (p *postgres) Ping() error {
	return p.db.Ping()
}
