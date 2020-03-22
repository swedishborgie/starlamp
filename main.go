package main

//go:generate pkger

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/markbates/pkger"
	"github.com/swedishborgie/starlamp/lightctl"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var state = State{
	AwakeTime:   "07:00:00",
	AsleepTime:  "18:00:00",
	AwakeColor:  lightctl.LightStateGreen,
	AsleepColor: lightctl.LightStateBlue,
}

func main() {
	reset()
	startTicker()

	server := echo.New()
	server.Use(middleware.CORS())
	server.GET("/", func(c echo.Context) error {
		idx, err := pkger.Open("/html/index.html")
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(idx)
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, string(body))
	})
	server.GET("/status", getStatus)
	server.POST("/status", setStatus)
	server.POST("/reset", resetColor)
	server.GET("/*", getStaticAsset)
	pkger.Include("/html")
	log.Fatalf("failed to start server: %s", server.Start(":8080"))
}

func startTicker() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case now := <-ticker.C:
				// If we're on a new day, we need to recalculate times.
				if now.Day() != state.NextAwakeTime.Day() {
					if err := recalculate(); err != nil {
						log.Printf("problem recalculating times: %s", err)
					}
				} else {
					recalculateState()
				}
			}
		}
	}()
}

func getStaticAsset(c echo.Context) error {
	filePath := url.PathEscape(c.Param("*"))
	file, err := pkger.Open("/html/" + filePath)
	if err == os.ErrNotExist {
		return c.NoContent(http.StatusNotFound)
	}
	return c.Stream(http.StatusOK, getMimeType(filePath), file)
}

func getMimeType(filePath string) string {
	filePath = strings.ToLower(filePath)
	if strings.HasSuffix(filePath, ".png") {
		return "image/png"
	} else if strings.HasSuffix(filePath, ".jpg") {
		return "image/png"
	} else if strings.HasSuffix(filePath, ".htm") || strings.HasSuffix(filePath, ".html") {
		return "text/html"
	} else if strings.HasSuffix(filePath, ".css") {
		return "text/css"
	} else if strings.HasSuffix(filePath, ".js") {
		return "application/javascript"
	}
	return "application/octet-stream"
}

func reset() {
	lightctl.Reset()
	state.CurrentState = StateUnknown
	// Figure out sleep states
	if err := recalculate(); err != nil {
		log.Fatalf("problem calculating times: %s", err)
	}
}

func getStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, state)
}

func setStatus(c echo.Context) error {
	if err := c.Bind(&state); err != nil {
		return err
	}
	if err := recalculate(); err != nil {
		return err
	}
	if state.CurrentColor != lightctl.GetState() {
		lightctl.SetState(state.CurrentColor)
	}
	return c.JSON(http.StatusOK, state)
}

func resetColor(c echo.Context) error {
	reset()
	return c.NoContent(http.StatusOK)
}

func recalculate() error {
	if err := recalculateTimes(); err != nil {
		return err
	}
	recalculateState()
	return nil
}

func recalculateState() {
	now := time.Now()
	oldState := state.CurrentState
	if now.Before(state.NextAwakeTime) || now.After(state.NextAsleepTime) {
		state.CurrentState = StateAsleep
	} else {
		state.CurrentState = StateAwake
	}

	//If the state changed, changed the color
	if oldState != state.CurrentState {
		if state.CurrentState == StateAwake {
			lightctl.SetState(state.AwakeColor)
		} else {
			lightctl.SetState(state.AsleepColor)
		}
	}
}

func recalculateTimes() error {
	awake, err := time.Parse("15:04:05", state.AwakeTime)
	if err != nil {
		return err
	}
	asleep, err := time.Parse("15:04:05", state.AsleepTime)
	if err != nil {
		return err
	}
	now := time.Now()
	state.NextAwakeTime = time.Date(now.Year(), now.Month(), now.Day(), awake.Hour(), awake.Minute(), awake.Second(), 0, time.Local)
	state.NextAsleepTime = time.Date(now.Year(), now.Month(), now.Day(), asleep.Hour(), asleep.Minute(), asleep.Second(), 0, time.Local)

	return nil
}
