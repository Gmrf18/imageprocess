package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

type TargetSize struct {
	Name  string
	Width int
	IsPct bool
}

const (
	ColorReset  = "\033[0m"
	ColorCyan   = "\033[36m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorBold   = "\033[1m"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(ColorCyan + ColorBold + "🚀 Optimizador de Imágenes WebP" + ColorReset)
	fmt.Println("-------------------------------")

	// 1. Preguntar directorios
	fmt.Print(ColorYellow + "📁 Directorio de entrada [default: input]: " + ColorReset)
	inDir, _ := reader.ReadString('\n')
	inDir = strings.TrimSpace(inDir)
	if inDir == "" {
		inDir = "input"
	}

	fmt.Print(ColorYellow + "📁 Directorio de salida [default: output]: " + ColorReset)
	outDir, _ := reader.ReadString('\n')
	outDir = strings.TrimSpace(outDir)
	if outDir == "" {
		outDir = "output"
	}

	// 2. Mostrar menú de tamaños
	fmt.Println(ColorCyan + ColorBold + "\n📐 Selecciona los tamaños deseados (separados por coma, ej: 1,2,3):" + ColorReset)
	options := []TargetSize{
		{Name: "320px (Móvil S)", Width: 320, IsPct: false},
		{Name: "468px (Móvil L)", Width: 468, IsPct: false},
		{Name: "768px (Tablet)", Width: 768, IsPct: false},
		{Name: "1080px (Full HD)", Width: 1080, IsPct: false},
		{Name: "1440px (2K)", Width: 1440, IsPct: false},
		{Name: "2160px (4K)", Width: 2160, IsPct: false},
	}

	for i, opt := range options {
		fmt.Printf(ColorCyan+"%d) %s\n"+ColorReset, i+1, opt.Name)
	}
	fmt.Println(ColorCyan + "7) Personalizado (tamaño en px)" + ColorReset)

	fmt.Print(ColorYellow + "\nOpción(es): " + ColorReset)
	choicesStr, _ := reader.ReadString('\n')
	choicesStr = strings.TrimSpace(choicesStr)

	selectedChoices := strings.Split(choicesStr, ",")
	var targets []TargetSize

	for _, c := range selectedChoices {
		idx, err := strconv.Atoi(strings.TrimSpace(c))
		if err != nil || idx < 1 || idx > 7 {
			continue
		}
		if idx == 7 {
			fmt.Print(ColorYellow + "Introduce los tamaños en px (ej: 800, 1200): " + ColorReset)
			pxStr, _ := reader.ReadString('\n')
			pxList := strings.Split(strings.TrimSpace(pxStr), ",")
			for _, p := range pxList {
				width, _ := strconv.Atoi(strings.TrimSpace(p))
				if width > 0 {
					targets = append(targets, TargetSize{Name: fmt.Sprintf("%dpx", width), Width: width, IsPct: false})
				}
			}
		} else {
			targets = append(targets, options[idx-1])
		}
	}

	if len(targets) == 0 {
		log.Fatal("No seleccionaste ningún tamaño válido.")
	}

	// 3. Menú de calidad
	fmt.Println(ColorCyan + ColorBold + "\n💎 Selecciona la calidad de salida:" + ColorReset)
	fmt.Println(ColorCyan + "1) Sin pérdida / Pérdida mínima (Lossless)" + ColorReset)
	fmt.Println(ColorCyan + "2) Calidad Máxima (95%)" + ColorReset)
	fmt.Println(ColorCyan + "3) Calidad Alta (80%) [Predeterminado]" + ColorReset)
	fmt.Println(ColorCyan + "4) Calidad Media-Alta (65%)" + ColorReset)
	fmt.Println(ColorCyan + "5) Calidad Media (45%)" + ColorReset)
	fmt.Println(ColorCyan + "6) Calidad Baja / Menor peso (25%)" + ColorReset)

	fmt.Print(ColorYellow + "\nOpción (1-6): " + ColorReset)
	qStr, _ := reader.ReadString('\n')
	qStr = strings.TrimSpace(qStr)

	qualityLevel, _ := strconv.Atoi(qStr)
	if qualityLevel < 1 || qualityLevel > 6 {
		qualityLevel = 3 // Default
	}

	var lossless bool
	var quality float32

	switch qualityLevel {
	case 1:
		lossless = true
		quality = 100
	case 2:
		lossless = false
		quality = 95
	case 3:
		lossless = false
		quality = 80
	case 4:
		lossless = false
		quality = 65
	case 5:
		lossless = false
		quality = 45
	case 6:
		lossless = false
		quality = 25
	default:
		lossless = false
		quality = 80
	}

	// 4. Menú de formato de salida
	fmt.Println(ColorCyan + ColorBold + "\n🗂️  Selecciona el formato de salida:" + ColorReset)
	fmt.Println(ColorCyan + "1) WebP [Predeterminado]" + ColorReset)
	fmt.Println(ColorCyan + "2) JPEG" + ColorReset)

	fmt.Print(ColorYellow + "\nOpción (1-2): " + ColorReset)
	fStr, _ := reader.ReadString('\n')
	fStr = strings.TrimSpace(fStr)

	outputFormat := "webp"
	if fStr == "2" {
		outputFormat = "jpeg"
	}

	if outputFormat == "jpeg" && lossless {
		fmt.Println(ColorYellow + "ℹ️  JPEG no admite modo lossless; se usará calidad 100." + ColorReset)
		lossless = false
		quality = 100
	}

	// 5. Procesar imágenes
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("Error al crear carpeta de salida: %v", err)
	}

	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatalf("Error al leer entrada: %v", err)
	}

	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".cr3": true}

	for _, entry := range files {
		if entry.IsDir() || !validExts[strings.ToLower(filepath.Ext(entry.Name()))] {
			continue
		}

		inputPath := filepath.Join(inDir, entry.Name())
		var src image.Image
		if strings.ToLower(filepath.Ext(entry.Name())) == ".cr3" {
			src, err = decodeCR3(inputPath)
		} else {
			src, err = imaging.Open(inputPath)
		}
		if err != nil {
			fmt.Printf(ColorRed+"❌ Error abriendo %s: %v\n"+ColorReset, entry.Name(), err)
			continue
		}

		bounds := src.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		orientation := "portrait"
		if w >= h {
			orientation = "landscape"
		}

		fmt.Printf(ColorCyan+"\n🖼️  Procesando: %s [%s]\n"+ColorReset, entry.Name(), orientation)
		baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

		for _, target := range targets {
			var finalWidth int
			var suffix string

			if target.IsPct {
				finalWidth = int(float64(w) * float64(target.Width) / 100.0)
				suffix = fmt.Sprintf("%dpct", target.Width)
			} else {
				finalWidth = target.Width
				suffix = strconv.Itoa(target.Width)
			}

			if finalWidth > w {
				fmt.Printf(ColorYellow+"  ⚠️ El tamaño solicitado (%dpx) es mayor al original. Limitando a %dpx para evitar pixelado.\n"+ColorReset, finalWidth, w)
				finalWidth = w
			}

			if finalWidth <= 0 {
				finalWidth = 1
			}

			// Redimensionar manteniendo ratio
			dst := imaging.Resize(src, finalWidth, 0, imaging.Lanczos)

			ext := "webp"
			if outputFormat == "jpeg" {
				ext = "jpg"
			}
			outName := fmt.Sprintf("%s_%s_%s.%s", baseName, orientation, suffix, ext)
			outPath := filepath.Join(outDir, outName)

			err := saveImage(outPath, dst, outputFormat, lossless, quality)
			if err != nil {
				fmt.Printf(ColorRed+"  ⚠️ Error guardando %s: %v\n"+ColorReset, outName, err)
			} else {
				fmt.Printf(ColorGreen+"  ✅ Generado: %s\n"+ColorReset, outName)
			}
		}
	}
	fmt.Println(ColorGreen + ColorBold + "\n✨ ¡Proceso finalizado!" + ColorReset)
}

func saveImage(path string, img image.Image, format string, lossless bool, quality float32) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: int(quality)})
	default:
		return webp.Encode(f, img, &webp.Options{Lossless: lossless, Quality: quality})
	}
}

var exiftoolPath string
var exiftoolChecked bool

func decodeCR3(path string) (image.Image, error) {
	if !exiftoolChecked {
		p, err := exec.LookPath("exiftool")
		if err == nil {
			exiftoolPath = p
		}
		exiftoolChecked = true
	}
	if exiftoolPath == "" {
		return nil, fmt.Errorf("exiftool no encontrado en PATH (necesario para .cr3)")
	}

	for _, tag := range []string{"-JpgFromRaw", "-PreviewImage"} {
		var buf bytes.Buffer
		cmd := exec.Command(exiftoolPath, "-b", tag, path)
		cmd.Stdout = &buf
		if err := cmd.Run(); err != nil {
			continue
		}
		if buf.Len() == 0 {
			continue
		}
		img, err := imaging.Decode(&buf)
		if err == nil {
			return img, nil
		}
	}
	return nil, fmt.Errorf("no se pudo extraer JPEG embebido del CR3")
}
