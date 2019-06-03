package postgres

import (
	"github.com/DusanKasan/payments/internal/app/payments/model"
	"github.com/go-pg/pg"
	errors "golang.org/x/xerrors"
)

type Payment struct {
	Type           string                 `sql:",notnull"`
	ID             string                 `sql:",notnull"`
	Version        int16                  `sql:",notnull"`
	OrganisationID string                 `sql:",notnull"`
	Attributes     map[string]interface{} `sql:",notnull"`
}

func (p *Payment) fromModel(payment *model.Payment) {
	p.ID = payment.ID
	p.Version = payment.Version
	p.OrganisationID = payment.OrganisationID
	p.Type = payment.Type
	p.Attributes = payment.Attributes
}

func (p *Payment) toModel() *model.Payment {
	return &model.Payment{
		ID:             p.ID,
		Version:        p.Version,
		OrganisationID: p.OrganisationID,
		Type:           p.Type,
		Attributes:     p.Attributes,
	}
}

func New(dsn string) (*storage, error) {
	options, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, errors.Opaque(err)
	}

	return &storage{pg.Connect(options)}, nil
}

type storage struct {
	db *pg.DB
}

func (s *storage) Create(payment model.Payment) (*model.Payment, error) {
	var p Payment
	p.fromModel(&payment)

	_, err := s.db.Model(&p).Returning("*").Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return nil, errors.Errorf("unable to create payment: %v: %w", err, model.ErrIdAlreadyExist)
		}

		return nil, errors.Opaque(err)
	}

	return p.toModel(), nil
}

func (s *storage) ReadOne(ID string) (*model.Payment, error) {
	payment := Payment{ID: ID}
	if err := s.db.Select(&payment); err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.Errorf("unable to read payment: %v: %w", err, model.ErrIdDoesNotExist)
		}
		return nil, errors.Opaque(err)
	}

	return payment.toModel(), nil
}

func (s *storage) ReadSortedByID(afterID *string, count int16) ([]model.Payment, error) {
	var pp []Payment

	expr := s.db.Model(&pp).OrderExpr("id ASC").Limit(int(count))
	if afterID != nil {
		expr = expr.Where("id > ?", afterID)
	}

	if err := expr.Select(); err != nil {
		return nil, errors.Opaque(err)
	}

	var m []model.Payment
	for _, p := range pp {
		m = append(m, *p.toModel())
	}

	return m, nil
}

func (s *storage) Update(ID string, payment model.Payment) (*model.Payment, error) {
	var p Payment
	p.fromModel(&payment)
	p.ID = ID

	if err := s.db.Update(&p); err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.Errorf("unable to update payment: %v: %w", err, model.ErrIdDoesNotExist)
		}
		return nil, errors.Opaque(err)
	}

	return p.toModel(), nil
}

func (s *storage) Delete(ID string) error {
	if err := s.db.Delete(&model.Payment{ID: ID}); err != nil {
		return errors.Opaque(err)
	}

	return nil
}
