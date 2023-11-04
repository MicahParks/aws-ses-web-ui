package postgres

import (
	"context"
	"fmt"
	netMail "net/mail"
	"time"

	"github.com/MicahParks/httphandle/middleware/ctxkey"
	jt "github.com/MicahParks/jsontype"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MicahParks/aws-ses-web-ui/model"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Postgres {
	return Postgres{pool: pool}
}
func (p Postgres) WriteCompose(ctx context.Context, c model.Compose) (model.Compose, error) {
	tx := ctx.Value(ctxkey.Tx).(pgx.Tx)

	var err error
	if c.UUID == uuid.Nil {
		c.UUID, err = uuid.NewRandom()
		if err != nil {
			return model.Compose{}, fmt.Errorf("failed to generate new compose UUID: %w", err)
		}
	}

	rTo := recipientsSlice(c.RecipientTo)
	rCC := recipientsSlice(c.RecipientCC)
	rBCC := recipientsSlice(c.RecipientBCC)

	//language=postgresql
	const query = `
INSERT INTO compose.message (uuid, from_addr, recipient_to, recipient_cc, recipient_bcc, subject, body)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created
`
	var id int64
	var created time.Time
	err = tx.QueryRow(ctx, query, c.UUID, c.From.Get().String(), rTo, rCC, rBCC, c.Subject, c.Body).Scan(&id, &created)
	if err != nil {
		return model.Compose{}, fmt.Errorf("failed to write new compose: %w", err)
	}
	c.ID = id
	c.Created = created

	return c, nil
}
func (p Postgres) UpdateComposeSESAccepted(ctx context.Context, u uuid.UUID, sesAccepted bool) error {
	tx := ctx.Value(ctxkey.Tx).(pgx.Tx)

	//language=postgresql
	const query = `
UPDATE compose.message
SET ses_accepted = $1
WHERE uuid = $2
`
	_, err := tx.Exec(ctx, query, sesAccepted, u)
	if err != nil {
		return fmt.Errorf("failed to update compose SES accepted: %w", err)
	}

	return nil
}

func recipientsSlice(s []*jt.JSONType[*netMail.Address]) []string {
	ptrs := make([]string, 0, len(s))
	for _, v := range s {
		ptrs = append(ptrs, v.Get().String())
	}
	return ptrs
}
