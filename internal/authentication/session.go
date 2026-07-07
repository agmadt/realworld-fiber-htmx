package authentication

import (
	"time"

	"realworld-fiber-htmx/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type gormStorage struct {
	db         *gorm.DB
	gcInterval time.Duration
	done       chan struct{}
}

type sessionEntry struct {
	Key    string `gorm:"column:k;primaryKey;size:64"`
	Value  []byte `gorm:"column:v"`
	Expiry int64  `gorm:"column:e"`
}

func (sessionEntry) TableName() string {
	return "fiber_storage"
}

func (s *gormStorage) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}

	var entry sessionEntry
	result := s.db.Where("k = ?", key).First(&entry)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	if entry.Expiry != 0 && entry.Expiry <= time.Now().Unix() {
		return nil, nil
	}

	return entry.Value, nil
}

func (s *gormStorage) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}

	var expiry int64
	if exp != 0 {
		expiry = time.Now().Add(exp).Unix()
	}

	return s.db.Table("fiber_storage").Exec(
		"INSERT OR REPLACE INTO fiber_storage (k, v, e) VALUES (?, ?, ?)",
		key, val, expiry,
	).Error
}

func (s *gormStorage) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}
	return s.db.Where("k = ?", key).Delete(&sessionEntry{}).Error
}

func (s *gormStorage) Reset() error {
	return s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&sessionEntry{}).Error
}

func (s *gormStorage) Close() error {
	s.done <- struct{}{}
	return nil
}

func (s *gormStorage) gcTicker() {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case t := <-ticker.C:
			s.gc(t)
		}
	}
}

func (s *gormStorage) gc(t time.Time) {
	s.db.Exec("DELETE FROM fiber_storage WHERE e <= ? AND e != 0", t.Unix())
}

var StoredAuthenticationSession *session.Store

func SessionStart() {

	db := database.Get()
	db.AutoMigrate(&sessionEntry{})
	db.Exec("CREATE INDEX IF NOT EXISTS idx_e ON fiber_storage (e)")

	store := &gormStorage{
		db:         db,
		gcInterval: 10 * time.Second,
		done:       make(chan struct{}),
	}
	go store.gcTicker()

	authSession := session.New(session.Config{
		Storage: store,
	})

	StoredAuthenticationSession = authSession
}

func AuthStore(c *fiber.Ctx, userID uint) {
	session, err := StoredAuthenticationSession.Get(c)
	if err != nil {
		panic(err)
	}

	session.Set("authentication", userID)
	if err := session.Save(); err != nil {
		panic(err)
	}
}

func AuthGet(c *fiber.Ctx) (bool, uint) {
	session, err := StoredAuthenticationSession.Get(c)
	if err != nil {
		panic(err)
	}

	value := session.Get("authentication")
	if value == nil {
		return false, 0
	}

	return true, value.(uint)
}

func AuthDestroy(c *fiber.Ctx) {
	session, err := StoredAuthenticationSession.Get(c)
	if err != nil {
		panic(err)
	}

	session.Delete("authentication")
	if err := session.Save(); err != nil {
		panic(err)
	}
}
