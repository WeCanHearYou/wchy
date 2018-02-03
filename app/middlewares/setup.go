package middlewares

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/worker"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage/postgres"
)

// Noop does nothing
func Noop() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			return next(c)
		}
	}
}

//WorkerSetup current context with some services
func WorkerSetup(logger log.Logger) worker.MiddlewareFunc {
	db := dbx.NewWithLogger(logger)
	emailer := app.NewEmailer(logger)
	return func(next worker.Job) worker.Job {
		return func(c *worker.Context) (err error) {
			trx, err := db.Begin()
			if err != nil {
				return err
			}

			c.SetServices(&app.Services{
				Tenants: postgres.NewTenantStorage(trx),
				Users:   postgres.NewUserStorage(trx),
				Ideas:   postgres.NewIdeaStorage(trx),
				Tags:    postgres.NewTagStorage(trx),
				Emailer: emailer,
			})

			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v\n%v", r, string(debug.Stack()))

					if trx != nil {
						trx.Rollback()
					}
				}
			}()
			if err = next(c); err != nil {
				panic(err)
			}
			trx.Commit()
			return err
		}
	}
}

//WebSetup current context with some services
func WebSetup(logger log.Logger) web.MiddlewareFunc {
	db := dbx.NewWithLogger(logger)
	emailer := app.NewEmailer(logger)
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			path := log.Magenta(c.Request.Method + " " + c.Request.URL.String())

			start := time.Now()
			c.Logger().Infof("%s started", path)

			trx, err := db.Begin()
			if err != nil {
				return err
			}

			c.SetActiveTransaction(trx)

			c.SetServices(&app.Services{
				Tenants: postgres.NewTenantStorage(trx),
				OAuth:   &oauth.HTTPService{},
				Users:   postgres.NewUserStorage(trx),
				Ideas:   postgres.NewIdeaStorage(trx),
				Tags:    postgres.NewTagStorage(trx),
				Emailer: emailer,
			})

			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("%v\n%v", r, string(debug.Stack()))
					c.Failure(err)
					c.Logger().Infof("%s finished in %s", path, log.Magenta(time.Since(start).String()))
					if trx != nil {
						trx.Rollback()
					}
				}
			}()

			if err = next(c); err != nil {
				panic(err)
			}

			if err = trx.Commit(); err != nil {
				panic(err)
			}

			c.Logger().Infof("%s finished in %s", path, log.Magenta(time.Since(start).String()))
			c.Logger().Debugf("Transaction committed")
			return err
		}
	}
}
