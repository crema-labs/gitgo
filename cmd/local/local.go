package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/crema-labs/gitgo/pkg/model"
	"github.com/crema-labs/gitgo/pkg/server"
	"github.com/crema-labs/gitgo/pkg/store"
)

func main() {
	db, err := store.NewSQLiteStore("./grants.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
		grant := &model.Grant{
			GrantID:     randomString(),
			GrantAmount: "2000",
			Status:      "open",
			Contributions: map[string]float64{
				"0x01489AF266B02AF5C727E0C0AF76d36Ea63CE1971Cd24A5A234d6E485b8c9d65": 100,
			},
		}

		err := db.InsertGrant(grant)
		if err != nil {
			log.Fatalf("Failed to insert grant: %v", err)
		}

	}
	x, err := db.GetAllGrants()
	if err != nil {
		log.Fatalf("Failed to get all grants: %v", err)
	}

	fmt.Println(x)

	s := server.NewServer(db)
	log.Fatal(s.Start(":8080"))
}

func randomString() string {
	//32 byte hex string
	data := [32]byte{}

	_, err := rand.Read(data[:])
	if err != nil {
		log.Fatalf("Failed to read random data: %v", err)
	}

	return hex.EncodeToString(data[:])

}
