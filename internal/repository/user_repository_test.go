package repository

import (
	"context"
	"testing"

	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/testutil"
	"github.com/House-lovers7/speadwear-go/pkg/utils"
)

func TestUserRepository_Create(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &domain.User{
				Name:           "Test User",
				Email:          "test@example.com",
				PasswordDigest: "hashedpassword",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &domain.User{
				Name:           "Another User",
				Email:          "test@example.com", // Same email as above
				PasswordDigest: "hashedpassword",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.user.ID == 0 {
				t.Error("Create() did not set user ID")
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test user
	user := fixtures.CreateUser()

	tests := []struct {
		name    string
		id      uint
		want    bool
		wantErr bool
	}{
		{
			name: "existing user",
			id:   user.ID,
			want: true,
		},
		{
			name: "non-existing user",
			id:   99999,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want && got == nil {
				t.Error("FindByID() returned nil for existing user")
			}
			if !tt.want && got != nil {
				t.Error("FindByID() returned user for non-existing ID")
			}
			if got != nil && got.ID != tt.id {
				t.Errorf("FindByID() returned wrong user ID = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test user
	user := fixtures.CreateUser(func(u *domain.User) {
		u.Email = "findbyemail@example.com"
	})

	tests := []struct {
		name    string
		email   string
		want    bool
		wantErr bool
	}{
		{
			name:  "existing email",
			email: user.Email,
			want:  true,
		},
		{
			name:  "non-existing email",
			email: "nonexistent@example.com",
			want:  false,
		},
		{
			name:  "case sensitive email",
			email: "FINDBYEMAIL@example.com",
			want:  false, // Assuming case-sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByEmail(ctx, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want && got == nil {
				t.Error("FindByEmail() returned nil for existing email")
			}
			if !tt.want && got != nil {
				t.Error("FindByEmail() returned user for non-existing email")
			}
			if got != nil && got.Email != tt.email {
				t.Errorf("FindByEmail() returned wrong email = %v, want %v", got.Email, tt.email)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test user
	user := fixtures.CreateUser()
	originalName := user.Name
	originalEmail := user.Email

	// Update user
	user.Name = "Updated Name"
	user.Email = "updated@example.com"
	user.Picture = "/uploads/avatar.jpg"

	err := repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify update
	updated, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID() after update error = %v", err)
	}

	if updated.Name == originalName {
		t.Error("Update() did not update name")
	}
	if updated.Email == originalEmail {
		t.Error("Update() did not update email")
	}
	if updated.Picture != user.Picture {
		t.Error("Update() did not update picture")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test user
	user := fixtures.CreateUser()

	// Delete user
	err := repo.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deletion
	deleted, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID() after delete error = %v", err)
	}
	if deleted != nil {
		t.Error("Delete() did not delete user")
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create multiple users
	for i := 0; i < 5; i++ {
		fixtures.CreateUser()
	}

	tests := []struct {
		name    string
		limit   int
		offset  int
		want    int
		wantErr bool
	}{
		{
			name:   "get all users",
			limit:  10,
			offset: 0,
			want:   5,
		},
		{
			name:   "get with pagination",
			limit:  2,
			offset: 0,
			want:   2,
		},
		{
			name:   "get with offset",
			limit:  10,
			offset: 3,
			want:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := repo.FindAll(ctx, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(users) != tt.want {
				t.Errorf("FindAll() returned %d users, want %d", len(users), tt.want)
			}
		})
	}
}

func TestUserRepository_Count(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create multiple users
	numUsers := 3
	for i := 0; i < numUsers; i++ {
		fixtures.CreateUser()
	}

	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if count != int64(numUsers) {
		t.Errorf("Count() = %d, want %d", count, numUsers)
	}
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test user
	user := fixtures.CreateUser(func(u *domain.User) {
		u.Email = "exists@example.com"
	})

	tests := []struct {
		name    string
		email   string
		want    bool
		wantErr bool
	}{
		{
			name:  "existing email",
			email: user.Email,
			want:  true,
		},
		{
			name:  "non-existing email",
			email: "notexists@example.com",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.ExistsByEmail(ctx, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExistsByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if exists != tt.want {
				t.Errorf("ExistsByEmail() = %v, want %v", exists, tt.want)
			}
		})
	}
}

func TestUserRepository_PasswordAuthentication(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create user with hashed password
	plainPassword := "testpassword123"
	hashedPassword, err := utils.HashPassword(plainPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := &domain.User{
		Name:           "Auth Test User",
		Email:          "auth@example.com",
		PasswordDigest: hashedPassword,
		Activated:      true,
	}

	err = repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test authentication
	found, err := repo.FindByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("FindByEmail() error = %v", err)
	}

	// Test correct password
	if !utils.CheckPassword(plainPassword, found.PasswordDigest) {
		t.Error("Password authentication failed with correct password")
	}

	// Test incorrect password
	if utils.CheckPassword("wrongpassword", found.PasswordDigest) {
		t.Error("Password authentication succeeded with incorrect password")
	}
}