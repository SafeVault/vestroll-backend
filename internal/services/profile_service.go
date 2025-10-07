package services

import (
    "context"
    "regexp"
    "strings"
    "time"

    "github.com/codeZe-us/vestroll-backend/internal/models"
    "github.com/codeZe-us/vestroll-backend/internal/repository"
)

// ProfileService orchestrates validation and persistence for the onboarding profile

type ProfileService struct {
    repo *repository.ProfileRepository
}

// GetProfile returns the current profile for a user (or an empty initialized one)
func (s *ProfileService) GetProfile(ctx context.Context, userID string) (models.UserProfile, error) {
    existing, err := s.repo.Get(ctx, userID)
    if err != nil { return models.UserProfile{}, err }
    if existing != nil { return *existing, nil }
    return models.UserProfile{UserID: userID}, nil
}

func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
    return &ProfileService{repo: repo}
}

// UpdateAccountType sets the account type and recomputes completion
func (s *ProfileService) UpdateAccountType(ctx context.Context, req models.AccountTypeRequest) (models.UserProfile, error) {
    prof := s.getOrInit(ctx, req.UserID)
    prof.AccountType = req.AccountType
    s.computeCompletion(&prof)
    prof.UpdatedAt = time.Now()
    if err := s.repo.Save(ctx, prof); err != nil { return models.UserProfile{}, err }
    return prof, nil
}

// UpdatePersonalDetails validates and stores personal details
func (s *ProfileService) UpdatePersonalDetails(ctx context.Context, req models.PersonalDetailsRequest) (models.UserProfile, error) {
    if err := validateDOB(req.Data.DateOfBirth); err != nil { return models.UserProfile{}, err }
    if err := validateDialCode(req.Data.DialCode); err != nil { return models.UserProfile{}, err }
    if err := validatePhone(req.Data.Phone); err != nil { return models.UserProfile{}, err }
    if req.Data.Gender != "" && !oneOf(strings.ToLower(req.Data.Gender), []string{"male","female","other"}) {
        return models.UserProfile{}, ErrValidation("gender must be male, female, or other")
    }

    prof := s.getOrInit(ctx, req.UserID)
    // Normalize gender lower-case
    data := req.Data
    data.Gender = strings.ToLower(data.Gender)
    prof.Personal = &data
    s.computeCompletion(&prof)
    prof.UpdatedAt = time.Now()
    if err := s.repo.Save(ctx, prof); err != nil { return models.UserProfile{}, err }
    return prof, nil
}

// UpdateAddress validates and stores address
func (s *ProfileService) UpdateAddress(ctx context.Context, req models.AddressRequest) (models.UserProfile, error) {
    if strings.TrimSpace(req.Data.Country) == "" { return models.UserProfile{}, ErrValidation("country is required") }
    if strings.TrimSpace(req.Data.Street) == "" { return models.UserProfile{}, ErrValidation("street is required") }
    if strings.TrimSpace(req.Data.City) == "" { return models.UserProfile{}, ErrValidation("city is required") }
    if req.Data.PostalCode != "" {
        // Allow simple alphanumeric 3-12
        if ok, _ := regexp.MatchString(`^[A-Za-z0-9\- ]{3,12}$`, req.Data.PostalCode); !ok {
            return models.UserProfile{}, ErrValidation("postal_code format is invalid")
        }
    }

    prof := s.getOrInit(ctx, req.UserID)
    data := req.Data
    prof.Address = &data
    s.computeCompletion(&prof)
    prof.UpdatedAt = time.Now()
    if err := s.repo.Save(ctx, prof); err != nil { return models.UserProfile{}, err }
    return prof, nil
}

func (s *ProfileService) getOrInit(ctx context.Context, userID string) models.UserProfile {
    existing, _ := s.repo.Get(ctx, userID)
    if existing != nil { return *existing }
    return models.UserProfile{UserID: userID}
}

func (s *ProfileService) computeCompletion(p *models.UserProfile) {
    steps := 0
    total := 3 // account type, personal, address
    if p.AccountType != "" { steps++ }
    if p.Personal != nil && p.Personal.FirstName != "" && p.Personal.LastName != "" && p.Personal.DateOfBirth != "" && p.Personal.Phone != "" && p.Personal.DialCode != "" { steps++ }
    if p.Address != nil && p.Address.Country != "" && p.Address.Street != "" && p.Address.City != "" { steps++ }
    p.CompletionPercent = int(float64(steps) / float64(total) * 100)
    p.Completed = p.CompletionPercent == 100
}

func validateDOB(d string) error {
    // Expect YYYY-MM-DD and be in the past
    re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
    if !re.MatchString(d) {
        return ErrValidation("date_of_birth must be YYYY-MM-DD")
    }
    t, err := time.Parse("2006-01-02", d)
    if err != nil { return ErrValidation("date_of_birth parse error") }
    if !t.Before(time.Now()) { return ErrValidation("date_of_birth must be in the past") }
    // Basic age check: at least 16 years old
    if t.After(time.Now().AddDate(-16, 0, 0)) {
        return ErrValidation("user must be at least 16 years old")
    }
    return nil
}

func validateDialCode(dc string) error {
    if ok, _ := regexp.MatchString(`^\+\d{1,4}$`, dc); !ok {
        return ErrValidation("dial_code must look like +234")
    }
    return nil
}

func validatePhone(p string) error {
    if ok, _ := regexp.MatchString(`^[0-9]{7,20}$`, p); !ok {
        return ErrValidation("phone must be 7-20 digits")
    }
    return nil
}

func oneOf(val string, allowed []string) bool {
    for _, a := range allowed { if val == a { return true } }
    return false
}
