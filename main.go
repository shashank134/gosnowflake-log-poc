package main

import (
	gs "github.com/snowflakedb/gosnowflake"
)

func main() {
	// "Victim": an ordinary Go service enabling driver logging as documented
	lg := gs.GetLogger()
	_ = lg.SetLogLevel("info")
	// Ordinary connection setup — parse the DSN exactly as real apps do
	_, _ = gs.ParseDSN("svc-user:H1-VICTIM-PASSWORD@acme-prod.snowflakecomputing.com/db?passcode=H1VICTIM-MFA-987654")
}
