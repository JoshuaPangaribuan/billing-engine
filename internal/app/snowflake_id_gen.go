package app

import (
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkguid"
	"github.com/hashicorp/go-multierror"
)

func (app *App) initSnowflakeGen() {
	snowflakeGen, err := pkguid.NewSnowflake()
	if err != nil {
		app.err = multierror.Append(app.err, err)
	}

	app.snowflakeGen = snowflakeGen
}
