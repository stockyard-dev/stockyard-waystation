package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Trip struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Destination string `json:"destination"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Budget int `json:"budget"`
	Itinerary string `json:"itinerary"`
	Status string `json:"status"`
	Notes string `json:"notes"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"waystation.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS trips(id TEXT PRIMARY KEY,name TEXT NOT NULL,destination TEXT DEFAULT '',start_date TEXT DEFAULT '',end_date TEXT DEFAULT '',budget INTEGER DEFAULT 0,itinerary TEXT DEFAULT '[]',status TEXT DEFAULT 'planning',notes TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Trip)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO trips(id,name,destination,start_date,end_date,budget,itinerary,status,notes,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Destination,e.StartDate,e.EndDate,e.Budget,e.Itinerary,e.Status,e.Notes,e.CreatedAt);return err}
func(d *DB)Get(id string)*Trip{var e Trip;if d.db.QueryRow(`SELECT id,name,destination,start_date,end_date,budget,itinerary,status,notes,created_at FROM trips WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Destination,&e.StartDate,&e.EndDate,&e.Budget,&e.Itinerary,&e.Status,&e.Notes,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Trip{rows,_:=d.db.Query(`SELECT id,name,destination,start_date,end_date,budget,itinerary,status,notes,created_at FROM trips ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Trip;for rows.Next(){var e Trip;rows.Scan(&e.ID,&e.Name,&e.Destination,&e.StartDate,&e.EndDate,&e.Budget,&e.Itinerary,&e.Status,&e.Notes,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM trips WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM trips`).Scan(&n);return n}
