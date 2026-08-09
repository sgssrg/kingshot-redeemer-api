package player

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
	"resty.dev/v3"
)

// UpdateAllPlayer updates real-time gamedata at the momement of all players in the database and makes a local copy to your database with the current data available ingame.
// @Summary Updates the Db with latest Player
// @Description Loads all players from the database and update every single of them, streaming per-player results as Server-Sent Events (event: update-res) followed by a summary (event: update-fin) about the rows it replaced with recent information. Rate limited.
// @Tags Player
// @Produce text/event-stream
// @Success 200 {string} string "update-res and update-fin event" example(event: update-res\ndata: {"updated":true,"player":{"pid":123,"kid":789,"dNname":"updated-player-name","pfp":"<updated-player-pfp-url>,"alliance":"DES"}}\n\n event: update-fin\ndata: {"count":123 })
// @Router /player/update/all [patch]

func UpdateAllPlayer(dbApp *db.App) echo.HandlerFunc {
	return func(c *echo.Context) error {

		reqContext := c.Request().Context()
		res := c.Response()

		// 1. Set required SSE headers
		res.Header().Set("Content-Type", "text/event-stream")
		res.Header().Set("Cache-Control", "no-cache")
		res.Header().Set("Connection", "keep-alive")
		res.Header().Set("X-Accel-Buffering", "no") // Prevents Nginx/reverse proxies from buffering the stream

		// Ensure the underlying ResponseWriter supports HTTP flushing
		flusher, ok := res.(http.Flusher)
		if !ok {
			slog.Error("Streaming unsupported: response writer does not implement http.Flusher")
			return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported")
		}

		client := resty.New()
		defer client.Close()

		metaJsonData := model.UpdateMetaSSE{
			Count: 0,
		}

		ctx := context.Background()
		p, _ := dbApp.Queries.GetAllPlayers(ctx)
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-reqContext.Done():
					return
				case <-ticker.C:
					fmt.Fprintf(res, ": keep-alive\n\n")
					flusher.Flush()
				}
			}
		}()

		for _, pi := range p {
			select {
			case <-reqContext.Done():
				slog.Info("Client disconnected, stopping process")
				return nil
			default:
			}

			// Fetching player using Scraper
			// /scraper/player/

			var pD model.PlayerInfo
			pid := strconv.Itoa(int(pi.Pid))
			slog.Info("Updating Started for - " + pid)
			url := os.Getenv("host") + "/scraper/player/" + pid
			pF, _ := client.R().SetResult(&pD).SetRetryCount(3).SetHeader("Content-Type", "application/json").Get(url)
			if pF.StatusCode() == 202 {
				p, _ := dbApp.Queries.UpdatePlayer(ctx, db.UpdatePlayerParams{
					Pid:      pi.Pid,
					Kid:      int64(pD.Kid),
					Dname:    sql.NullString{String: pD.Dname, Valid: true},
					Pfp:      sql.NullString{String: pD.Pfp, Valid: true},
					Alliance: sql.NullString{String: pD.Alliance, Valid: true},
				})
				fmt.Println(p)
				slog.Info("Player Updated with PiD - " + pid)
				updateEvent := model.UpdatePlayerSSE{
					Updated: true,
					Player: model.PlayerInfo{
						Pid:      uint(p.Pid),
						Kid:      uint(p.Kid),
						Dname:    p.Dname.String,
						Pfp:      p.Pfp.String,
						Alliance: p.Alliance.String,
					},
				}

				// 2. Marshal payload to JSON
				jsonData, err := json.Marshal(updateEvent)
				if err != nil {
					slog.Error("Failed to marshal event", "error", err)
					continue
				}

				// 3. Format as standard SSE payload: "data: <json>\n\n"
				if _, err := fmt.Fprintf(res, "event: update-res\ndata: %s\n\n", jsonData); err != nil {
					return err
				}

				// 4. Immediately flush buffer to the client
				flusher.Flush()
				// sleep(1)
				metaJsonData.Count++
			} else {
				continue
			}

		}
		metaJsonBytes, err := json.Marshal(metaJsonData)
		if err != nil {
			slog.Error("Failed to marshal event", "error", err)
			return nil
		}

		// FIX 2: Fixed duplicate "data: data:" -> "data: %s\n\n"
		fmt.Fprintf(res, "event: update-fin\ndata: %s\n\n", metaJsonBytes)
		flusher.Flush()

		return nil
	}
}

// UpdatePlayer updates real-time gamedata at the momement of a single players in the database.
// @Summary Updates the Db with latest Player Data of that player with same PlayerID
// @Description Updates a single Player with the playerID provided in cia teh param
// @Tags Player
// @Produce application/json
// @Param pid path int true "Player ID"
// @Success 200 {object} model.UpdatePlayerSSE
// @Router /player/update/ [patch]
func UpdatePlayer(dbApp *db.App) echo.HandlerFunc {
	return func(c *echo.Context) error {
		pidStr := c.Param("pid")
		pidInt, err := strconv.Atoi(pidStr)
		if err != nil {
			fmt.Printf(err.Error())
		}
		ctx := context.Background()
		client := resty.New()
		defer client.Close()
		// Fetching Player Details via Self API

		var pD model.PlayerInfo
		slog.Info("Updating Started for - " + pidStr)
		url := os.Getenv("host") + "/scraper/player/" + pidStr
		pF, _ := client.R().SetResult(&pD).SetRetryCount(3).SetHeader("Content-Type", "application/json").Get(url)
		if pF.StatusCode() == 202 {
			p, _ := dbApp.Queries.UpdatePlayer(ctx, db.UpdatePlayerParams{
				Pid:      int64(pidInt),
				Kid:      int64(pD.Kid),
				Dname:    sql.NullString{String: pD.Dname, Valid: true},
				Pfp:      sql.NullString{String: pD.Pfp, Valid: true},
				Alliance: sql.NullString{String: pD.Alliance, Valid: true},
			})
			fmt.Println(p)
			slog.Info("Player Updated with PiD - " + pidStr)
			updateEvent := model.UpdatePlayerSSE{
				Updated: true,
				Player: model.PlayerInfo{
					Pid:      uint(p.Pid),
					Kid:      uint(p.Kid),
					Dname:    p.Dname.String,
					Pfp:      p.Pfp.String,
					Alliance: p.Alliance.String,
				},
			}

			c.JSON(http.StatusAccepted, updateEvent)
		}

		return nil
	}

}
func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Second)
}
