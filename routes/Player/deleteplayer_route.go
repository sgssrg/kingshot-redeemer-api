package player

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
)

// DeletePlayer removes a player by PiD.
// @Summary Delete a player
// @Description Removes a player from the database by player ID.
// @Tags Player
// @Produce json
// @Param pid path int true "Player ID"
// @Success 202 {object} model.DeletePlayerResponse "Player deleted" example({"pid":123,"deleted":true})
// @Failure 500 {object} model.DeletePlayerErrorResponse "Delete failed" example({"pid":123,"deleted":false,"error":"database error"})
// @Router /player/del/{pid} [delete]
func DeletePlayer(dbApp *db.App) echo.HandlerFunc {
	return func(c *echo.Context) error {
		pid := c.Param("pid")

		slog.Info(pid)
		ctx := context.Background()
		pidInt, _ := strconv.Atoi(pid)

		rows, err := dbApp.Queries.DeletePlayer(ctx, int64(pidInt))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.DeletePlayerErrorResponse{
				Pid:     pid,
				Deleted: false,
				Error:   err.Error(),
			})
		}

		if rows == 0 {
			// Nothing deleted
			return c.JSON(http.StatusNotFound, model.DeletePlayerResponse{
				Pid:     pid,
				Deleted: false,
				Message: "There was no Player to delete with this PlayerID",
			})
		}

		// Row deleted
		return c.JSON(http.StatusAccepted, model.DeletePlayerResponse{
			Pid:     pid,
			Deleted: true,
			Message: "Player Deleted Sucessfully!",
		})
	}
}
