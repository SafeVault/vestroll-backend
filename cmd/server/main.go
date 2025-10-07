package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alicebob/miniredis/v2"
	"github.com/codeZe-us/vestroll-backend/internal/config"
	"github.com/codeZe-us/vestroll-backend/internal/database"
	"github.com/codeZe-us/vestroll-backend/internal/handlers"
	authhandlers "github.com/codeZe-us/vestroll-backend/internal/handlers/auth"
	"github.com/codeZe-us/vestroll-backend/internal/handlers/auth"
	"github.com/codeZe-us/vestroll-backend/internal/middleware"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/codeZe-us/vestroll-backend/internal/services/email_service"
	"github.com/codeZe-us/vestroll-backend/internal/services/sms_service"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
)

func main() {
	// Load configuration
	cfg := config.Load()

	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "VestRoll API is running",
			"version": "1.0.0",
		})
	})

	// Initialize PostgreSQL client
	postgresClient, err := database.NewPostgresClient(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgresClient.Close()

	// Initialize Redis client
	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
		// Start embedded in-memory Redis for local dev/testing
		if mini, merr := miniredis.Run(); merr == nil {
			log.Printf("Embedded Redis (miniredis) started at %s", mini.Addr())
			redisClient = redis.NewClient(&redis.Options{Addr: mini.Addr()})
		} else {
			log.Printf("Failed to start embedded Redis: %v", merr)
		}
	}

	// Initialize services only if Redis is available
	var otpHandler *handlers.OTPHandler
var businessProfileHandler *handlers.BusinessProfileHandler
	var profileHandler *handlers.ProfileHandler
	var pinHandler *handlers.PINHandler
	var emailVerificationHandler *authhandlers.EmailVerificationHandler
	if redisClient != nil {
		// Initialize repositories
		otpRepo := repository.NewOTPRepository(redisClient, cfg.OTP.TTL)
		businessRepo := repository.NewBusinessProfileRepository(redisClient, 0)
		profileRepo := repository.NewProfileRepository(redisClient, 0)
		pinRepo := repository.NewPinRepository(redisClient, 0)

		// Initialize services
		smsService := sms_service.NewSMSService(cfg.Twilio)
		emailService := email_service.NewEmailService(cfg.SMTP)
		smsService := sms_service.NewSMSService(cfg.Twilio)
		emailService := email_service.NewEmailService(cfg.SMTP)
		otpService := services.NewOTPService(otpRepo, smsService, emailService, cfg.OTP)
		businessService := services.NewBusinessProfileService(businessRepo)
		profileService := services.NewProfileService(profileRepo)
		pinService := services.NewPINService(pinRepo)

		// Email verification wiring
		userRepo := repository.NewUserRepository(redisClient)

		// Initialize handlers
		otpHandler = handlers.NewOTPHandler(otpService)
		businessProfileHandler = handlers.NewBusinessProfileHandler(businessService)
		profileHandler = handlers.NewProfileHandler(profileService)
		pinHandler = handlers.NewPINHandler(pinService)
		emailVerificationHandler = authhandlers.NewEmailVerificationHandler(otpService, userRepo)
	}

	// Initialize authentication services (if PostgreSQL is available)
	var loginHandler *auth.LoginHandler
	if db != nil {
		// Initialize user repository and auth services
		userRepo := repository.NewUserRepository(db)
		jwtService := services.NewJWTService(cfg.JWT)
		authService := services.NewAuthService(userRepo, jwtService)

		// Initialize login handler
		loginHandler = auth.NewLoginHandler(authService)
	}

	// API routes
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			// Apply rate limiting to OTP endpoints
			auth.Use(middleware.OTPRateLimitMiddleware())
			

			// User authentication endpoints (always available)
			userHandler.RegisterRoutes(auth)

			// Login endpoints (only if PostgreSQL is available)
			if loginHandler != nil {
				// Apply specific rate limiting for login
				loginGroup := auth.Group("/")
				loginGroup.Use(middleware.LoginRateLimitMiddleware())
				auth_routes.RegisterLoginRoutes(loginGroup, loginHandler)
			}

			
			// OTP endpoints (only if Redis is available)
			if otpHandler != nil {
				otpHandler.RegisterRoutes(auth)
			}

			// PIN endpoints (only if Redis is available)
			if pinHandler != nil {
				pinHandler.RegisterRoutes(auth)
			}

			// Google OAuth endpoints
			googleOAuth := authhandlers.NewGoogleOAuth(
				cfg.Google.ClientID,
				cfg.Google.ClientSecret,
				cfg.Google.RedirectURL,
			)
			authhandlers.RegisterGoogleAuthRoutes(auth, googleOAuth)

			// Apple OAuth endpoints
			appleOAuth := authhandlers.NewAppleOAuth(
				cfg.Apple.ClientID,
				cfg.Apple.ClientSecret,
				cfg.Apple.RedirectURL,
			)
			authhandlers.RegisterAppleAuthRoutes(auth, appleOAuth)

			// Existing auth endpoints
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - Coming soon"})
			})
			// Register endpoint (placeholder for future implementation)
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Register endpoint - Coming soon"})
			})

		}

		employees := v1.Group("/employees")
		{
			employees.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Get employees - Coming soon"})
			})
		}

		payroll := v1.Group("/payroll")
		{
			payroll.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Payroll management - Coming soon"})
			})
		}

		profile := v1.Group("/profile")
		{
			if businessProfileHandler != nil {
				businessProfileHandler.RegisterRoutes(profile)
			}
			if profileHandler != nil {
				profileHandler.RegisterRoutes(profile)
			}
		}
	}

	// Also expose non-versioned /api routes for compatibility
	api := r.Group("/api")
	{
		// Profile routes
		profile := api.Group("/profile")
		{
			if businessProfileHandler != nil {
				businessProfileHandler.RegisterRoutes(profile)
			}
			if profileHandler != nil {
				profileHandler.RegisterRoutes(profile)
			}
		}
	}

	fmt.Println(" VestRoll Backend starting on :8080")
	fmt.Println(" Health check: http://localhost:8080/health")
	fmt.Println(" API Base: http://localhost:8080/api/v1")

	if redisClient != nil {
		fmt.Println(" Using Redis backend (real or embedded)")
		fmt.Println(" OTP Endpoints:")
		fmt.Println("   POST /api/v1/auth/send-otp")
		fmt.Println("   POST /api/v1/auth/verify-otp")
		fmt.Println(" Profile Endpoints:")
		fmt.Println("   POST /api/v1/profile/business-details")
		fmt.Println("   POST /api/profile/business-details")
		fmt.Println("   POST /api/v1/profile/account-type")
		fmt.Println("   POST /api/v1/profile/personal-details")
		fmt.Println("   POST /api/v1/profile/address")
		fmt.Println("   GET  /api/v1/profile/status?user_id=")
		fmt.Println(" PIN Endpoints:")
		fmt.Println("   POST /api/v1/auth/setup-pin")
		fmt.Println("   POST /api/v1/auth/login-pin")
		fmt.Println(" Email Verification:")
		fmt.Println("   POST /api/auth/verify-email")
	} else {
		fmt.Println(" OTP endpoints disabled (Redis not available)")
		fmt.Println(" Profile endpoints disabled (Redis not available)")
	}

	log.Fatal(r.Run(":8080"))
}
