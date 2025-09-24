package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/paulmach/orb/encoding/wkb"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

func (r *Repository) SaveMessage(ctx context.Context, mes *model.ParsedMessage) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadUncommitted,
	})
	if err != nil {
		return err
	}
	query := `
        INSERT INTO messages(
                             region,
            sid, dof, atd, ata, dep_coords_normalize, arr_coords_normalize,
            dep_coordinate, arr_coordinate, arr_region_rf, opr, reg, typ, rmk, min_alt, max_alt
        )
        VALUES ((SELECT d.gid FROM district_shapes as d WHERE st_contains(d.geom,ST_SetSRID(ST_GeomFromWKB($7),0))),$1, $2, $3, $4, $5, $6, ST_GeomFromWKB($7), ST_GeomFromWKB($8), $9, $10, $11, $12, $13, $14, $15)
        ON CONFLICT (sid,atd, dep_coordinate, arr_coordinate) DO NOTHING;
    `
	slog.Info("Executing insert query", "sid", mes.SID)
	if mes.SID == "7771464892" {
		fmt.Println(mes.SID)
	}
	_, err = tx.Exec(ctx, query,
		mes.SID, mes.DOF, mes.ATD, mes.ATA, mes.DepCoords, mes.ArrCoords,
		wkb.Value(mes.DepLatLon), wkb.Value(mes.ArrLatLon), mes.ArrRegionRF,
		mes.OPR, mes.REG, mes.TYP, mes.RMK, mes.MinAlt, mes.MaxAlt)
	if err != nil {
		tx.Rollback(ctx)
		slog.Error("Failed to execute query", "sid", mes.SID, "err", err)
		return err
	}
	if len(mes.ZoneLatLon) > 0 {
		insertFlightCood := `
            INSERT INTO flight_coordinates(sid, coordinate) 
            VALUES ($1, ST_GeomFromWKB($2))
        `
		for _, coord := range mes.ZoneLatLon {
			if coord[1] > 90 || coord[1] < -90 || coord[0] > 180 || coord[0] < -180 {
				slog.Error("Invalid coordinate", "sid", mes.SID, "coord", coord)
				continue
			}
			_, err := tx.Exec(ctx, insertFlightCood, mes.SID, wkb.Value(coord))
			if err != nil {
				tx.Rollback(ctx)
				slog.Error("Failed to insert flight coordinates", "sid", mes.SID, "err", err)
				return err
			}
		}
	}
	tx.Commit(ctx)
	return nil
}
