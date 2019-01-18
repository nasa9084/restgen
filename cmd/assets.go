package main

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetsfbaa88640bf76efe7dacb8b0d8ddfc156a2acb1b = "package main\n\nimport (\n\t\"context\"\n\t\"log\"\n\t\"net/http\"\n\t\"os\"\n\t\"os/signal\"\n\t\"syscall\"\n\t\"time\"\n\n\tflags \"github.com/jessevdk/go-flags\"\n\t\"github.com/pkg/errors\"\n)\n\ntype options struct {\n\tAddr string `short:\"a\" long:\"address\" default:\":8080\" description:\"address to listen to\"`\n}\n\nfunc main() {\n\tif err := _main(); err != nil {\n\t\tlog.Printf(\"error: %v\", err)\n\t\tos.Exit(1)\n\t}\n\tos.Exit(0)\n}\n\nfunc _main() error {\n\tvar opts options\n\tif _, err := flags.Parse(&opts); err != nil {\n\t\treturn errors.Wrap(err, \"parsing flags\")\n\t}\n\n\tctx, cancel := context.WithCancel(context.Background())\n\tdefer cancel()\n\n\tsrv := http.Server{\n\t\tAddr:    opts.Addr,\n\t\tHandler: NewRouter(),\n\t}\n\n\tsigCh := make(chan os.Signal, 1)\n\tsignal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)\n\tgo func(ctx context.Context) {\n\t\tfor {\n\t\t\tselect {\n\t\t\tcase <-ctx.Done():\n\t\t\t\treturn\n\t\t\tcase <-sigCh:\n\t\t\t\tctx, cancel := context.WithTimeout(ctx, 1*time.Minute)\n\t\t\t\tdefer cancel()\n\t\t\t\tif err := srv.Shutdown(ctx); err != nil {\n\t\t\t\t\tlog.Print(err)\n\t\t\t\t}\n\t\t\t\treturn\n\t\t\t}\n\t\t}\n\t}(ctx)\n\n\tif err := srv.ListenAndServe(); err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}\n"
var _Assetsf506750bdd1f5f056cd3153f6e7b0355f4ecdb21 = "---\nopenapi: 3.0.2\ninfo:\n  title: %s\n  version: v0.0.1\npaths:\n  /:\n    get:\n      description: health check\n      operationId: HealthCheck\n      responses:\n        '200':\n          $ref: '#/components/responses/HealthCheckResponse'\ncomponents:\n  responses:\n    HealthCheckResponse:\n      description: response for HealthCheck\n      content:\n        application/json:\n          schema:\n            type: object\n            properties:\n              status:\n                type: string\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"assets"}, "/assets": []string{"server_main.go.tmpl", "default_spec.yaml.tmpl"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1547722785, 1547722785565408846),
		Data:     nil,
	}, "/assets": &assets.File{
		Path:     "/assets",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1547706804, 1547706804047334094),
		Data:     nil,
	}, "/assets/server_main.go.tmpl": &assets.File{
		Path:     "/assets/server_main.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1547776847, 1547776847384358551),
		Data:     []byte(_Assetsfbaa88640bf76efe7dacb8b0d8ddfc156a2acb1b),
	}, "/assets/default_spec.yaml.tmpl": &assets.File{
		Path:     "/assets/default_spec.yaml.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1547704898, 1547704898394350439),
		Data:     []byte(_Assetsf506750bdd1f5f056cd3153f6e7b0355f4ecdb21),
	}}, "")
