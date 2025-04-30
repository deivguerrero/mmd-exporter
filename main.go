// main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watch := flag.Bool("watch", false, "Observa cambios en archivos .mmd del directorio y los exporta automáticamente")
	flag.Parse()

	args := flag.Args()
	dirOrFile := "."

	if len(args) == 1 {
		dirOrFile = args[0]
	}

	info, err := os.Stat(dirOrFile)
	if os.IsNotExist(err) {
		fmt.Printf("El path '%s' no existe.\n", dirOrFile)
		os.Exit(1)
	}

	if info.IsDir() {
		if *watch {
			watchDirectory(dirOrFile)
		} else {
			processDirectory(dirOrFile)
		}
	} else if strings.HasSuffix(info.Name(), ".mmd") {
		processFile(dirOrFile)
	} else {
		fmt.Printf("Archivo no válido: %s\n", dirOrFile)
	}
}

func processDirectory(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("No se pudo leer el directorio: %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".mmd") {
			fullPath := filepath.Join(dir, entry.Name())
			processFile(fullPath)
		}
	}
}

func processFile(path string) {
	baseName := strings.TrimSuffix(path, ".mmd")

	fmt.Printf("Procesando %s...\n", path)

	// Generar PNG
	pngCmd := exec.Command("mmdc", "-i", path, "-o", baseName+".png", "--scale", "2")
	pngCmd.Stdout = os.Stdout
	pngCmd.Stderr = os.Stderr
	if err := pngCmd.Run(); err != nil {
		fmt.Printf("❌ Error generando PNG: %v\n", err)
	}

	// Generar SVG
	svgCmd := exec.Command("mmdc", "-i", path, "-o", baseName+".svg")
	svgCmd.Stdout = os.Stdout
	svgCmd.Stderr = os.Stderr
	if err := svgCmd.Run(); err != nil {
		fmt.Printf("❌ Error generando SVG: %v\n", err)
	}
}

func watchDirectory(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("❌ No se pudo iniciar el watcher: %v\n", err)
		os.Exit(1)
	}
	defer watcher.Close()

	if err := watcher.Add(dir); err != nil {
		fmt.Printf("❌ No se pudo observar el directorio: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🔁 Observando '%s' por cambios en archivos .mmd...\n", dir)

	debouncers := make(map[string]time.Time)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 &&
					strings.HasSuffix(event.Name, ".mmd") &&
					!isDebounced(event.Name, debouncers) {
					fmt.Printf("📝 Cambio detectado en %s\n", event.Name)
					processFile(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("❌ Error en watcher: %v\n", err)
			}
		}
	}()

	select {} // Mantiene el watcher vivo
}

func isDebounced(path string, debouncers map[string]time.Time) bool {
	now := time.Now()
	last, exists := debouncers[path]
	if exists && now.Sub(last) < 500*time.Millisecond {
		return true
	}
	debouncers[path] = now
	return false
}