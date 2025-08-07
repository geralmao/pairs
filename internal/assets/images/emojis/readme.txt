*** NO SE AÑADEN LOS SPRITES DE EMOJIS PARA REDUCIR ESPACIO ***
Este directorio contiene los sprites de los emojis descargados de https://openmoji.org.

Para evitar imágenes repetidas pero con composición de nombres diferentes puede ejecutar 
el programa go que se incluye al final de este texto, que renombra los archivos con
contenido identico pero con nombre diferente colocando el prefijo "XXXX_" con lo que se 
puede eliminar estos ficheros si lo desea.



*** EMOJI SPRITES ARE NOT INCLUDED IN THIS REPOSITORY TO REDUCE REPO SIZE ***
This directory is intended to hold the emoji sprite images downloaded from https://openmoji.org.

To avoid storing duplicate images with different filenames, you can use the provided Go script
(at the end of this file). It scans the directory for files with identical content and renames
duplicates by prefixing them with "XXXX_". You can then delete those "XXXX_" files if you want
to reduce redundancy.



~~~go
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Printf("Buscando y calculando MD5 de imágenes en: %s\n", dirPath)

	var imageExtensions = map[string]bool{
		".png": true,
	}
	md5Sums := make(map[string][]string)

	// Recorrer el directorio para calcular el MD5 de cada imagen
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if imageExtensions[ext] {
				fmt.Printf("  Calculando MD5 para: %s\n", filepath.Base(path))
				hash, err := calculateMD5(path)
				if err != nil {
					fmt.Printf("    Error al calcular MD5 para '%s': %v\n", filepath.Base(path), err)
					return nil // Continuar con otros archivos
				}
				md5Sums[hash] = append(md5Sums[hash], path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error al escanear directorio o calcular MD5: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAnálisis de MD5 completado. Identificando duplicados...")

	duplicatesFound := 0
	// Iterar sobre el mapa de sumas MD5 para encontrar duplicados
	for md5Hash, files := range md5Sums {
		if len(files) > 1 {
			duplicatesFound++
			fmt.Printf("\n--- Duplicados encontrados para MD5: %s ---\n", md5Hash)
			fmt.Printf("  Original: '%s'\n", filepath.Base(files[0])) // El primer archivo es el "original"

			// Renombrar los archivos duplicados (a partir del segundo)
			for i := 1; i < len(files); i++ {
				duplicateFile := files[i]
				fmt.Printf("  Duplicado: '%s'\n", filepath.Base(duplicateFile))

				err := renameFileWithPrefix(duplicateFile, "XXXX_")
				if err != nil {
					fmt.Printf("    ❌ Error al renombrar '%s': %v\n", filepath.Base(duplicateFile), err)
				} else {
					fmt.Printf("    ✔️ Renombrado '%s' a 'XXXX_%s'\n", filepath.Base(duplicateFile), filepath.Base(duplicateFile))
				}
			}
		}
	}

	if duplicatesFound == 0 {
		fmt.Println("\nNo se encontraron imágenes duplicadas.")
	} else {
		fmt.Printf("\nProceso completado. Se encontraron y marcaron %d grupos de imágenes duplicadas.\n", duplicatesFound)
	}
}

// calculateMD5 calcula la suma MD5 de un archivo.
func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el archivo '%s': %w", filePath, err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("error al leer el archivo '%s' para MD5: %w", filePath, err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// renameFileWithPrefix renombra un archivo añadiendo un prefijo.
func renameFileWithPrefix(filePath, prefix string) error {
	dir := filepath.Dir(filePath)
	oldName := filepath.Base(filePath)
	newName := prefix + oldName
	newPath := filepath.Join(dir, newName)

	return os.Rename(filePath, newPath)
}

const dirPath string = "/code/pairs/internal/assets/images/emojis"
~~~
