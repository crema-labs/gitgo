package server_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/crema-labs/gitgo/pkg/model"
	"github.com/crema-labs/gitgo/pkg/store"
)

var testDB *store.SQLiteStore

func TestMain(m *testing.M) {
	// Set up
	var err error
	testDB, err = store.NewSQLiteStore(":memory:")
	if err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Tear down
	testDB.Close()

	os.Exit(code)
}

func setupTestData() error {
	grant := &model.Grant{
		GrantID:     "123",
		GrantAmount: "1000",
		Status:      "open",
		Contributions: map[string]float64{
			"0x123": 500,
			"0x456": 500,
		},
	}

	contributionsJSON, err := json.Marshal(grant.Contributions)
	if err != nil {
		return err
	}

	_, err = testDB.DB().Exec(
		"INSERT INTO grants (grantid, grant_amount, status, contributions) VALUES (?, ?, ?, ?)",
		grant.GrantID, grant.GrantAmount, grant.Status, string(contributionsJSON),
	)
	return err
}

func clearTestData() error {
	_, err := testDB.DB().Exec("DELETE FROM grants")
	return err
}

func TestGetGrant(t *testing.T) {
	if err := setupTestData(); err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}
	defer clearTestData()

	tests := []struct {
		name          string
		grantID       string
		expectedGrant *model.Grant
		expectError   bool
	}{
		{
			name:    "Existing Grant",
			grantID: "123",
			expectedGrant: &model.Grant{
				GrantID:     "123",
				GrantAmount: "1000",
				Status:      "open",
				Contributions: map[string]float64{
					"0x123": 500,
					"0x456": 500,
				},
			},
			expectError: false,
		},
		{
			name:          "Non-existing Grant",
			grantID:       "456",
			expectedGrant: nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grant, err := testDB.GetGrant(tt.grantID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				if err != store.ErrGrantNotFound {
					t.Errorf("Expected ErrGrantNotFound, but got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if grant == nil {
					t.Fatalf("Expected grant, but got nil")
				}
				if grant.GrantID != *&tt.expectedGrant.GrantID {
					t.Errorf("Unexpected grant: got %v, want %v", grant, tt.expectedGrant)
				}
			}
		})
	}
}

func TestUpdateGrantStatus(t *testing.T) {
	if err := setupTestData(); err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}
	defer clearTestData()

	tests := []struct {
		name        string
		grantID     string
		newStatus   string
		expectError bool
	}{
		{
			name:        "Existing Grant",
			grantID:     "123",
			newStatus:   "closed",
			expectError: false,
		},
		{
			name:        "Non-existing Grant",
			grantID:     "456",
			newStatus:   "closed",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testDB.UpdateGrantStatus(tt.grantID, tt.newStatus)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				updatedGrant, err := testDB.GetGrant(tt.grantID)
				if err != nil {
					t.Errorf("Failed to get updated grant: %v", err)
				}
				if updatedGrant.Status != tt.newStatus {
					t.Errorf("Grant status not updated: got %v, want %v", updatedGrant.Status, tt.newStatus)
				}
			}
		})
	}
}
