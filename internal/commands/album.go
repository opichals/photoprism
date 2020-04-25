package commands

import (
	"context"
	"fmt"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/query"

	"github.com/urfave/cli"
)

var AlbumCommand = cli.Command{
	Name:    "album",
	Aliases: []string{""},
	Usage:   "Adds photos to an album",
	Flags:   albumFlags,
	Action:  albumAction,
}

var albumFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "album",
		Value: "",
		Usage: "Album to add photos to",
	},
}

// albumAction adds photos to an album.
func albumAction(ctx *cli.Context) error {
	conf := config.NewConfig(ctx)

	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := conf.Init(cctx); err != nil {
		return err
	}

	conf.MigrateDb()

	db := conf.Db()
	q := query.New(db)
	result, err := q.Albums(form.AlbumSearch{Name: ctx.String("album")})

	if err != nil {
		log.Errorf("album: %s", err)
		return nil
	}
	a := result[0]

	log.Debugf("Album %+v\n", a.AlbumName)

	photos, _, err := q.PhotosByFilenames(ctx.Args())
	log.Debugf("Photo %+v\n", photos[0].PhotoTitle)

	var added []*entity.PhotoAlbum

	for _, p := range photos {
		added = append(added, entity.NewPhotoAlbum(p.PhotoUUID, a.AlbumUUID).FirstOrCreate(db))
	}

	if len(added) == 1 {
		log.Info(fmt.Sprintf("one photo added to %s", a.AlbumName))
	} else {
		log.Info(fmt.Sprintf("%d photos added to %s", len(added), a.AlbumName))
	}
	return nil
}
