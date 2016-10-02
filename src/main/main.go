package main

import(
    "log"
    "net/http"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "github.com/labstack/echo/engine/standard"
)
// site url for the html template
// https://html5up.net/Landed

func hello(c echo.Context) error {
    return c.String(http.StatusOK, "yallo!\n")
}

func adminMain(c echo.Context) error {
    return c.String(http.StatusOK, "admin area")
}

func login(c echo.Context) error {
    userName := c.QueryParam("username")
    password := c.QueryParam("password")

    if userName == "bluebot" && password == "topsecret" {
        cookie := new(echo.Cookie)
        cookie.SetName("blue_bot")
        cookie.SetValue("secret_token")
        cookie.SetExpires(time.Now().Add(time.Second * 30))
        c.SetCookie(cookie)
        return c.String(http.StatusOK, "you have a cookie now")
    }

    return c.String(http.StatusUnauthorized, "somthing is fishy here, who are you?")
}

func auth() echo.MiddlewareFunc {
    return middleware.BasicAuth(func(user, pass string) bool {
        if user == "bot" && pass == "secret" {
            return true
        }

        return false
    })
}

func checkPass(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        if false {

            return next(c)
        }

        return c.JSON(401, map[string]string{
            "error": "go back you are not the one",
        })
    }
}

func logReq(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        log.Printf("Request type [%s] %s\n", c.Request().Method(), c.Request().URL().Path())
        return next(c)
    }
}

func sendName(c echo.Context) error {
    name := c.Param("name")

    params := c.QueryParams()
    b, _ := json.MarshalIndent(params, "", "   ")

    return c.String(http.StatusOK, fmt.Sprintf("this the value of name: %s\nthis is the values of all the params: %s\n", name, string(b)))
}

func main() {

    e := echo.New()

    e.Use(
        middleware.Secure(),
    )

    e.File("/", "./index.html")
    e.Static("/assets", "assets")
    e.Static("/images", "images")
    e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
        Level: 5,
    }))

    e.Use(logReq)
    g := e.Group("/admin", checkPass)
    g.GET("/main", adminMain)

    e.GET("/hello", hello, checkPass)
    e.GET("/login", login)
    e.GET("/main/:name", sendName)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    log.Printf("Starting server on port: %s\n", port)
    e.Run(standard.New(":" + port))
}