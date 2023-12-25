package sqlstore

import (
	"database/sql"

	"github.com/coltoneshaw/mattermost-plugin-customers/server/app"
	sq "github.com/mattermost/squirrel"
	"github.com/pkg/errors"
)

// customerStore holds the information needed to fulfill the methods in the store interface.
type customerStore struct {
	pluginAPI      PluginAPIClient
	store          *SQLStore
	queryBuilder   sq.StatementBuilderType
	customerSelect sq.SelectBuilder
}

type sqlCustomers struct {
	app.Customer
}

// NewCustomerStore creates a new store for customers ServiceImpl.
func NewCustomerStore(pluginAPI PluginAPIClient, sqlStore *SQLStore) app.CustomerStore {
	customerSelect := sqlStore.builder.
		Select("ci.ID", "ci.Name", "ci.Type").
		From("CRM_Customers as ci")

	return &customerStore{
		pluginAPI:      pluginAPI,
		store:          sqlStore,
		queryBuilder:   sqlStore.builder,
		customerSelect: customerSelect,
	}
}

func (s *customerStore) Get(id string) (app.Customer, error) {
	if id == "" {
		return app.Customer{}, errors.New("ID cannot be empty")
	}

	tx, err := s.store.db.Beginx()
	if err != nil {
		return app.Customer{}, errors.Wrap(err, "could not begin transaction")
	}
	defer s.store.finalizeTransaction(tx)
	var rawCustomers sqlCustomers
	err = s.store.getBuilder(tx, &rawCustomers, s.customerSelect.Where(sq.Eq{"ci.ID": id}))
	if err == sql.ErrNoRows {
		return app.Customer{}, errors.Wrapf(app.ErrNotFound, "customer does not exist for id '%s'", id)
	} else if err != nil {
		return app.Customer{}, errors.Wrapf(err, "failed to get customer by id '%s'", id)
	}

	customer, err := toCustomers(rawCustomers)
	if err != nil {
		return app.Customer{}, err
	}

	if err = tx.Commit(); err != nil {
		return app.Customer{}, errors.Wrap(err, "could not commit transaction")
	}

	return customer, nil
}
func toCustomers(rawCustomers sqlCustomers) (app.Customer, error) {
	customers := rawCustomers.Customer
	return customers, nil
}
