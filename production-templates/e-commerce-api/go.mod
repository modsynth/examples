module github.com/modsynth/e-commerce-api

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/joho/godotenv v1.5.1
	github.com/stripe/stripe-go/v76 v76.11.0
	golang.org/x/crypto v0.17.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/driver/sqlite v1.5.4
	gorm.io/gorm v1.25.5
)

// Modsynth modules (will be available once published)
// Uncomment these when the modules are actually published
// require (
// 	github.com/modsynth/auth-module v0.1.0
// 	github.com/modsynth/db-module v0.1.0
// 	github.com/modsynth/cache-module v0.1.0
// 	github.com/modsynth/logging-module v0.1.0
// 	github.com/modsynth/file-storage-module v0.1.0
// 	github.com/modsynth/payment-module v0.1.0
// 	github.com/modsynth/search-module v0.1.0
// 	github.com/modsynth/monitoring-module v0.1.0
// )
