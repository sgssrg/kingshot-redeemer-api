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
// @Success 202 {object} model.DeletePlayerResponseSuccess "Player Delete"
// @Failure 400 {object} model.DeletePlayerValidationError "Invalid PiD"
// @Failure 500 {object} model.DeletePlayerResponseFailure "Delete Failed"
// @Router /player/del/{pid} [delete]
func DeletePlayer(dbApp *db.App) echo.HandlerFunc {
	return func(c *echo.Context) error {
		pid := c.Param("pid")
		slog.Info("DeletePlayer called", "pid", pid)

		pidInt, err := strconv.Atoi(pid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.DeletePlayerValidationError{
				WrongField: "pid",
				Value:      pid,
				Message:    "Invalid PiD, cannot be converted to Integer",
			})
		}

		ctx := context.Background()
		pD, err := dbApp.Queries.DeletePlayer(ctx, int64(pidInt))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.DeletePlayerResponseFailure{
				Pid:     pid,
				Deleted: false,
				Message: err.Error(),
			})
		}
		// Row deleted
		return c.JSON(http.StatusAccepted, model.DeletePlayerResponseSuccess{
			Pid:     pid,
			Deleted: true,
			Player: model.PlayerInfo{
				Pid:      uint(pidInt),
				Kid:      uint(pD.Kid),
				Dname:    pD.Dname.String,
				Pfp:      pD.Pfp.String,
				Alliance: pD.Alliance.String,
			},
			Message: "Player deleted successfully",
		})
	}
}
