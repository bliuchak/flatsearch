package estate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog"

	"github.com/bliuchak/flatsearch/internal/platform/storage"
	"github.com/bliuchak/flatsearch/internal/sreality"
)

type Estate struct {
	sreality  *sreality.SrealityAPI
	pgStorage *storage.Postgres
	logger    zerolog.Logger
}

func NewEstate(sreality *sreality.SrealityAPI, pgStorage *storage.Postgres, logger zerolog.Logger) *Estate {
	return &Estate{sreality: sreality, pgStorage: pgStorage, logger: logger}
}

func (e *Estate) CompareAndSave(ctx context.Context, estates []sreality.Estate) error {
	var addCounter uint
	var newEstates []string

	for _, estate := range estates {
		extID := strconv.Itoa(estate.HashID)
		name := strings.Replace(estate.Name, "\u00a0", " ", 1) // remove nbsp from string
		url := "https://www.sreality.cz/detail/pronajem/ostatni/garazove-stani/" + estate.Seo.Locality + "/" + strconv.Itoa(estate.HashID)
		price := strconv.Itoa(int(estate.Price))

		savedEstate, err := e.pgStorage.GetEstateByExtID(ctx, extID)
		if errors.Is(err, sql.ErrNoRows) {
			// add new advert if not exists
			_, err = e.pgStorage.InsertEstate(ctx, extID, name, url, price)
			if err != nil {
				return fmt.Errorf("insert estate: %s", err)
			}
		}
		if err != nil {
			return fmt.Errorf("get estate: %s", err)
		}

		// update estate if exists
		_, err = e.pgStorage.UpdateEstate(ctx, savedEstate.ExtID, name, url, price)
		if err != nil {
			return fmt.Errorf("update estate: %s", err)
		}

		addCounter++

		newEstates = append(newEstates, url+" - "+price)
	}

	e.logger.Info().
		Uint("added", addCounter).
		Msg("save and compare status")

	return nil
}
