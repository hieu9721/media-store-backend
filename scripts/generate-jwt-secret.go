package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	// Generate 32 bytes (256 bits)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating secret: %v\n", err)
		os.Exit(1)
	}

	secret := hex.EncodeToString(bytes)

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║           JWT SECRET GENERATOR                                 ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Generated JWT Secret (256-bit):")
	fmt.Println()
	fmt.Printf("JWT_SECRET=%s\n", secret)
	fmt.Println()
	fmt.Println("Add this to your .env file:")
	fmt.Println("1. Copy the line above")
	fmt.Println("2. Replace the existing JWT_SECRET in .env")
	fmt.Println("3. Restart your server")
	fmt.Println()
	fmt.Println("⚠️  WARNING:")
	fmt.Println("- Keep this secret safe!")
	fmt.Println("- Never commit .env to version control")
	fmt.Println("- Use different secrets for dev/staging/production")
	fmt.Println()
}
