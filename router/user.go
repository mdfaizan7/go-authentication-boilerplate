package router

import (
	db "notion-clone/server/database"
	"notion-clone/server/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes func sets up all the user routes
func SetupUserRoutes() {
	USER.Post("/", CreateUser)
}

// CreateUser route registers a User into the database
func CreateUser(c *fiber.Ctx) error {
	u := new(models.User)

	if err := c.BodyParser(u); err != nil {
		return c.JSON(`{"error":true, "input":"Please review your input"`)
	}

	// validate if the email, username and password are in correct format
	errors := ValidateRegister(u)
	if errors.Err {
		return c.JSON(errors)
	}

	// validate if email and username are unique
	user := new(models.User)
	db.DB.Where(&models.User{Email: u.Email}).Find(&user)
	if u.Email == user.Email {
		errors.Err, errors.Email = true, "Email is already registered"
	}
	user = new(models.User)
	db.DB.Where(&models.User{Username: u.Username}).Find(&user)
	if u.Username == user.Username {
		errors.Err, errors.Username = true, "Username is already registered"
	}

	if errors.Err {
		return c.JSON(errors)
	}

	// Hashing the password with the default cost of 10
	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		panic(err)
	}
	u.Password = string(hashedPassword)

	db.DB.Create(&u)

	return c.JSON(u)
}

//
