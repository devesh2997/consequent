package sqlclients

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
)

// Stmt is an aggregate prepared statement.
// It holds a prepared statement for each underlying physical db.
type Stmt interface {
	Close() error
	Exec(...interface{}) (sql.Result, error)
	Query(...interface{}) (*sql.Rows, error)
	QueryRow(...interface{}) *sql.Row
}

type stmt struct {
	db    *SQLXClusterDB
	stmts []*sqlx.Stmt
}

// Close closes the statement by concurrently closing all underlying
// statements concurrently, returning the first non nil error.
func (s *stmt) Close() error {
	return scatter(len(s.stmts), func(i int) error {
		return s.stmts[i].Close()
	})
}

// Exec executes a prepared statement with the given arguments
// and returns a Result summarizing the effect of the statement.
// Exec uses the master as the underlying physical db.
func (s *stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.stmts[0].Exec(args...)
}

// Query executes a prepared query statement with the given
// arguments and returns the query results as a *sql.Rows.
// Query uses a slave as the underlying physical db.
func (s *stmt) Query(args ...interface{}) (*sql.Rows, error) {
	return s.stmts[s.db.slave(len(s.db.dbs))].Query(args...)
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error
// will be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *sql.Row's Scan scans the first selected row and discards the rest.
// QueryRow uses a slave as the underlying physical db.
func (s *stmt) QueryRow(args ...interface{}) *sql.Row {
	return s.stmts[s.db.slave(len(s.db.dbs))].QueryRow(args...)
}

/*
	SQLXClusterDB is a wrapper around sqlx.SQLXClusterDB (which in itself wraps over sql.SQLXClusterDB) with multiple underlying databases.
	Reads and writes are automatically directed to the correct db (master or slave)
*/
type SQLXClusterDB struct {
	dbs          []*sqlx.DB // all available databases. first will be master, rest slaves
	currentSlave uint64     // index of next slave db to be used
}

/*
	getSQLXClientDatabase concurrently opens each underlying pre-existing sql.DB first
	one being used as the master and the rest as slaves.
*/
func getSQLXClientDatabase(driverName string, sqlDBs []*sql.DB) (*SQLXClusterDB, error) {
	db := &SQLXClusterDB{dbs: make([]*sqlx.DB, len(sqlDBs))}

	err := scatter(len(db.dbs), func(i int) (err error) {
		sqlxDB := sqlx.NewDb(sqlDBs[i], driverName)
		// By default, StructScan (converting query result rows to structs)
		// will return an error if a column does not map to a field in the destination.
		// Unsafe method a new copy with this safety turned off
		sqlxDB = sqlxDB.Unsafe()
		if err := sqlxDB.Ping(); err != nil {
			return err
		}

		db.dbs[i] = sqlxDB
		return err
	})

	if err != nil {
		return nil, err

	}

	return db, nil
}

// Close closes all physical databases concurrently, releasing any open resources.
func (db *SQLXClusterDB) Close() error {
	return scatter(len(db.dbs), func(i int) error {
		return db.dbs[i].Close()
	})
}

// Driver returns the physical database's underlying driver.
func (db *SQLXClusterDB) Driver() driver.Driver {
	return db.Master().Driver()
}

// Begin starts a transaction on the master. The isolation level is dependent on the driver.
func (db *SQLXClusterDB) Begin() (*sql.Tx, error) {
	return db.Master().Begin()
}

// Beginx starts a transaction on the master (using *sqlx.Tx). The isolation level is dependent on the driver.
func (db *SQLXClusterDB) Beginx() (*sqlx.Tx, error) {
	return db.Master().Beginx()
}

/*
	BeginTx starts a transaction with the provided context on the master.
	The provided TxOptions is optional and may be nil if defaults should be used.
	If a non-default isolation level is used that the driver doesn't support,
	an error will be returned.
*/
func (db *SQLXClusterDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.Master().BeginTx(ctx, opts)
}

/*
	BeginTxx starts a transaction with the provided context on the master (using *sqlx.Tx).
	The provided TxOptions is optional and may be nil if defaults should be used.
	If a non-default isolation level is used that the driver doesn't support,
	an error will be returned.
*/
func (db *SQLXClusterDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return db.Master().BeginTxx(ctx, opts)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
// Exec uses the master as the underlying physical db.
func (db *SQLXClusterDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Master().Exec(query, args...)
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
// Exec uses the master as the underlying physical db.
func (db *SQLXClusterDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.Master().ExecContext(ctx, query, args...)
}

// Ping verifies if a connection to each physical database is still alive,
// establishing a connection if necessary.
func (db *SQLXClusterDB) Ping() error {
	return scatter(len(db.dbs), func(i int) error {
		return db.dbs[i].Ping()
	})
}

// PingContext verifies if a connection to each physical database is still
// alive, establishing a connection if necessary.
func (db *SQLXClusterDB) PingContext(ctx context.Context) error {
	return scatter(len(db.dbs), func(i int) error {
		return db.dbs[i].PingContext(ctx)
	})
}

// Preparex creates a prepared statement for later queries or executions
// on each physical database, concurrently.
func (db *SQLXClusterDB) Preparex(query string) (Stmt, error) {
	stmts := make([]*sqlx.Stmt, len(db.dbs))

	err := scatter(len(db.dbs), func(i int) (err error) {
		stmts[i], err = db.dbs[i].Preparex(query)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &stmt{db: db, stmts: stmts}, nil
}

// PrepareContext creates a prepared statement for later queries or executions
// on each physical database, concurrently.
//
// The provided context is used for the preparation of the statement, not for
// the execution of the statement.
func (db *SQLXClusterDB) PreparexContext(ctx context.Context, query string) (Stmt, error) {
	stmts := make([]*sqlx.Stmt, len(db.dbs))

	err := scatter(len(db.dbs), func(i int) (err error) {
		stmts[i], err = db.dbs[i].PreparexContext(ctx, query)
		return err
	})

	if err != nil {
		return nil, err
	}
	return &stmt{db: db, stmts: stmts}, nil
}

/*
	Query executes a query that returns rows, typically a SELECT.
	The args are for any placeholder parameters in the query.
	Query uses a slave as the physical db.
*/
func (db *SQLXClusterDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Slave().Query(query, args...)
}

/*
	Queryx executes a query that returns rows (*sqlx.Rows), typically a SELECT.
	The args are for any placeholder parameters in the query.
	Query uses a slave as the physical db.
*/
func (db *SQLXClusterDB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.Slave().Queryx(query, args...)
}

/*
	NamedQuery executes a query that returns rows (*sqlx.Rows), typically a SELECT.
	The arg is for any named placeholder parameters in the query.
	Query uses a slave as the physical db.
*/
func (db *SQLXClusterDB) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return db.Slave().NamedQuery(query, arg)
}

// NamedExec using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *SQLXClusterDB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return db.Master().NamedExec(query, arg)
}

/*
	QueryContext executes a query that returns rows, typically a SELECT.
	The args are for any placeholder parameters in the query.
	QueryContext uses a slave as the physical db.
*/
func (db *SQLXClusterDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.Slave().QueryContext(ctx, query, args...)
}

/*
	QueryxContext executes a query that returns rows (*sqlx.Rows)(*sqlx.Rows), typically a SELECT.
	The args are for any placeholder parameters in the query.
	QueryContext uses a slave as the physical db.
*/
func (db *SQLXClusterDB) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.Slave().QueryxContext(ctx, query, args...)
}

/*
	QueryRow executes a query that is expected to return at most one row.
	QueryRow always return a non-nil value.
	Errors are deferred until Row's Scan method is called.
	QueryRow uses a slave as the physical db.
*/
func (db *SQLXClusterDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.Slave().QueryRow(query, args...)
}

/*
	QueryRowx executes a query that is expected to return at most one row i.e (*sqlx.Row)
	QueryRowx always return a non-nil value.
	Errors are deferred until Row's Scan method is called.
	QueryRowx uses a slave as the physical db.
*/
func (db *SQLXClusterDB) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return db.Slave().QueryRowx(query, args...)
}

/*
	QueryRowContext executes a query that is expected to return at most one row.
	QueryRowContext always return a non-nil value.
	Errors are deferred until Row's Scan method is called.
	QueryRowContext uses a slave as the physical db.
*/
func (db *SQLXClusterDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.Slave().QueryRowContext(ctx, query, args...)
}

/*
	QueryRowxContext executes a query that is expected to return at most one row (*sqlx.Row).
	QueryRowContext always return a non-nil value.
	Errors are deferred until Row's Scan method is called.
	QueryRowContext uses a slave as the physical db.
*/
func (db *SQLXClusterDB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return db.Slave().QueryRowxContext(ctx, query, args...)
}

/*
	Select uses rows.Scan on scannable types and rows.StructScan on non-scannable types. It is roughly
	analagous to Query, where Select is useful for fetching a slice of results
	Select uses a slave as the physical db.
*/
func (db *SQLXClusterDB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.Slave().Select(dest, query, args...)
}

/*
	Get uses rows.Scan on scannable types and rows.StructScan on non-scannable types. It is roughly
	analagous to QueryRow, where Get is useful for fetching a single row
	Get uses a slave as the physical db.
*/
func (db *SQLXClusterDB) Get(dest interface{}, query string, args ...interface{}) error {
	return db.Slave().Get(dest, query, args...)
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool for each underlying physical db.
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns then the
// new MaxIdleConns will be reduced to match the MaxOpenConns limit
// If n <= 0, no idle connections are retained.
func (db *SQLXClusterDB) SetMaxIdleConns(n int) {
	for i := range db.dbs {
		db.dbs[i].SetMaxIdleConns(n)
	}
}

// SetMaxOpenConns sets the maximum number of open connections
// to each physical database.
// If MaxIdleConns is greater than 0 and the new MaxOpenConns
// is less than MaxIdleConns, then MaxIdleConns will be reduced to match
// the new MaxOpenConns limit. If n <= 0, then there is no limit on the number
// of open connections. The default is 0 (unlimited).
func (db *SQLXClusterDB) SetMaxOpenConns(n int) {
	for i := range db.dbs {
		db.dbs[i].SetMaxOpenConns(n)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// Expired connections may be closed lazily before reuse.
// If d <= 0, connections are reused forever.
func (db *SQLXClusterDB) SetConnMaxLifetime(d time.Duration) {
	for i := range db.dbs {
		db.dbs[i].SetConnMaxLifetime(d)
	}
}

// Slave returns one of the physical databases which is a slave
func (db *SQLXClusterDB) Slave() *sqlx.DB {
	return db.dbs[db.slave(len(db.dbs))]
}

// Master returns the master database
func (db *SQLXClusterDB) Master() *sqlx.DB {
	return db.dbs[0]
}

// Master returns the current slave database
// Every time this function is called, it returns the next
// slave database according to round robin
func (db *SQLXClusterDB) slave(n int) int {
	if n <= 1 {
		return 0
	}
	return int(1 + (atomic.AddUint64(&db.currentSlave, 1) % uint64(n-1)))
}
