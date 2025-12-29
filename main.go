// filebrowser-quantum-sync
// Author: Steve Zabka
// Author-URL: https://cryinkfly.com
// License:  Apache-2.0
//
// Version: 1.0.0

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	bolt "go.etcd.io/bbolt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	dbPath := "/db/database.db"
	outPath := "/config/users"

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatalf("DB open failed: %v", err)
	}
	defer db.Close()

	out, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("users create failed: %v", err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	seen := map[string]bool{}
	count := 0

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(_ []byte, b *bolt.Bucket) error {
			return b.ForEach(func(_, v []byte) error {
				var u User
				if json.Unmarshal(v, &u) != nil {
					return nil
				}
				if u.Username == "" || u.Password == "" {
					return nil
				}
				if !strings.HasPrefix(u.Password, "$2") {
					return nil
				}
				if seen[u.Username] {
					return nil
				}

				fmt.Fprintf(w, "%s:%s\n", u.Username, u.Password)
				seen[u.Username] = true
				count++
				return nil
			})
		})
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("users file written (%d users)\n", count)
}
