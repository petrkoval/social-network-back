package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type zerologPgxTracer struct {
	logger *zerolog.Logger
}

func (z *zerologPgxTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	z.logger.Debug().
		Str("query", data.SQL).
		Str("args", fmt.Sprintf("%v", data.Args)).
		Msg("tracing pgx query start")

	return ctx
}

func (z *zerologPgxTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	z.logger.Debug().
		Str("command tag", data.CommandTag.String()).
		Str("err", fmt.Sprintf("%v", data.Err)).
		Msg("tracing pgx query end")
}
