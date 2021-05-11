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
	Use:  "inituser dbPath",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(args[0])
	},
}

func main() {
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database; %s; %w", dbPath, err)
	}
	users := []model.User{
		{
			Account:  "taro@yamada.jp",
			Name:     "taro",
			Password: "yamada",
			Message:  "taro's account",
		},
		{
			Account:  "hanako@flower.com",
			Name:     "hanako",
			Password: "flower",
			Message:  "hanako's account",
		},
	}
	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("failed to create user; %w", err)
	}
	return nil
}
