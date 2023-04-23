package data

import (
	"crypto/tls"
	"database/sql"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// Connect abre a conexão com o banco de dados e a retorna
func Connect() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                     os.Getenv("DB_USER"),
		Passwd:                   os.Getenv("DB_SENHA"),
		Net:                      "tcp",
		Addr:                     "127.0.0.1:3306",
		DBName:                   "devbook",
		Params:                   map[string]string{},
		Collation:                "",
		Loc:                      &time.Location{},
		MaxAllowedPacket:         0,
		ServerPubKey:             "",
		TLSConfig:                "",
		TLS:                      &tls.Config{},
		Timeout:                  0,
		ReadTimeout:              0,
		WriteTimeout:             0,
		AllowAllFiles:            false,
		AllowCleartextPasswords:  false,
		AllowFallbackToPlaintext: false,
		AllowNativePasswords:     false,
		AllowOldPasswords:        false,
		CheckConnLiveness:        false,
		ClientFoundRows:          false,
		ColumnsWithAlias:         false,
		InterpolateParams:        false,
		MultiStatements:          false,
		ParseTime:                false,
		RejectReadOnly:           false,
	}
	db, erro := sql.Open("mysql", cfg.FormatDSN())
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}
