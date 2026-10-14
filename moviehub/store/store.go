package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate review")
	ErrNoMovie   = errors.New("movie not found")
)

type Store struct {
	Pool *pgxpool.Pool
}

func (s *Store) ListMembers(ctx context.Context) ([]Member, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT id, email, display_name, role FROM members ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Member{}
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.Email, &m.DisplayName, &m.Role); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// DeleteReview ลบได้เมื่อเป็นเจ้าของรีวิว หรือเป็นแอดมิน
func (s *Store) DeleteReview(ctx context.Context, reviewID, memberID int64, isAdmin bool) error {
	tag, err := s.Pool.Exec(ctx,
		`DELETE FROM reviews
         WHERE id = $1 AND (member_id = $2 OR $3)`, reviewID, memberID, isAdmin)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
