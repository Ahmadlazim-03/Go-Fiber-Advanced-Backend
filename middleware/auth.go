package middleware

import (
	"modul4crud/models"
	"modul4crud/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ValidateJWT middleware untuk validasi token JWT
func ValidateJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak ditemukan",
			})
		}

		// Extract token dari "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid",
			})
		}

		tokenString := tokenParts[1]

		// Parse dan validasi token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid",
			})
		}

		// Set user info ke context untuk digunakan di handler selanjutnya
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RequireRole middleware untuk validasi role-based access control
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "User tidak terautentikasi",
			})
		}

		role := userRole.(string)

		// Check apakah role user ada dalam allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error":          "Akses ditolak: role tidak memiliki permission",
			"required_roles": allowedRoles,
			"user_role":      role,
		})
	}
}

// RequireAdmin middleware khusus untuk admin only
func RequireAdmin() fiber.Handler {
	return RequireRole(models.RoleAdmin)
}

// RequireUser middleware khusus untuk user role (jika diperlukan)
func RequireUser() fiber.Handler {
	return RequireRole(models.RoleUser)
}

// RequireAdminOrUser middleware untuk admin atau user (semua yang authenticated)
func RequireAdminOrUser() fiber.Handler {
	return RequireRole(models.RoleAdmin, models.RoleUser)
}

// GetUserFromContext helper function untuk mengambil user info dari context
func GetUserFromContext(c *fiber.Ctx) *models.JWTClaims {
	return &models.JWTClaims{
		UserID:   c.Locals("user_id").(int),
		Username: c.Locals("username").(string),
		Role:     c.Locals("role").(string),
	}
}
