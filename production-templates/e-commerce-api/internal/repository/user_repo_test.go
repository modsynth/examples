package repository

import (
	"testing"

	"github.com/modsynth/e-commerce-api/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Auto-migrate schema
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &domain.User{
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
				FirstName:    "Test",
				LastName:     "User",
				Role:         domain.RoleCustomer,
				IsActive:     true,
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &domain.User{
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
				FirstName:    "Another",
				LastName:     "User",
				Role:         domain.RoleCustomer,
				IsActive:     true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.user.ID == 0 {
				t.Error("Create() did not set user ID")
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create test user
	testUser := &domain.User{
		Email:        "find@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "Find",
		LastName:     "Test",
		Role:         domain.RoleCustomer,
		IsActive:     true,
	}

	if err := repo.Create(testUser); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing user",
			email:   "find@example.com",
			wantErr: false,
		},
		{
			name:    "non-existing user",
			email:   "notfound@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindByEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("FindByEmail() returned nil user")
					return
				}
				if user.Email != tt.email {
					t.Errorf("FindByEmail() got email = %v, want %v", user.Email, tt.email)
				}
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create test user
	testUser := &domain.User{
		Email:        "findid@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "FindID",
		LastName:     "Test",
		Role:         domain.RoleCustomer,
		IsActive:     true,
	}

	if err := repo.Create(testUser); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      testUser.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("FindByID() returned nil user")
					return
				}
				if user.ID != tt.id {
					t.Errorf("FindByID() got ID = %v, want %v", user.ID, tt.id)
				}
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create test user
	testUser := &domain.User{
		Email:        "update@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "Original",
		LastName:     "Name",
		Role:         domain.RoleCustomer,
		IsActive:     true,
	}

	if err := repo.Create(testUser); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Update user
	testUser.FirstName = "Updated"
	testUser.LastName = "Name"

	if err := repo.Update(testUser); err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}

	// Verify update
	updatedUser, err := repo.FindByID(testUser.ID)
	if err != nil {
		t.Fatalf("failed to find updated user: %v", err)
	}

	if updatedUser.FirstName != "Updated" {
		t.Errorf("Update() did not update FirstName, got = %v, want %v", updatedUser.FirstName, "Updated")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create test user
	testUser := &domain.User{
		Email:        "delete@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "Delete",
		LastName:     "Test",
		Role:         domain.RoleCustomer,
		IsActive:     true,
	}

	if err := repo.Create(testUser); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	userID := testUser.ID

	// Delete user
	if err := repo.Delete(userID); err != nil {
		t.Errorf("Delete() error = %v", err)
		return
	}

	// Verify deletion
	_, err := repo.FindByID(userID)
	if err == nil {
		t.Error("Delete() did not delete user")
	}
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create multiple test users
	for i := 1; i <= 25; i++ {
		user := &domain.User{
			Email:        string(rune(i)) + "list@example.com",
			PasswordHash: "hashed_password",
			FirstName:    "List",
			LastName:     "Test",
			Role:         domain.RoleCustomer,
			IsActive:     true,
		}
		if err := repo.Create(user); err != nil {
			t.Fatalf("failed to create test user %d: %v", i, err)
		}
	}

	tests := []struct {
		name      string
		page      int
		limit     int
		wantCount int
		wantTotal int64
	}{
		{
			name:      "first page",
			page:      1,
			limit:     10,
			wantCount: 10,
			wantTotal: 25,
		},
		{
			name:      "second page",
			page:      2,
			limit:     10,
			wantCount: 10,
			wantTotal: 25,
		},
		{
			name:      "last page",
			page:      3,
			limit:     10,
			wantCount: 5,
			wantTotal: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, err := repo.List(tt.page, tt.limit)
			if err != nil {
				t.Errorf("List() error = %v", err)
				return
			}

			if len(users) != tt.wantCount {
				t.Errorf("List() got count = %v, want %v", len(users), tt.wantCount)
			}

			if total != tt.wantTotal {
				t.Errorf("List() got total = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}
