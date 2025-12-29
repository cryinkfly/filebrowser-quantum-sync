// filebrowser-quantum-sync
// Author: Steve Zabka
// Author-URL: https://cryinkfly.com
// License:  Apache-2.0
//
// Version: 1.0.1

package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "strings"
    "time"

    bolt "go.etcd.io/bbolt"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"` // bcrypt hash
}

// Interval zwischen den Syncs
const syncInterval = 1 * time.Minute

func main() {
    for {
        err := runSync()
        if err != nil {
            log.Printf("Sync error: %v", err)
        }
        log.Printf("Next sync in %v...", syncInterval)
        time.Sleep(syncInterval)
    }
}

func runSync() error {
    dbPath := "/db/database.db"
    htpasswdPath := "/sync/users"

    db, err := bolt.Open(dbPath, 0600, nil)
    if err != nil {
        return fmt.Errorf("opening DB: %w", err)
    }
    defer db.Close()

    f, err := os.Create(htpasswdPath)
    if err != nil {
        return fmt.Errorf("creating htpasswd file: %w", err)
    }
    defer f.Close()

    writer := bufio.NewWriter(f)
    defer writer.Flush()

    userCount := 0
    exportedUsers := make(map[string]bool)

    err = db.View(func(tx *bolt.Tx) error {
        return tx.ForEach(func(_ []byte, b *bolt.Bucket) error {
            return b.ForEach(func(_, v []byte) error {
                var user User
                if err := json.Unmarshal(v, &user); err != nil {
                    return nil
                }

                username := strings.TrimSpace(user.Username)
                hash := strings.TrimSpace(user.Password)

                if username == "" || hash == "" {
                    return nil
                }
                if !strings.HasPrefix(hash, "$2") || len(hash) < 50 {
                    fmt.Printf("WARN: Invalid hash for user %s: %q\n", username, hash)
                    return nil
                }

                if exportedUsers[username] {
                    fmt.Printf("WARN: User %s already exported, skipping\n", username)
                    return nil
                }

                line := fmt.Sprintf("%s:%s\n", username, hash)
                if _, err := writer.WriteString(line); err != nil {
                    return err
                }

                exportedUsers[username] = true
                userCount++
                return nil
            })
        })
    })
    if err != nil {
        return fmt.Errorf("reading DB: %w", err)
    }

    writer.Flush()
    fmt.Printf("htpasswd successfully created (%d users)\n", userCount)
    return nil
}
