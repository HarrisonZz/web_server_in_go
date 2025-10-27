package i2cdevice

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/HarrisonZz/web_server_in_go/internal/handler"
	"github.com/HarrisonZz/web_server_in_go/internal/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/io/i2c"
)

const (
	stm32Addr = 0x15

	//Registers
	LedCtrl  = 0x01
	LedQuery = 0x02

	//value
	LED_ON  = 0x01
	LED_OFF = 0x00
)

var (
	i2cDev *i2c.Device
	mu     sync.Mutex
)

func InitI2C() {

	logger.Info("I2C API initializing")

	var err error
	i2cDev, err = i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-2"}, stm32Addr)
	if err != nil {
		logger.Error("I2C Bus Open Failed !")
		return
	}
	if err := probeI2CDevice(i2cDev); err == nil {

		handler.RegisterRoute(http.MethodPost, "/led", ledHandler)
		handler.RegisterRoute(http.MethodGet, "/led", ledHandler)

	} else {
		logger.Info("Skipping /led route registration (I2C device not found)")
	}
}

func probeI2CDevice(dev *i2c.Device) error {

	err := dev.Write([]byte{})
	if err != nil {
		logger.Error(fmt.Sprintf("I2C device not responding: %v", err))
		return err
	}
	logger.Info("I2C device responded successfully.")
	return nil
}

func ledHandler(c *gin.Context) {

	switch c.Request.Method {
	case http.MethodPost:
		handleLedSet(c)
	case http.MethodGet:
		handleLedQuery(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not supported"})
	}

}

func handleLedQuery(c *gin.Context) {

	mu.Lock()
	defer mu.Unlock()

	// 讀取 1 byte
	buf := make([]byte, 1)
	if err := i2cDev.ReadReg(LedQuery, buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("I2C read failed: %v", err)})
		return
	}

	var state string
	switch buf[0] {

	case 0x00:

		state = "Off"

	case 0x01:

		state = "On"

	}

	c.JSON(http.StatusOK, gin.H{"LED State get": state})
}

func handleLedSet(c *gin.Context) {
	start := time.Now()
	var req struct {
		State string `json:"state"`
	}

	if err := c.BindJSON(&req); err != nil {
		logger.Warn(fmt.Sprintf(
			"[LED] %s invalid JSON from=%s error=%v",
			c.FullPath(),
			c.ClientIP(),
			err,
		))

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	state := strings.ToLower(req.State)
	var data byte

	switch state {
	case "on":
		data = LED_ON
	case "off":
		data = LED_OFF
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "state must be 'on' or 'off'"})

		logger.Warn(fmt.Sprintf(
			"[LED] %s invalid state='%s' from=%s",
			c.FullPath(),
			state,
			c.ClientIP(),
		))
		return
	}

	logger.Info(fmt.Sprintf(
		"[LED] Request received state=%s from=%s",
		state,
		c.ClientIP(),
	))

	mu.Lock()
	defer mu.Unlock()

	if err := i2cDev.WriteReg(LedCtrl, []byte{data}); err != nil {
		logger.Error(fmt.Sprintf(
			"[LED] WriteReg failed state=%s error=%v from=%s",
			state,
			err,
			c.ClientIP(),
		))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("I2C write failed: %v", err)})
		return
	}
	elapsed := time.Since(start)

	logger.Info(fmt.Sprintf(
		"[LED] State changed to %s via I2C duration=%v from=%s",
		state,
		elapsed,
		c.ClientIP(),
	))

	c.JSON(http.StatusOK, gin.H{"LED Status set": state})
}
