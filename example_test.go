package main

import (
	"context"
	"testing"

	"github.com/openshift-online/rh-trex-core/api"
	"github.com/openshift-online/rh-trex-core/dao"
	"github.com/openshift-online/rh-trex-core/services"
)

// Example resource using the core library
type User struct {
	api.Meta
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Mock DAO for testing
type mockUserDAO struct {
	users map[string]*User
}

func (m *mockUserDAO) Get(ctx context.Context, id string) (*User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserDAO) Create(ctx context.Context, user *User) (*User, error) {
	user.ID = "test-id"
	m.users[user.ID] = user
	return user, nil
}

func (m *mockUserDAO) Replace(ctx context.Context, user *User) (*User, error) {
	m.users[user.ID] = user
	return user, nil
}

func (m *mockUserDAO) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserDAO) List(ctx context.Context, query api.ListQuery) ([]User, error) {
	var users []User
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

func (m *mockUserDAO) Count(ctx context.Context, query api.ListQuery) (int, error) {
	return len(m.users), nil
}

func (m *mockUserDAO) FindByIDs(ctx context.Context, ids []string) ([]User, error) {
	var users []User
	for _, id := range ids {
		if user, exists := m.users[id]; exists {
			users = append(users, *user)
		}
	}
	return users, nil
}

// Mock event emitter
type mockEventEmitter struct{}

func (m *mockEventEmitter) EmitEvent(ctx context.Context, source, sourceID string, eventType api.EventType) error {
	return nil
}

// Test that demonstrates the core library functionality
func TestCoreLibraryBasicFunctionality(t *testing.T) {
	// Create mock dependencies
	mockDAO := &mockUserDAO{users: make(map[string]*User)}
	mockEvents := &mockEventEmitter{}

	// Create service using core library
	userService := services.NewBaseCRUDService[User](mockDAO, mockEvents, "Users")

	// Test Create
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	createdUser, err := userService.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if createdUser.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got %s", createdUser.Name)
	}

	// Test Get
	fetchedUser, err := userService.Get(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if fetchedUser.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %s", fetchedUser.Email)
	}

	// Test List
	list, err := userService.List(context.Background(), api.ListQuery{Page: 1, Size: 10})
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if list.Total != 1 {
		t.Errorf("Expected total 1, got %d", list.Total)
	}

	// Test event handlers
	if err := userService.OnUpsert(context.Background(), "test-id"); err != nil {
		t.Errorf("OnUpsert failed: %v", err)
	}

	if err := userService.OnDelete(context.Background(), "test-id"); err != nil {
		t.Errorf("OnDelete failed: %v", err)
	}
}