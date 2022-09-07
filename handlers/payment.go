package handlers

import (
	"time"
	// "fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/aidarkhanov/nanoid"
)

type User struct {
	ID    string    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Balance int `json:"balance" default:"0"`
}

type LoginRequest struct {
	Email    string
	Password string
}

type Amount struct {
	Amount int `json:"amount"`
}

type Balance struct {
	Status bool `json:"status"`
	Balance int `json:"balance"`
}
                      
var users []User
var currentUserId string = ""

func createJWTToken(user User) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = exp
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}

func Ping(c *fiber.Ctx) error {
	return c.SendString("OK  👍")
}


func SignUp(c *fiber.Ctx) error {
	id := nanoid.New()
	req := new(User)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid signup credentials!")
	}

	// save this info in the memory
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		ID: id,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	}

	users = append(users, *user)
	currentUserId = user.ID

	// create a jwt token
	token, exp, err := createJWTToken(*user)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
}

func Login(c *fiber.Ctx) error {
	user := User{}
	req := new(LoginRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}
	
		if req.Email == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid login credentials!")
		}
	
	// verify email and password
	for i, u := range users {
		if u.Email == req.Email {
			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
				return err
			}		
			user = users[i]
			break
		} else {
			return fiber.NewError(fiber.StatusBadRequest, "Password or E-mail do not match!")
		}
	}
		
	currentUserId = user.ID
	token, exp, err := createJWTToken(user)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
}

func Deposit(c *fiber.Ctx) error {
	req := new(Amount)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	if req.Amount == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Please input an amount!")
	}

	for i, u := range users {
		if u.ID == currentUserId {
			users[i].Balance += req.Amount 
			break
		}
	}

	return c.SendString("Deposit Successful")
}

func Transfer(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(Amount)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	// check if balance is sufficient before debiting
	for i, u := range users {
		if u.ID == currentUserId {
			if users[i].Balance < req.Amount {
				return fiber.NewError(fiber.StatusBadRequest, "Insufficient balance!")
			} else {
				users[i].Balance -= req.Amount
			}
			break
		}
	}

	// credit recipient
	for i, u := range users {
		if u.ID == id {
			users[i].Balance += req.Amount 
			break
		}
	}

	return c.SendString("Transfer Successful")
}

func TransferOut(c *fiber.Ctx) error {
	req := new(Amount)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	if req.Amount == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Please input an amount!")
	}

	// check if balance is sufficient before debiting
	for i, u := range users {
		if u.ID == currentUserId {
			if users[i].Balance < req.Amount {
				return fiber.NewError(fiber.StatusBadRequest, "Insufficient balance!")
			} else {
				users[i].Balance -= req.Amount
			}
			break
		}
	}

	return c.SendString("Transfer Successful")
}

func GetBalance(c *fiber.Ctx) error {
	var balance int

	if currentUserId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID!")
	}

	for i, u := range users {
		if u.ID == currentUserId {
			balance = users[i].Balance
			break
		}
	}

	userBalance := &Balance{
		Status: true,
		Balance: balance,
	}

	return c.JSON(userBalance)
}

func All(c *fiber.Ctx) error {
	return c.JSON(users)
}