package db

type DatabaseInterface interface {
	Create(model any) error
	Close() error
	Migrate()
}
