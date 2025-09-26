var applyCmd = &cobra.Command{
    Use:   "apply",
    Short: "Apply configuration",
    Run: func(cmd *cobra.Command, args []string) {
        engine, err := core.NewEngine()
        if err != nil {
            fmt.Printf("Failed to initialize engine: %v\n", err)
            os.Exit(1)
        }

        if err := engine.Apply("main.tf"); err != nil {
            fmt.Printf("Apply failed: %v\n", err)
            os.Exit(1)
        }
    },
}