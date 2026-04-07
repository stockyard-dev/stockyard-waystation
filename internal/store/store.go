package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type DB struct{ db *sql.DB }
type Trip struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Destination string `json:"destination"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Budget      int    `json:"budget"`
	Itinerary   string `json:"itinerary"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
	CreatedAt   string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "waystation.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS trips(id TEXT PRIMARY KEY,name TEXT NOT NULL,destination TEXT DEFAULT '',start_date TEXT DEFAULT '',end_date TEXT DEFAULT '',budget INTEGER DEFAULT 0,itinerary TEXT DEFAULT '[]',status TEXT DEFAULT 'planning',notes TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
	resource TEXT NOT NULL,
	record_id TEXT NOT NULL,
	data TEXT NOT NULL DEFAULT '{}',
	PRIMARY KEY(resource, record_id)
)`)
	return &DB{db: db}, nil
}
func (d *DB) Close() error { return d.db.Close() }
func genID() string        { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string          { return time.Now().UTC().Format(time.RFC3339) }
func (d *DB) Create(e *Trip) error {
	e.ID = genID()
	e.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO trips(id,name,destination,start_date,end_date,budget,itinerary,status,notes,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`, e.ID, e.Name, e.Destination, e.StartDate, e.EndDate, e.Budget, e.Itinerary, e.Status, e.Notes, e.CreatedAt)
	return err
}
func (d *DB) Get(id string) *Trip {
	var e Trip
	if d.db.QueryRow(`SELECT id,name,destination,start_date,end_date,budget,itinerary,status,notes,created_at FROM trips WHERE id=?`, id).Scan(&e.ID, &e.Name, &e.Destination, &e.StartDate, &e.EndDate, &e.Budget, &e.Itinerary, &e.Status, &e.Notes, &e.CreatedAt) != nil {
		return nil
	}
	return &e
}
func (d *DB) List() []Trip {
	rows, _ := d.db.Query(`SELECT id,name,destination,start_date,end_date,budget,itinerary,status,notes,created_at FROM trips ORDER BY created_at DESC`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Trip
	for rows.Next() {
		var e Trip
		rows.Scan(&e.ID, &e.Name, &e.Destination, &e.StartDate, &e.EndDate, &e.Budget, &e.Itinerary, &e.Status, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}
func (d *DB) Update(e *Trip) error {
	_, err := d.db.Exec(`UPDATE trips SET name=?,destination=?,start_date=?,end_date=?,budget=?,itinerary=?,status=?,notes=? WHERE id=?`, e.Name, e.Destination, e.StartDate, e.EndDate, e.Budget, e.Itinerary, e.Status, e.Notes, e.ID)
	return err
}
func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM trips WHERE id=?`, id)
	return err
}
func (d *DB) Count() int { var n int; d.db.QueryRow(`SELECT COUNT(*) FROM trips`).Scan(&n); return n }

func (d *DB) Search(q string, filters map[string]string) []Trip {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (name LIKE ?)"
		args = append(args, "%"+q+"%")
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(`SELECT id,name,destination,start_date,end_date,budget,itinerary,status,notes,created_at FROM trips WHERE `+where+` ORDER BY created_at DESC`, args...)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Trip
	for rows.Next() {
		var e Trip
		rows.Scan(&e.ID, &e.Name, &e.Destination, &e.StartDate, &e.EndDate, &e.Budget, &e.Itinerary, &e.Status, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Stats() map[string]any {
	m := map[string]any{"total": d.Count()}
	rows, _ := d.db.Query(`SELECT status,COUNT(*) FROM trips GROUP BY status`)
	if rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}
	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
