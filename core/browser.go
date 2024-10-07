package core

import (
	"EtherDrop/tools"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gookit/config/v2"
)

func initializeBrowser() *rod.Browser {
	extensionPath, _ := filepath.Abs("./extension/mini-app-android-spoof")
	browserPath := config.String("BROWSER_PATH")

	var launchOptions string

	if browserPath != "" {
		launchOptions = launcher.New().
			Bin(browserPath).
			Set("load-extension", extensionPath).
			Headless(true).
			Set("no-sandbox", "").
			Set("disable-gpu", "").
			Set("disable-background-services", "").
			Set("no-default-browser-check", "").
			Set("disable-popup-blocking", "").
			Set("disable-hardware-acceleration", "").
			Set("disable-application-cache", "").
			Set("disable-media-cache", "").
			Set("disable-infobars", "").
			Set("autoplay-policy", "no-user-gesture-required").
			Set("disable-background-timer-throttling", "").
			Set("disable-notifications", "").
			Set("log-level", "error").
			Set("disable-sync", "").
			Set("disable-accessibility", "").
			Set("disable-web-security", "").
			MustLaunch()
	} else {
		launchOptions = launcher.New().
			Set("load-extension", extensionPath).
			Headless(true).
			Set("no-sandbox", "").
			Set("disable-gpu", "").
			Set("disable-background-services", "").
			Set("no-default-browser-check", "").
			Set("disable-popup-blocking", "").
			Set("disable-hardware-acceleration", "").
			Set("disable-application-cache", "").
			Set("disable-media-cache", "").
			Set("disable-infobars", "").
			Set("autoplay-policy", "no-user-gesture-required").
			Set("disable-background-timer-throttling", "").
			Set("disable-notifications", "").
			Set("log-level", "error").
			Set("disable-sync", "").
			Set("disable-accessibility", "").
			Set("disable-web-security", "").
			MustLaunch()
	}

	browser := rod.New().ControlURL(launchOptions).MustConnect()

	return browser
}

func (c *Client) checkElement(page *rod.Page, selector string) bool {
	defer func() {
		if r := recover(); r != nil {
			tools.Logger("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.account.phone, r))
		}
	}()

	for attempt := 1; attempt <= 2; attempt++ {
		_, err := page.Timeout(5 * time.Second).Element(selector)

		if err == nil {
			return true
		} else if errors.Is(err, &rod.ElementNotFoundError{}) {
			if attempt == 3 {
				tools.Logger("warning", fmt.Sprintf("| %s | Element %v not found after %d attempts", c.account.phone, selector, attempt))
				return false
			}
		} else {
			panic(err)
		}
	}

	return false
}

func (c *Client) navigate(page *rod.Page, url string) {
	defer func() {
		if r := recover(); r != nil {
			tools.Logger("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.account.phone, r))
		}
	}()

	page.Timeout(3 * time.Second).Navigate(url)
	page.MustWaitLoad()
	page.MustWaitRequestIdle()
}

func (c *Client) clickElement(page *rod.Page, selector string) {
	defer func() {
		if r := recover(); r != nil {
			tools.Logger("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.account.phone, r))
		}
	}()

	c.checkElement(page, selector)

	page.Timeout(3 * time.Second).MustElement(selector).MustWaitVisible()

	page.Timeout(3 * time.Second).MustElement(selector).MustClick()

	page.MustWaitRequestIdle()
}

func (c *Client) inputText(page *rod.Page, value string, selector string) {
	defer func() {
		if r := recover(); r != nil {
			tools.Logger("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.account.phone, r))
		}
	}()

	c.checkElement(page, selector)

	page.Timeout(3 * time.Second).MustElement(selector).MustClick().MustInput(value)
}

func (c *Client) getText(page *rod.Page, selector string) string {
	defer func() {
		if r := recover(); r != nil {
			tools.Logger("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.account.phone, r))
		}
	}()

	c.checkElement(page, selector)

	text := page.Timeout(10 * time.Second).MustElement(selector).MustText()

	return text
}

func (c *Client) removeTextFormInput(page *rod.Page, selector string) {
	defer func() {
		if r := recover(); r != nil {
			tools.Logger("warning", fmt.Sprintf("| %s | Recovered from panic : %v", c.account.phone, r))
		}
	}()

	c.checkElement(page, selector)

	// Dapatkan elemen berdasarkan selector, periksa apakah elemen ada
	element, err := page.Element(selector)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Element not found: %v", c.account.phone, err))
		return
	}
	if element == nil {
		tools.Logger("error", fmt.Sprintf("| %s | Element is nil for selector: %s", c.account.phone, selector))
		return
	}

	// Periksa apakah elemen adalah input atau div
	tagName, err := element.Eval(`() => this.tagName.toLowerCase()`)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get tag name: %v", c.account.phone, err))
		return
	}

	switch tagName.Value.String() {
	case "input":
		// Jika elemen adalah input, hapus teks
		page.MustElement(selector).MustSelectAllText().MustInput("")
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to clear input text: %v", c.account.phone, err))
		}
	case "div":
		// Jika elemen adalah div, hapus teks menggunakan JavaScript
		_, err = element.Eval(`() => { this.textContent = ""; }`)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed To Remove Text From Input Field: %v", c.account.phone, err))
		} else {
			tools.Logger("info", fmt.Sprintf("| %s | Remove Text From Input Field", c.account.phone))
		}
	default:
		tools.Logger("info", fmt.Sprintf("| %s | The element is not an input or div, skipping", c.account.phone))
	}
}
