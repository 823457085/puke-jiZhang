package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		openid TEXT UNIQUE NOT NULL,
		nickname TEXT DEFAULT '',
		avatar_url TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_code TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		game_type TEXT DEFAULT 'generic',
		creator_id INTEGER NOT NULL,
		status TEXT DEFAULT 'active',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		closed_at DATETIME,
		FOREIGN KEY (creator_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS room_members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (room_id) REFERENCES rooms(id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		UNIQUE(room_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS bills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_id INTEGER NOT NULL,
		payer_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		note TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (room_id) REFERENCES rooms(id),
		FOREIGN KEY (payer_id) REFERENCES users(id),
		FOREIGN KEY (receiver_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS settlements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_id INTEGER NOT NULL,
		settled_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (room_id) REFERENCES rooms(id)
	);

	CREATE INDEX IF NOT EXISTS idx_rooms_code ON rooms(room_code);
	CREATE INDEX IF NOT EXISTS idx_bills_room ON bills(room_id);
	CREATE INDEX IF NOT EXISTS idx_room_members_room ON room_members(room_id);
	`

	_, err := DB.Exec(schema)
	return err
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
