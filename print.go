package main

import "fmt"

func PrintResult(result Result) {
	fmt.Printf("%-8s : %s\n", "URL", result.URL)

	if result.Error == "" {
		fmt.Printf("%-8s : %s\n",       "Status",   result.Status)
		fmt.Printf("%-8s : %d bytes\n", "Size",     result.Size)
		fmt.Printf("%-8s : %d ms\n",    "Duration", result.Duration)
	} else {
		fmt.Printf("%-8s : %s\n", "Error", result.Error)
	}

	fmt.Println("--\n")
}
