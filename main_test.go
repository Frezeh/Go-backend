
package main

import (
  "net/http/httptest"
  "net/http"
  "testing"
  // "bytes"
  // "encoding/json"
  "fmt"

  "github.com/Frezeh/Go-backend/handlers"
  "github.com/gofiber/fiber/v2"
  // "github.com/gofiber/fiber/v2/middleware/cors"
  "github.com/stretchr/testify/assert" // add Testify package
)

type Tests struct {
  name string
  server *httptest.Server
  // response *http.Response
  expectedCode int 
}

func TestPingRoute(t *testing.T) {
  tests := []struct {
    description  string 
    route        string 
    expectedCode int 
  }{
    {
      description:  "get HTTP status 200",
      route:        "/api",
      expectedCode: 200,
    },
    {
      description:  "get HTTP status 404, when route is not exists",
      route:        "/not-found",
      expectedCode: 404,
    },
  }
  app := fiber.New()
  app.Get("/api", handlers.Ping)

  for _, test := range tests {
    req := httptest.NewRequest("GET", test.route, nil)
    // rec := httptest.NewRecorder()
    resp, _ := app.Test(req, 1)
    // fmt.Println(rec.Body.String())
    assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
  }
}

func TestBalanceRouteBeforeSignUp(t *testing.T) {
  tests := []struct {
    description  string 
    route        string 
    expectedCode int 
  }{
    {
      description:  "get HTTP status 400",
      route:        "/api/balance",
      expectedCode: 400,
    },
  }
  app := fiber.New()
  app.Get("/api/balance", handlers.GetBalance)

  for _, test := range tests {
    req := httptest.NewRequest("GET", test.route, nil)
    resp, _ := app.Test(req, 1)
    assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
  }
}