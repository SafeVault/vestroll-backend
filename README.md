# VestRoll Backend

Profile onboarding flow implemented.

How to run locally
- Go 1.25+
- No Redis required: the server auto-starts an embedded Redis for local dev.

Steps:
1) Start the server
   - go run ./cmd/server
2) Health check
   - GET http://localhost:8080/health

Profile setup endpoints
- POST /api/v1/profile/account-type
  - Body: {"user_id":"<id>", "account_type":"freelancer|contractor"}
- POST /api/v1/profile/personal-details
  - Body: {"user_id":"<id>", "data":{"first_name":"...","last_name":"...","gender":"male|female|other","date_of_birth":"YYYY-MM-DD","dial_code":"+234","phone":"8012345678"}}
- POST /api/v1/profile/address
  - Body: {"user_id":"<id>", "data":{"country":"NG","street":"...","city":"...","postal_code":"100001"}}
- GET /api/v1/profile/status?user_id=<id>

Validation rules
- Account type: freelancer|contractor
- Personal: first_name, last_name required; gender optional male|female|other; date_of_birth must be YYYY-MM-DD and >=16y; dial_code like +NNN; phone 7–20 digits
- Address: country, street, city required; postal_code optional (3–12 alnum/space/-)

Completion tracking
- 33% after account type
- 66% after personal details
- 100% after address (completed=true)

Testing
- Unit tests cover service validation and completion logic with miniredis:
  - go test ./...
