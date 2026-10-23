# gosnowflake log-leak demo (H1 Snowflake PoC)

CI demo: the Go Snowflake driver writes the DSN query string (incl. passcode) to stderr at the documented INFO logging level. GitHub Actions captures stderr into this run's **public** log — readable by anyone on the internet without any account.
