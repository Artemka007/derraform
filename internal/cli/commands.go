package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing...")
		// Инициализация state файла
	},
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show execution plan",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating plan...")
		// Парсинг конфига и план изменений
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Applying configuration...")
		// Применение изменений
	},
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy infrastructure",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Destroying infrastructure...")
		// Удаление ресурсов
	},
}
