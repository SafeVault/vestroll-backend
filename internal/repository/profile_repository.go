package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/codeZe-us/vestroll-backend/internal/models"
    "github.com/go-redis/redis/v8"
)

// ProfileRepository persists user onboarding profiles in Redis
// Key pattern: user_profile:{user_id}
type ProfileRepository struct {
    client *redis.Client
    ttl    time.Duration
}

func NewProfileRepository(client *redis.Client, ttl time.Duration) *ProfileRepository {
    return &ProfileRepository{client: client, ttl: ttl}
}

func (r *ProfileRepository) key(userID string) string {
    return fmt.Sprintf("user_profile:%s", userID)
}

func (r *ProfileRepository) Save(ctx context.Context, profile models.UserProfile) error {
    data, err := json.Marshal(profile)
    if err != nil { return err }
    return r.client.Set(ctx, r.key(profile.UserID), data, r.ttl).Err()
}

func (r *ProfileRepository) Get(ctx context.Context, userID string) (*models.UserProfile, error) {
    val, err := r.client.Get(ctx, r.key(userID)).Bytes()
    if err != nil {
        if err == redis.Nil { return nil, nil }
        return nil, err
    }
    var p models.UserProfile
    if err := json.Unmarshal(val, &p); err != nil { return nil, err }
    return &p, nil
}