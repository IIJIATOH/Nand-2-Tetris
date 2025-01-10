package main

import (
	"assembler/internal/mypackage"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var pathToFile string

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "MyApp is a CLI application",
	Long:  `MyApp is a CLI application for translating assembly programs into Hack binary code.`,
}

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate assembly program into Hack binary code",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Working, like i think. Word is :%s", pathToFile)
		mypackage.Translate(pathToFile)
	},
}

func main() {
	// Определение флагов
	translateCmd.Flags().StringVarP(&pathToFile, "path", "p", "C:/Users/mikhailovpa.DESKTOP-OKO95JV/Downloads/nand/nand2tetris/projects/6/max/Max.asm", "A path to file to translate")

	// Добавление подкоманды к корневой команде
	rootCmd.AddCommand(translateCmd)

	// Выполнение корневой команды
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
