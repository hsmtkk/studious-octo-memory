package main

import (
	"fmt"
	"log"

	"github.com/hsmtkk/studious-octo-memory/model"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var command = &cobra.Command{
	Use: "migrate dbPath",
	Run: func(cmd *cobra.Command, args []string) {
		migrate(args[0])
	},
}

func main() {
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}

func migrate(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database; %s; %w", dbPath, err)
	}
	db.AutoMigrate(&model.User{}, &model.Group{}, &model.Post{}, &model.Comment{})
	return nil
}
