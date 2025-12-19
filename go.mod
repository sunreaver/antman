module github.com/sunreaver/antman/v4

go 1.24.0

replace (
	github.com/godoes/gorm-oracle => github.com/ggicegg/gorm-oceanbase-oracle v0.0.4
	github.com/mattn/go-oci8 => github.com/ggicegg/go-oceanbase-oci8 v0.0.1
)

require (
	gitee.com/opengauss/openGauss-connector-go-pq v1.0.4
	github.com/IBM/sarama v1.43.2
	github.com/go-redis/redis/v8 v8.11.5
	github.com/godoes/gorm-dameng v0.1.1
	github.com/godoes/gorm-oracle v1.6.9
	github.com/mattn/go-oci8 v0.1.1
	github.com/pkg/errors v0.9.1
	github.com/sijms/go-ora/v2 v2.8.18
	github.com/sunreaver/logger/v3 v3.0.3
	github.com/sunreaver/tomlanalysis v1.0.0
	gorm.io/driver/mysql v1.5.6
	gorm.io/driver/postgres v1.5.7
	gorm.io/driver/sqlite v1.5.5
	gorm.io/gorm v1.25.10
	gorm.io/plugin/dbresolver v1.5.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/eapache/go-resiliency v1.6.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
