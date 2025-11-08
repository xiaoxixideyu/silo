module silo

go 1.24.3

require (
	dario.cat/mergo v1.0.2
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/ekrucio/apm-agent-go/module/apmgoredisv9/v2 v2.0.0-20240214105830-6f3c22353857
	github.com/go-co-op/gocron/v2 v2.17.0
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.28.0
	github.com/go-redsync/redsync/v4 v4.14.0
	github.com/go-resty/resty/v2 v2.16.5
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.7.0
	github.com/jackc/pgconn v1.14.3
	github.com/jmattheis/goverter v1.9.2
	github.com/labstack/echo/v4 v4.13.4
	github.com/labstack/gommon v0.4.2
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.14.1
	github.com/rs/zerolog v1.34.0
	github.com/samber/lo v1.52.0
	github.com/samber/slog-common v0.19.0
	github.com/spf13/viper v1.15.0
	github.com/swaggo/echo-swagger v1.4.1
	github.com/swaggo/swag v1.16.6
	github.com/topfreegames/pitaya/v2 v2.11.18
	go.elastic.co/apm/module/apmechov4/v2 v2.7.1
	go.elastic.co/apm/module/apmgormv2/v2 v2.7.1
	go.elastic.co/apm/v2 v2.7.1
	golang.org/x/exp v0.0.0-20251017212417-90e834f514db
	golang.org/x/sync v0.17.0
	golang.org/x/time v0.11.0
	golang.org/x/tools v0.38.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dave/jennifer v1.6.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/go-sysinfo v1.7.1 // indirect
	github.com/elastic/go-windows v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.10 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/google/subcommands v1.2.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx/v4 v4.18.2 // indirect
	github.com/jackc/pgx/v5 v5.5.4 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/swaggo/files/v2 v2.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.elastic.co/apm/module/apmhttp/v2 v2.7.1 // indirect
	go.elastic.co/apm/module/apmsql/v2 v2.7.1 // indirect
	go.elastic.co/fastjson v1.5.1 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.5.4 // indirect
	howett.net/plist v1.0.0 // indirect
)
