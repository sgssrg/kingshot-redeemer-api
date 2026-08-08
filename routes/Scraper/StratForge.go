package scraper

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"

	"gitlab.com/ribonin/apis/kingshot-redeem/routes/Scraper/lib"
)

func StratForgePlayerScraperAPI(c *echo.Context) error {
	fid := c.Param("fid")
	slog.Info("Fetch Started for fid -" + fid)

	pInfo, errInfo, err := lib.StratForgePlayerScraper(fid)
	if err != nil {
		slog.Error("Error scraping player", "fid", fid, "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Player unable to be scraped",
		})
	}

	if errInfo != nil {
		return c.JSON(http.StatusConflict, errInfo)
	}

	return c.JSON(http.StatusAccepted, pInfo)
}
