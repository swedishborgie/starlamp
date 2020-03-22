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

type SleepState int

const (
	StateUnknown SleepState = 0
	StateAwake   SleepState = 1
	StateAsleep  SleepState = 2
)

var (
	awakeTime       = "07:00:00"
	asleepTime      = "18:00:00"
	todayAwakeTime  time.Time
	todayAsleepTime time.Time
	awakeColor      = lightctl.LightStateGreen
	asleepColor     = lightctl.LightStateBlue
	currentState    = StateUnknown
)

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
	server.POST("/awake/:time", setAwakeTime)
	server.GET("/awake", getAwakeTime)
	server.POST("/asleep/:time", setAsleepTime)
	server.GET("/asleep", getAsleepTime)
	server.GET("/color", getColor)
	server.POST("/color/:color", setColor)
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
				if now.Day() != todayAwakeTime.Day() {
					recalculate()
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
	currentState = StateUnknown
	// Figure out sleep states
	if err := recalculate(); err != nil {
		log.Fatalf("problem calculating times: %s", err)
	}
}

func getStatus(c echo.Context) error {
	status := &struct {
		AwakeTime      string
		AsleepTime     string
		NextAwakeTime  time.Time
		NextAsleepTime time.Time
		AwakeColor     lightctl.LightState
		AsleepColor    lightctl.LightState
		CurrentState   SleepState
		CurrentColor   lightctl.LightState
	}{
		awakeTime,
		asleepTime,
		todayAwakeTime,
		todayAsleepTime,
		awakeColor,
		asleepColor,
		currentState,
		lightctl.GetState(),
	}
	return c.JSON(http.StatusOK, status)
}

func resetColor(c echo.Context) error {
	reset()
	return c.NoContent(http.StatusOK)
}

func getAsleepTime(c echo.Context) error {
	return c.JSON(http.StatusOK, asleepTime)
}

func getAwakeTime(c echo.Context) error {
	return c.JSON(http.StatusOK, awakeTime)
}

func setAwakeTime(c echo.Context) error {
	tmp := awakeTime
	awakeTime = c.Param("time")
	if err := recalculate(); err != nil {
		awakeTime = tmp
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.NoContent(http.StatusOK)
}

func setAsleepTime(c echo.Context) error {
	tmp := asleepTime
	asleepTime = c.Param("time")
	if err := recalculate(); err != nil {
		asleepTime = tmp
		return c.JSON(http.StatusBadRequest, err)
	}
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
	oldState := currentState
	if now.Before(todayAwakeTime) || now.After(todayAsleepTime) {
		currentState = StateAsleep
	} else {
		currentState = StateAwake
	}

	//If the state changed, changed the color
	if oldState != currentState {
		if currentState == StateAwake {
			lightctl.SetState(awakeColor)
		} else {
			lightctl.SetState(asleepColor)
		}
	}
}

func recalculateTimes() error {
	awake, err := time.Parse("15:04:05", awakeTime)
	if err != nil {
		return err
	}
	asleep, err := time.Parse("15:04:05", asleepTime)
	if err != nil {
		return err
	}
	now := time.Now()
	todayAwakeTime = time.Date(now.Year(), now.Month(), now.Day(), awake.Hour(), awake.Minute(), awake.Second(), 0, time.Local)
	todayAsleepTime = time.Date(now.Year(), now.Month(), now.Day(), asleep.Hour(), asleep.Minute(), asleep.Second(), 0, time.Local)

	return nil
}

func getColor(c echo.Context) error {
	return c.JSON(http.StatusOK, lightctl.GetState())
}

func setColor(c echo.Context) error {
	state := lightctl.ParseLightState(c.Param("color"))
	lightctl.SetState(state)
	return c.NoContent(http.StatusOK)
}

func (s SleepState) String() string {
	switch s {
	case StateAwake:
		return "awake"
	case StateAsleep:
		return "asleep"
	default:
		return "unknown"
	}
}
