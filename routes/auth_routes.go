package routes

import (
	"modul4crud/services"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures authentication and user management routes
// Note: Public routes (register/login) are now defined in routes.go
// to ensure they are registered BEFORE the protected API group
func SetupAuthRoutes(api fiber.Router, authService *services.AuthService) {
	// This function is no longer used as auth routes are now defined directly in routes.go
	// Keeping this file for future expansion or reference
	// All routes are now in routes.go:
	// - Public: /api/register, /api/login, /auth/register, /auth/login
	// - Protected: /api/profile, /api/logout
	// - Admin: /api/users/*
}
