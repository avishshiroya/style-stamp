package middlewares

import (
	"fmt"
	"strings"
	"style-stamp/app/models"
	"style-stamp/config"
	"style-stamp/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Jti      string `json:"jti"`
}

func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  401,
				"message": "Token Required",
			})
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" || tokenParts[1] == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  401,
				"message": "Not valid format of token",
			})
		}
		token := tokenParts[1]
		claims, err := utils.VerifyAccessToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  401,
				"message": "Invalid or expired token",
			})
		}
		jti, ok := claims["jti"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  401,
				"message": "Invalid or expired token",
			})
		}
		var device models.Device
		result := config.DB.Where("jti=?", jti).First(&device)
		if result.Error == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  401,
				"message": "Token is expired",
			})
		}

		userIdStr, _ := claims["id"].(string)
		userId, _ := uuid.Parse(userIdStr)
		fmt.Print(userId)
		var user UserResponse
		userResult := config.DB.Model(&models.User{}).Select("id,username").Where("id=?", userId).Find(&user)
		if userResult.Error != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  401,
				"message": "User Not Found",
			})
		}
		fmt.Print(user)
		user.Jti = jti
		c.Locals("user", user)
		return c.Next()
	}
}
