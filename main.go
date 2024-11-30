package main

import (
	"log"
	"os"

	"github.com/dangquangdon/pgrest_rpc_client_generator/generator"
	"github.com/urfave/cli/v2"
)

var gen *generator.Generator

func main() {
	app := &cli.App{
		Name:  "pgr-gen",
		Usage: "Generate HTTP Client which calls the RPC functions defined in PostgREST server.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "destination",
				Value:   ".",
				Aliases: []string{"d"},
				Usage:   "Destination directory where the code is going to be stored",
			},
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "Base url of the PostgREST server",
			},
			&cli.StringFlag{
				Name:     "client-id",
				Aliases:  []string{"c"},
				Required: true,
				Usage:    "Value for the User-Agent Header",
			},
		},
		Before: func(ctx *cli.Context) error {
			dest := ctx.String("destination")
			url := ctx.String("url")
			clientId := ctx.String("client-id")
			gen = generator.NewGeneartor(url, dest, clientId)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate all the codes for types and requests",
				Action: func(ctx *cli.Context) error {
					err := gen.GenerateTypes()
					if err != nil {
						log.Fatal(err)
					}

					return gen.GenerateRequests()
				},
			},
			{
				Name:  "generate-types",
				Usage: "Generate only the types",
				Action: func(ctx *cli.Context) error {
					return gen.GenerateTypes()
				},
			},
			{
				Name:  "generate-requests",
				Usage: "Generate only the functions for http requests",
				Action: func(ctx *cli.Context) error {
					return gen.GenerateRequests()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
