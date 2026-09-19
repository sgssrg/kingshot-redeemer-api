package alliance

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
)

func GetAllUniqueAllianceFromDB(dbApp *db.App) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := context.Background()
		alliances, err := dbApp.Queries.GetAllUniqueAlliance(ctx)
		if err != nil {
			// Handle database error (e.g., connection lost)
			slog.Warn(err.Error())
			return c.JSON(http.StatusInternalServerError, model.GetAllUniqueAllianceFromDBResponse{ErrorState: true, Message: "Unable to Fetch Alliances"})
		}

		if len(alliances) == 0 {
			fmt.Println("No alliances found.")
			return c.JSON(http.StatusOK, model.GetAllUniqueAllianceFromDBResponse{ErrorState: true, Message: "No Alliances or Player in DB."})
		} else {
			var result []string
			for _, a := range alliances {
				if a.Valid {
					result = append(result, a.String)
				}
			}
			return c.JSON(http.StatusOK, model.GetAllUniqueAllianceFromDBResponse{ErrorState: false, Message: "Fetched Unique Alliances List From DB.", AllianceList: result})
		}
	}
}
