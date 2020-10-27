package router

import (
	db "go-authentication-boilerplate/database"
	"go-authentication-boilerplate/models"
	"go-authentication-boilerplate/util"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes func sets up all the user routes
func SetupUserRoutes() {
	USER.Post("/signup", CreateUser) // Sign Up a user
	USER.Post("/signin", LoginUser)  // Sign In a user
}

// CreateUser route registers a User into the database
func CreateUser(c *fiber.Ctx) error {
	u := new(models.User)

	if err := c.BodyParser(u); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	// validate if the email, username and password are in correct format
	errors := util.ValidateRegister(u)
	if errors.Err {
		return c.JSON(errors)
	}

	if count := db.DB.Where(&models.User{Email: u.Email}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Err, errors.Email = true, "Email is already registered"
	}
	if count := db.DB.Where(&models.User{Username: u.Username}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Err, errors.Username = true, "Username is already registered"
	}
	if errors.Err {
		return c.JSON(errors)
	}

	// Hashing the password with a random salt
	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	u.Password = string(hashedPassword)

	if err := db.DB.Create(&u).Error; err != nil {
		return c.JSON(fiber.Map{"error": true, "general": "Something went wrong, please try again later. 😕"})
	}

	// setting up the authorization cookies
	accessToken, refreshToken := util.GenerateTokens(u)
	accessCookie, refreshCookie := getAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.SendString("User registered successfully")

}

// LoginUser route logins a user in the app
func LoginUser(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	u := new(models.User)
	if res := db.DB.Where(&models.User{Email: input.Identity}).Or(&models.User{Username: input.Identity}).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// setting up the authorization cookies
	accessToken, refreshToken := util.GenerateTokens(u)
	accessCookie, refreshCookie := getAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.SendString("User succesfullt logged in")

}

func getAuthCookies(accessToken, refreshToken string) (*fiber.Cookie, *fiber.Cookie) {
	accessCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(15 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	return accessCookie, refreshCookie
}
