package cmd

import (
	"context"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"

	_ "github.com/getfider/fider/app/services/blob/fs"
	_ "github.com/getfider/fider/app/services/blob/s3"
	_ "github.com/getfider/fider/app/services/blob/sql"
	_ "github.com/getfider/fider/app/services/email/mailgun"
	_ "github.com/getfider/fider/app/services/email/smtp"
	_ "github.com/getfider/fider/app/services/httpclient"
	_ "github.com/getfider/fider/app/services/log/console"
	_ "github.com/getfider/fider/app/services/log/file"
	_ "github.com/getfider/fider/app/services/log/sql"
	_ "github.com/getfider/fider/app/services/oauth"
	_ "github.com/getfider/fider/app/services/sqlstore/postgres"
	_ "github.com/getfider/fider/app/services/webhook"
)

//RunServer starts the Fider Server
//Returns an exitcode, 0 for OK and 1 for ERROR
func RunServer() int {
	svcs := bus.Init()
	ctx := log.WithProperty(context.Background(), log.PropertyKeyTag, "BOOTSTRAP")
	for _, s := range svcs {
		log.Debugf(ctx, "Service '@{ServiceCategory}.@{ServiceName}' has been initialized.", dto.Props{
			"ServiceCategory": s.Category(),
			"ServiceName":     s.Name(),
		})
	}

	bus.Publish(ctx, &cmd.PurgeExpiredNotifications{})
	copyEtcFiles(ctx)

	e := routes(web.New())

	go e.Start(":" + env.Config.Port)
	return listenSignals(e)
}

// on startup, copy all etc/ files from configured blob storage into local etc/ folder
// this can be used to avoid having to mount volumes on ephemeral environments
func copyEtcFiles(ctx context.Context) {
	q := &query.ListBlobs{Prefix: "etc/"}
	if err := bus.Dispatch(ctx, q); err != nil {
		panic(errors.Wrap(err, "failed to list etc/ blobs"))
	}

	if len(q.Result) == 0 {
		log.Debug(ctx, "No etc/ files to copy")
		return
	}

	for _, blobKey := range q.Result {
		getBlob := &query.GetBlobByKey{Key: blobKey}
		if err := bus.Dispatch(ctx, getBlob); err != nil {
			panic(errors.Wrap(err, "failed to get blob by key '%s'", blobKey))
		}
		if err := os.MkdirAll(path.Dir(blobKey), 0774); err != nil {
			panic(errors.Wrap(err, "failed to create dir"))
		}
		if err := os.WriteFile(blobKey, getBlob.Result.Content, 0774); err != nil {
			panic(errors.Wrap(err, "failed to write blob '%s'", blobKey))
		}

		log.Debugf(ctx, "Copied '@{BlobKey}' to etc/ folder.", dto.Props{
			"BlobKey": blobKey,
		})
	}
}

func listenSignals(e *web.Engine) int {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, append([]os.Signal{syscall.SIGTERM, syscall.SIGINT}, extraSignals...)...)
	for {
		s := <-signals
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			err := e.Stop()
			if err != nil {
				return 1
			}
			return 0
		default:
			ret := handleExtraSignal(s, e)
			if ret >= 0 {
				return ret
			}
		}
	}
}
