package migrations

import (
	"log"
	"os/exec"
)

func RunMigrations(dsn string) error {
	cmd := exec.Command("goose", "postgres", dsn, "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Migration output: %s", output)
		return err
	}
	log.Printf("Migrations applied successfully")
	return nil
}
