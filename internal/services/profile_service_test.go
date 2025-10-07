package services

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	redis "github.com/go-redis/redis/v8"
)

func setupProfileService(t *testing.T) (*ProfileService, func()) {
	mini, err := miniredis.Run()
	if err != nil { t.Fatalf("failed to start miniredis: %v", err) }
	rdb := redis.NewClient(&redis.Options{Addr: mini.Addr()})
	repo := repository.NewProfileRepository(rdb, 0)
	service := NewProfileService(repo)
	cleanup := func() { mini.Close() }
	return service, cleanup
}

func TestAccountTypeAndCompletion(t *testing.T) {
	svc, cleanup := setupProfileService(t)
	defer cleanup()

	ctx := context.Background()
	prof, err := svc.UpdateAccountType(ctx, models.AccountTypeRequest{UserID: "u1", AccountType: "freelancer"})
	if err != nil { t.Fatalf("UpdateAccountType error: %v", err) }
	if prof.CompletionPercent != 33 { t.Fatalf("expected 33, got %d", prof.CompletionPercent) }
}

func TestPersonalDOBValidation(t *testing.T) {
	svc, cleanup := setupProfileService(t)
	defer cleanup()
	ctx := context.Background()
	req := models.PersonalDetailsRequest{UserID: "u1", Data: models.PersonalDetails{
		FirstName: "A", LastName: "B", DateOfBirth: "21-07-1995", DialCode: "+1", Phone: "1234567",
	}}
	_, err := svc.UpdatePersonalDetails(ctx, req)
	if err == nil { t.Fatalf("expected DOB validation error, got nil") }
}

func TestAddressPostalValidation(t *testing.T) {
	svc, cleanup := setupProfileService(t)
	defer cleanup()
	ctx := context.Background()
	req := models.AddressRequest{UserID: "u1", Data: models.Address{Country: "NG", Street: "S", City: "C", PostalCode: "@@@"}}
	_, err := svc.UpdateAddress(ctx, req)
	if err == nil { t.Fatalf("expected postal validation error, got nil") }
}

func TestFullFlowCompletionAndPersistence(t *testing.T) {
	svc, cleanup := setupProfileService(t)
	defer cleanup()
	ctx := context.Background()

	_, err := svc.UpdateAccountType(ctx, models.AccountTypeRequest{UserID: "u2", AccountType: "contractor"})
	if err != nil { t.Fatalf("account type error: %v", err) }

	_, err = svc.UpdatePersonalDetails(ctx, models.PersonalDetailsRequest{UserID: "u2", Data: models.PersonalDetails{
		FirstName: "Ada", LastName: "Lovelace", Gender: "female", DateOfBirth: "1995-07-21", DialCode: "+234", Phone: "8012345678",
	}})
	if err != nil { t.Fatalf("personal details error: %v", err) }

	prof, err := svc.UpdateAddress(ctx, models.AddressRequest{UserID: "u2", Data: models.Address{Country: "NG", Street: "12 Marina Rd", City: "Lagos", PostalCode: "100001"}})
	if err != nil { t.Fatalf("address error: %v", err) }
	if !prof.Completed || prof.CompletionPercent != 100 { t.Fatalf("expected completed=100, got %v %d", prof.Completed, prof.CompletionPercent) }

	// Read back
	got, err := svc.GetProfile(ctx, "u2")
	if err != nil { t.Fatalf("GetProfile error: %v", err) }
	if got.Address == nil || got.Personal == nil { t.Fatalf("expected persisted personal and address") }
}
