package server_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/crema-labs/gitgo/pkg/model"
	"github.com/crema-labs/gitgo/pkg/store"
)

var testDB *store.SQLiteStore

func TestMain(m *testing.M) {
	var err error
	testDB, err = store.NewSQLiteStore("./sql")
	if err != nil {
		panic(err)
	}
	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func setupTestData() (*model.Grant, error) {
	grant := &model.Grant{
		GrantID:     "123",
		GrantAmount: "1000",
		Status:      "open",
		Contributions: map[string]float64{
			"0x123": 500,
			"0x456": 500,
		},
	}
	err := testDB.InsertGrant(grant)
	return grant, err
}

func clearTestData() error {
	_, err := testDB.DB().Exec("DELETE FROM grants")
	return err
}

func TestInsertGrant(t *testing.T) {
	defer clearTestData()

	grant := &model.Grant{
		GrantID:     "test123",
		GrantAmount: "2000",
		Status:      "open",
		Contributions: map[string]float64{
			"0x789": 1000,
			"0xabc": 1000,
		},
	}

	err := testDB.InsertGrant(grant)
	if err != nil {
		t.Fatalf("Failed to insert grant: %v", err)
	}

	// Verify the grant was inserted correctly
	insertedGrant, err := testDB.GetGrant(grant.GrantID)
	if err != nil {
		t.Fatalf("Failed to get inserted grant: %v", err)
	}

	if !reflect.DeepEqual(grant, insertedGrant) {
		t.Errorf("Inserted grant does not match original. Got %+v, want %+v", insertedGrant, grant)
	}
}

func TestGetGrant(t *testing.T) {
	existingGrant, err := setupTestData()
	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}
	defer clearTestData()

	tests := []struct {
		name          string
		grantID       string
		expectedGrant *model.Grant
		expectError   bool
		expectedError error
	}{
		{
			name:          "Existing Grant",
			grantID:       existingGrant.GrantID,
			expectedGrant: existingGrant,
			expectError:   false,
		},
		{
			name:          "Non-existing Grant",
			grantID:       "456",
			expectedGrant: nil,
			expectError:   true,
			expectedError: store.ErrGrantNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grant, err := testDB.GetGrant(tt.grantID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				} else if err != tt.expectedError {
					t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if !reflect.DeepEqual(grant, tt.expectedGrant) {
					t.Errorf("Grants do not match. Got %+v, want %+v", grant, tt.expectedGrant)
				}
			}
		})
	}
}

func TestUpdateGrantStatus(t *testing.T) {
	existingGrant, err := setupTestData()
	if err != nil {
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
			grantID:     existingGrant.GrantID,
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
				} else {
					updatedGrant, err := testDB.GetGrant(tt.grantID)
					if err != nil {
						t.Errorf("Failed to get updated grant: %v", err)
					} else if updatedGrant.Status != tt.newStatus {
						t.Errorf("Grant status not updated: got %v, want %v", updatedGrant.Status, tt.newStatus)
					}
				}
			}
		})
	}
}
