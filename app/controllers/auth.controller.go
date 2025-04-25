package controllers

import (
	"fmt"
	"strings"
	"style-stamp/app/models"
	"style-stamp/config"
	"style-stamp/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SignUp(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": "Invalid credentials", "error": err.Error()})
	}

	var exists models.User
	result := config.DB.Where("username=?", user.Username).First(&exists) // check if username already exists
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": fiber.StatusNotFound, "message": "user name allready exists", "error": result.Error})
	}

	user.Password = utils.HashPassword(user.Password)
	result = config.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": fiber.StatusInternalServerError, "error": result.Error})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": fiber.StatusCreated, "message": "User signup successful"})
}

func SignIn(c *fiber.Ctx) error {
	var input models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": "Invalid credentials", "error": err.Error()})
	}

	var user models.User
	result := config.DB.Where("username=?", input.Username).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": fiber.StatusNotFound, "message": "Invalid credentials", "error": result.Error})
	}
	if !utils.ComparePassword(input.Password, user.Password) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": "Invalid credentials", "error": result.Error})
	}
	auth := utils.AuthDto{
		ID:       user.ID,
		Username: user.Username,
		Jti:      uuid.New(),
	}
	accessToken, err := utils.CreateAccessToken(auth)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": "Invalid credentials", "error": err})
	}

	refreshToken, err := utils.CreateRefreshToken(auth)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest, "message": "Invalid credentials", "error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": fiber.StatusOK, "message": "Signin successful ", "accessToken": accessToken, "refreshToken": refreshToken})
}

func GetDetails(c *fiber.Ctx) error {
	data := c.Locals("user")
	fmt.Println(data)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": fiber.StatusOK, "message": "user details fetch successful ", "data": data})
}

func LogOut(c *fiber.Ctx) error {
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
	device := models.Device{
		Jti: jti,
	}
	result := config.DB.Create(&device)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  500,
			"message": "Internal Server Error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "Logout successful.",
	})
}

func RefreshToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	claims, err := utils.VerifyRefreshToken(tokenParts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
	jti, ok := claims["jti"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid uuid"})
	}
	var device models.Device
	findDevice := config.DB.Where("jti=?", jti).First(&device)
	if findDevice.Error == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
	authData := models.Device{
		Jti: jti,
	}
	result := config.DB.Create(&authData)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failer to refresh the token"})
	}
	idStr, _ := claims["id"].(string)
	id, _ := uuid.Parse(idStr)
	username, _ := claims["username"].(string)

	authDto := utils.AuthDto{
		ID:       id,
		Username: username,
		Jti:      uuid.New(),
	}

	token, _ := utils.CreateAccessToken(authDto)
	refreshToken, _ := utils.CreateRefreshToken(authDto)

	return c.JSON(fiber.Map{"status": 200, "message": "User Logout Successful.", "accessToken": token, "refreshToken": refreshToken})
}
