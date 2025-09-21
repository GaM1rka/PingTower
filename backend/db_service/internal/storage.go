package internal

import (
	"database/sql"
	"db_service/configs"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/lib/pq"
)

type Storage struct {
	psql *sql.DB
	ch   *sql.DB
}

type PingLog struct {
	SiteID   int       `json:"id"`
	ReqTime  time.Time `json:"req_time"`
	RespTime int64     `json:"resp_time"`
	Status   string    `json:"status"`
	Site     string    `json:"site"`
}

func NewStorage(psqlDB, chDB *sql.DB) *Storage {
	return &Storage{
		psql: psqlDB,
		ch:   chDB,
	}
}
func InitPostgreSQL() (*sql.DB, error) {
	connStr := "host=postgres_db port=5432 user=postgres password=postgres dbname=ping_db sslmode=disable"
	configs.DBLogger.Println("🔄 Attempting to connect to PostgreSQL...")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		configs.DBLogger.Printf("Failed to open DB connection: %v", err)
		return nil, fmt.Errorf("failed to open DB connection: %v", err)
	}

	// Добавляем ретраи с логированием
	var connected bool
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			connected = true
			configs.DBLogger.Println("✅ Successfully connected to PostgreSQL")
			break
		}
		configs.DBLogger.Printf("⏳ Waiting for PostgreSQL... attempt %d, error: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if !connected {
		configs.DBLogger.Printf("Could not connect to PostgreSQL after retries: %v", err)
		return nil, fmt.Errorf("could not connect to PostgreSQL after retries: %v", err)
	}

	err = createPostgreSQLTables(db)
	if err != nil {
		configs.DBLogger.Printf("Failed to create tables: %v", err)
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

func InitClickHouse() (*sql.DB, error) {
	configs.DBLogger.Println("🔄 Attempting to connect to ClickHouse...")

	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"clickhouse_db:9000"},
		Auth: clickhouse.Auth{Database: "default", Username: "default", Password: ""},
		Settings: clickhouse.Settings{
			"max_memory_usage":          0, // безлимит для запроса (или поставь, например, 1e9)
			"max_memory_usage_for_user": 0,
		},
		DialTimeout: 10 * time.Second,
	})

	// Добавляем ретраи для ClickHouse
	var connected bool
	for i := 0; i < 10; i++ {
		err := conn.Ping()
		if err == nil {
			connected = true
			configs.DBLogger.Println("✅ Successfully connected to ClickHouse")
			break
		}
		configs.DBLogger.Printf("⏳ Waiting for ClickHouse... attempt %d, error: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if !connected {
		configs.DBLogger.Println("Could not connect to ClickHouse after retries")
		return nil, fmt.Errorf("could not connect to ClickHouse after retries")
	}

	err := createClickHouseTable(conn)
	if err != nil {
		configs.DBLogger.Printf("Failed to create ClickHouse table: %v", err)
		return nil, fmt.Errorf("failed to create ClickHouse table: %v", err)
	}

	return conn, nil
}

func createPostgreSQLTables(db *sql.DB) error {
	// Таблица пользователей
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Таблица сайтов пользователей
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_sites (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			site VARCHAR(255) NOT NULL,
			check_interval INTEGER DEFAULT 60,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, site)
		)
	`)
	return err
}

func createClickHouseTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ping_logs (
			user_id UInt32,
			site String,
			req_time DateTime,
			resp_time Int64,
			status String
		) ENGINE = MergeTree()
		ORDER BY (user_id, site, req_time)
		PARTITION BY toYYYYMM(req_time)
	`)
	return err
}

// Методы работы с PostgreSQL
func (s *Storage) CreateUser(email, password string) (int, error) {
	var id int
	err := s.psql.QueryRow(`
		INSERT INTO users (email, password) 
		VALUES ($1, $2) 
		RETURNING id
	`, email, password).Scan(&id)
	return id, err
}

func (s *Storage) AddUserSite(userID int, site string) error {
	_, err := s.psql.Exec(`
		INSERT INTO user_sites (user_id, site) 
		VALUES ($1, $2)
		ON CONFLICT (user_id, site) DO NOTHING
	`, userID, site)
	return err
}

func (s *Storage) AddSiteWithCheck(userID int, site string, checkInterval int) error {
	// Добавляем сайт
	err := s.AddUserSite(userID, site)
	if err != nil {
		return err
	}

	// Добавляем начальную запись в ClickHouse
	_, err = s.ch.Exec(`
		INSERT INTO ping_logs (user_id, site, req_time, resp_time, status)
		VALUES (?, ?, ?, ?, ?)
	`, userID, site, time.Now(), 0, "initial")

	return err
}

// Методы работы с ClickHouse
func (s *Storage) GetSiteLogs(userID, siteID int) ([]PingLog, error) {
	// 1) достаём URL сайта по siteID
	var site string
	err := s.psql.QueryRow(`
        SELECT site FROM user_sites WHERE id = $1 AND user_id = $2
    `, siteID, userID).Scan(&site)
	if err != nil {
		return nil, err
	}

	// 2) читаем логи из ClickHouse
	rows, err := s.ch.Query(`
        SELECT req_time, resp_time, status, site
        FROM ping_logs
        WHERE user_id = ? AND site = ?
        ORDER BY req_time DESC
    `, userID, site)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3) маппим в структуру и проставляем SiteID
	var logs []PingLog
	for rows.Next() {
		var log PingLog
		if err := rows.Scan(&log.ReqTime, &log.RespTime, &log.Status, &log.Site); err != nil {
			return nil, err
		}
		log.SiteID = siteID // ← вот здесь добавляем ID
		logs = append(logs, log)
	}
	return logs, nil
}

func (s *Storage) GetAllUserLogs(userID int) ([]PingLog, error) {
	// URL -> site_id (чтобы вернуть id сайта)
	siteIDByURL := map[string]int{}
	rs, err := s.psql.Query(`SELECT id, site FROM user_sites WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	for rs.Next() {
		var sid int
		var url string
		if err := rs.Scan(&sid, &url); err != nil {
			rs.Close()
			return nil, err
		}
		siteIDByURL[url] = sid
	}
	rs.Close()

	// берём "первую" (последнюю по времени) запись на каждый site
	rows, err := s.ch.Query(`
        SELECT req_time, resp_time, status, site
        FROM ping_logs
        WHERE user_id = ?
        ORDER BY site ASC, req_time DESC
        LIMIT 1 BY site
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []PingLog
	for rows.Next() {
		var log PingLog
		if err := rows.Scan(&log.ReqTime, &log.RespTime, &log.Status, &log.Site); err != nil {
			return nil, err
		}
		log.SiteID = siteIDByURL[log.Site] // 0, если сайт удалили из user_sites
		logs = append(logs, log)
	}
	return logs, nil
}

// Уже, вероятно, есть — оставляю для полноты:
func (s *Storage) GetUserIDByEmail(email string) (int, error) {
	var id int
	err := s.psql.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&id)
	return id, err
}

func (s *Storage) GetUserEmail(userID int) (string, error) {
	var email string
	err := s.psql.QueryRow(`SELECT email FROM users WHERE id = $1`, userID).Scan(&email)
	return email, err
}

func (s *Storage) GetAllUsersSites() ([]UserSites, error) {
	rows, err := s.psql.Query(`
		SELECT u.id AS user_id,
			COALESCE(us.id, 0) AS site_id,          -- <-- добавили COALESCE
			COALESCE(us.site, '') AS site,
			COALESCE(us.check_interval, 0) AS check_interval
		FROM users u
		LEFT JOIN user_sites us ON us.user_id = u.id
		ORDER BY u.id, us.id
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// собираем по пользователям
	byUser := map[int][]SiteInfo{}
	userOrder := []int{}
	for rows.Next() {
		var uid, sid, interval int
		var url string
		if err := rows.Scan(&uid, &sid, &url, &interval); err != nil {
			return nil, err
		}
		if _, seen := byUser[uid]; !seen {
			userOrder = append(userOrder, uid)
			byUser[uid] = []SiteInfo{}
		}
		if sid != 0 && url != "" {
			byUser[uid] = append(byUser[uid], SiteInfo{
				ID: sid, URL: url, CheckInterval: interval,
			})
		}
	}
	// формируем итог
	out := make([]UserSites, 0, len(userOrder))
	for _, uid := range userOrder {
		out = append(out, UserSites{UserID: uid, Sites: byUser[uid]})
	}
	return out, nil
}

// Добавление лога пинга
func (s *Storage) AddPingLog(userID int, site string, respTime int64, status string) error {
	_, err := s.ch.Exec(`
		INSERT INTO ping_logs (user_id, site, req_time, resp_time, status)
		VALUES (?, ?, ?, ?, ?)
	`, userID, site, time.Now(), respTime, status)
	return err
}
