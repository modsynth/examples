package service

import (
	"testing"

	"github.com/modsynth/e-commerce-api/internal/config"
	"github.com/modsynth/e-commerce-api/internal/domain"
	"github.com/modsynth/e-commerce-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

func setupTestConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret-key",
			AccessTTL:  900000000000,  // 15 minutes in nanoseconds
			RefreshTTL: 604800000000000, // 7 days in nanoseconds
		},
	}
}

func TestAuthService_Register(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()
	userRepo := repository.NewUserRepository(db)
	authService := NewAuthService(userRepo, cfg)

	tests := []struct {
		name    string
		req     *domain.RegisterRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid registration",
			req: &domain.RegisterRequest{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req: &domain.RegisterRequest{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "Duplicate",
				LastName:  "User",
			},
			wantErr: true,
			errMsg:  "email already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := authService.Register(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("Register() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			// Validate user was created correctly
			if user == nil {
				t.Error("Register() returned nil user")
				return
			}

			if user.Email != tt.req.Email {
				t.Errorf("Register() email = %v, want %v", user.Email, tt.req.Email)
			}

			if user.Role != domain.RoleCustomer {
				t.Errorf("Register() role = %v, want %v", user.Role, domain.RoleCustomer)
			}

			if !user.IsActive {
				t.Error("Register() user is not active")
			}

			// Verify password was hashed
			if user.PasswordHash == tt.req.Password {
				t.Error("Register() did not hash password")
			}

			// Verify password hash is valid
			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(tt.req.Password))
			if err != nil {
				t.Error("Register() password hash is invalid")
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()
	userRepo := repository.NewUserRepository(db)
	authService := NewAuthService(userRepo, cfg)

	// Create test user
	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	testUser := &domain.User{
		Email:        "login@example.com",
		PasswordHash: string(hashedPassword),
		FirstName:    "Login",
		LastName:     "Test",
		Role:         domain.RoleCustomer,
		IsActive:     true,
	}

	if err := userRepo.Create(testUser); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		req     *domain.LoginRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid credentials",
			req: &domain.LoginRequest{
				Email:    "login@example.com",
				Password: password,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			req: &domain.LoginRequest{
				Email:    "notfound@example.com",
				Password: password,
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
		{
			name: "invalid password",
			req: &domain.LoginRequest{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := authService.Login(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("Login() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			// Validate response
			if resp == nil {
				t.Error("Login() returned nil response")
				return
			}

			if resp.AccessToken == "" {
				t.Error("Login() access token is empty")
			}

			if resp.RefreshToken == "" {
				t.Error("Login() refresh token is empty")
			}

			if resp.User == nil {
				t.Error("Login() user is nil")
				return
			}

			if resp.User.Email != tt.req.Email {
				t.Errorf("Login() user email = %v, want %v", resp.User.Email, tt.req.Email)
			}
		})
	}
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()
	userRepo := repository.NewUserRepository(db)
	authService := NewAuthService(userRepo, cfg)

	// Create inactive test user
	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	inactiveUser := &domain.User{
		Email:        "inactive@example.com",
		PasswordHash: string(hashedPassword),
		FirstName:    "Inactive",
		LastName:     "User",
		Role:         domain.RoleCustomer,
		IsActive:     false, // Inactive user
	}

	if err := userRepo.Create(inactiveUser); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	req := &domain.LoginRequest{
		Email:    "inactive@example.com",
		Password: password,
	}

	_, err := authService.Login(req)
	if err == nil {
		t.Error("Login() should fail for inactive user")
		return
	}

	if err.Error() != "account is inactive" {
		t.Errorf("Login() error = %v, want 'account is inactive'", err.Error())
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()
	userRepo := repository.NewUserRepository(db)
	authService := NewAuthService(userRepo, cfg)

	// Create test user
	testUser := &domain.User{
		Email:        "refresh@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "Refresh",
		LastName:     "Test",
		Role:         domain.RoleCustomer,
		IsActive:     true,
	}

	if err := userRepo.Create(testUser); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Generate valid refresh token
	refreshToken, err := authService.(*authService).generateRefreshToken(testUser)
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid refresh token",
			token:   refreshToken,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := authService.RefreshToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Error("RefreshToken() returned nil response")
					return
				}

				if resp.AccessToken == "" {
					t.Error("RefreshToken() access token is empty")
				}

				if resp.RefreshToken == "" {
					t.Error("RefreshToken() refresh token is empty")
				}
			}
		})
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	db := setupTestDB(t)
	cfg := setupTestConfig()
	userRepo := repository.NewUserRepository(db)
	authService := NewAuthService(userRepo, cfg)

	// Create test user
	testUser := &domain.User{
		Email:        "getuser@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "GetUser",
		LastName:     "Test",
		Role:         domain.RoleCustomer,
		IsActive:     true,
	}

	if err := userRepo.Create(testUser); err != nil {
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
			user, err := authService.GetUserByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("GetUserByID() returned nil user")
					return
				}

				if user.ID != tt.id {
					t.Errorf("GetUserByID() user ID = %v, want %v", user.ID, tt.id)
				}
			}
		})
	}
}
