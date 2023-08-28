package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log/syslog"
	"os/user"
	"syscall"

	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

var (
	connString string
	sLog       *syslog.Writer
)

func getGUID() string {
	guid := uuid.New()
	sLog.Info(fmt.Sprintf("getGUID(): %s", guid.String()))
	return guid.String()

}
func readConfig() (*Config, error) {
	data, err := os.ReadFile("dupfile.conf")
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

func connectToDatabase(connString string) (*sql.DB, error) {

	db, err := sql.Open("mssql", connString)
	if err != nil {
		sLog.Err(err.Error())
		return nil, err
	}
	return db, nil
}

func insertFileInfo(db *sql.DB, filename, location, extension, hashValue, permission, owner, guid string, size int64, created, modified time.Time, isDuplicate int) error {
	insertSQL := `
	INSERT INTO FileInfo (Filename, Location, Size, CreatedDate, ModifiedDate, FileExtension, HashValue, Permissions, FileOwner, GUID, IsDuplicate) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(insertSQL, filename, location, size, created, modified, extension, hashValue, permission, owner, guid, isDuplicate)
	if err != nil {
		sLog.Err(err.Error())
	}
	sLog.Info(fmt.Sprintf("insertFileInfo: {%s} %s/%s [%d] %s", guid, location, filename, size, hashValue))

	return err
}

func getFileOwner(path string) (string, error) {
	// Get file owner and permissions
	fileInfo, err := os.Stat(path)
	if err != nil {
		sLog.Err(err.Error())
		return "", err
	}
	//fileMode := fileInfo.Mode().String()

	uid := fmt.Sprintf("%d", fileInfo.Sys().(*syscall.Stat_t).Uid)
	userInfo, err := user.LookupId(uid)
	if err != nil {
		sLog.Err(err.Error())
		return "", err
	}

	return userInfo.Username, nil
}

func processFiles(dir string, db *sql.DB) error {
	hashes := make(map[string][]string)
	sessionGUID := getGUID()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Error(fmt.Sprintf("filepath.Walk Error: %s", err.Error()))
			//return err
		} else {

			if !info.IsDir() {
				owner, _ := getFileOwner(path)
				data, err := os.ReadFile(path)
				if err != nil {
					sLog.Err(fmt.Sprintf("ERROR os.ReadFile(%s): %s", path, err.Error()))

				} else {

					hash := sha256.Sum256(data)
					hashStr := hex.EncodeToString(hash[:])

					hashes[hashStr] = append(hashes[hashStr], path)

					filename := info.Name()
					permissions := fmt.Sprintf("%o", info.Mode().Perm())
					location := path
					size := info.Size()
					created := info.ModTime()
					modified := info.ModTime()
					extension := strings.TrimPrefix(filepath.Ext(filename), ".")
					isDuplicate := 0
					if len(hashes[hashStr]) > 1 {
						isDuplicate = 1
						sLog.Info(fmt.Sprintf("Duplicate file found (%s): %s", path, hashStr))
					}
					//if len(hashes[hashStr]) > 1 {
					// This means the current file is a duplicate, so insert it into the database
					if err := insertFileInfo(db, filename, location, extension, hashStr, permissions, owner, sessionGUID, size, created, modified, isDuplicate); err != nil {
						sLog.Err(fmt.Sprintf("insertFileInfo Error: %s", err.Error()))
						//return err
					}
				}
			}
		}
		return nil
	})
	sLog.Info(fmt.Sprintf("processFiles Error: %s", err.Error()))
	return nil
}

func init() {

	cfg, err := readConfig()
	if err != nil {
		log.Printf("Error reading configuration: %s", err.Error())
		panic(err.Error())
	}
	sysLog, err := syslog.Dial(cfg.Syslog.Protocol, fmt.Sprintf("%s:%d", cfg.Syslog.RemoteHost, cfg.Syslog.Port),
		syslog.LOG_INFO|syslog.LOG_DAEMON|syslog.LOG_DEBUG|syslog.LOG_ERR, cfg.Syslog.Tag)
	if err != nil {
		log.Fatal("Fatal syslog not started", err)
	}
	sLog = sysLog

	sLog.Info(fmt.Sprintf("%s:%d", cfg.Syslog.RemoteHost, cfg.Syslog.Port))
	connString = fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName)

	fmt.Fprintf(sLog, "dupfiles started")
	sLog.Emerg("And this is a daemon emergency with demotag.")
	sLog.Info("Info dupfiles started")
}

func main() {

	startTime := time.Now()

	ptrPath := flag.String("path", ".", "Path to iterate")
	flag.Parse()

	db, err := connectToDatabase(connString)
	if err != nil {
		sLog.Info(fmt.Sprintf("Error connecting to database: %s", err.Error()))
		return
	}
	defer db.Close()

	if err := processFiles(*ptrPath, db); err != nil {
		sLog.Info(fmt.Sprintf("Error processing files: %s", err.Error()))
		return
	}

	endTime := time.Now()

	sLog.Info(fmt.Sprintf("Files processed successfully! started at %s and completed at %s. Took %s", startTime, endTime, endTime.Sub(startTime)))
}
