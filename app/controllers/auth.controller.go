package controllers

import (
	"style-stamp/app/models"
	"style-stamp/config"
	"style-stamp/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SignUp(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest,"message": "Invalid credentials", "error": err.Error()})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": fiber.StatusBadRequest,"message": "Invalid credentials", "error": err.Error()})
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
		ID: user.ID,
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
