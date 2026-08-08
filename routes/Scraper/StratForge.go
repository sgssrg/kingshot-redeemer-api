package scraper

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	"gitlab.com/ribonin/apis/kingshot-redeem/model"
	"gitlab.com/ribonin/apis/kingshot-redeem/routes/Scraper/lib"
)

// @Summary Finds a Player using PlayerID
// @Description Scrapes data from stratforge.tools for Player Details (robots.txt dated 08-10-2026)
// @Tags Scraper
// @Param player path int true "Player ID to Fetch"
// @Produce json
// @Success 202 {object} model.ScrapePlayerInfo "Created Player examples" example({"pid":123,"kid":123,"dName":"<player-name>","pfp":"player-pfp-link","alliance":"<player-alliance-tag>"})
// @Failure 409 {object} model.CustomScrapePlayerErrInfo{} "Conflict error" example({"pid":123,"type":0,"message":"What caused the error"})
// @Failure 500 {object} model.CustomScrapePlayerErrInfo{} example({"pid":123,"type":1,"message":"Player unable to be scraped / PlayerID is wrong"})

func StratForgePlayerScraperAPI(c *echo.Context) error {
	fid := c.Param("fid")
	slog.Info("Fetch Started for fid -" + fid)
	pidInt, _ := strconv.Atoi(fid)

	pInfo, errInfo, err := lib.StratForgePlayerScraper(fid)
	if err != nil {
		slog.Error("Error scraping player", "fid", fid, "err", err)
		return c.JSON(http.StatusInternalServerError, model.CustomScrapePlayerErrInfo{
			Pid:     pidInt,
			Type:    1,
			Message: "Player unable to be scraped / PlayerID is wrong",
		})
	}

	if errInfo != nil {
		return c.JSON(http.StatusConflict, errInfo)
	}
	slog.Info("Player ", "Fetched", pInfo)
	return c.JSON(http.StatusAccepted, pInfo)
}
