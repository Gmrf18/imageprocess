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

	"github.com/AlecAivazis/survey/v2"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

type TargetSize struct {
	Name  string
	Width int
	IsPct bool
}

var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fd := os.Stdin.Fd()
	interactive := isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)

	fmt.Println(bold(cyan("🚀 Optimizador de Imágenes WebP")))
	fmt.Println("-------------------------------")

	// 1. Preguntar directorios
	var inDir string
	for {
		fmt.Print(yellow("📁 Directorio de entrada [default: input]: "))
		inDir, _ = reader.ReadString('\n')
		inDir = strings.TrimSpace(inDir)
		if inDir == "" {
			inDir = "input"
		}

		info, err := os.Stat(inDir)
		if err != nil || !info.IsDir() {
			fmt.Printf("%s La carpeta '%s' no existe o no es un directorio.\n", red("❌"), inDir)
			continue
		}
		break
	}

	fmt.Print(yellow("📁 Directorio de salida [default: output]: "))
	outDir, _ := reader.ReadString('\n')
	outDir = strings.TrimSpace(outDir)
	if outDir == "" {
		outDir = "output"
	}

	options := []TargetSize{
		{Name: "240px (Watch/Legacy)", Width: 240, IsPct: false},
		{Name: "280px (Galaxy Fold cerrado)", Width: 280, IsPct: false},
		{Name: "320px (Móvil S)", Width: 320, IsPct: false},
		{Name: "468px (Móvil L)", Width: 468, IsPct: false},
		{Name: "768px (Tablet)", Width: 768, IsPct: false},
		{Name: "1080px (Full HD)", Width: 1080, IsPct: false},
		{Name: "1440px (2K)", Width: 1440, IsPct: false},
		{Name: "2160px (4K)", Width: 2160, IsPct: false},
	}

	// 2. Selección de tamaños, calidad y formato
	targets := selectSizes(reader, options, interactive)
	lossless, quality := selectQuality(reader, interactive)
	outputFormat := selectFormat(reader, interactive)

	if outputFormat == "jpeg" && lossless {
		fmt.Println(yellow("ℹ️  JPEG no admite modo lossless; se usará calidad 100."))
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
			fmt.Printf("%s Error abriendo %s: %v\n", red("❌"), entry.Name(), err)
			continue
		}

		bounds := src.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		orientation := "portrait"
		if w >= h {
			orientation = "landscape"
		}

		fmt.Printf("\n%s Procesando: %s [%s]\n", cyan("🖼️ "), entry.Name(), orientation)
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
				fmt.Printf("%s El tamaño solicitado (%dpx) es mayor al original. Limitando a %dpx para evitar pixelado.\n", yellow("  ⚠️ "), finalWidth, w)
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
				fmt.Printf("%s Error guardando %s: %v\n", red("  ⚠️ "), outName, err)
			} else {
				fmt.Printf("%s Generado: %s\n", green("  ✅ "), outName)
			}
		}
	}
	fmt.Println(bold(green("\n✨ ¡Proceso finalizado!")))
}

// selectSizes pide los tamaños de salida. En modo interactivo usa checkboxes
// (multi-selección); si no hay terminal interactiva, cae al menú numérico.
func selectSizes(reader *bufio.Reader, options []TargetSize, interactive bool) []TargetSize {
	var targets []TargetSize

	if interactive {
		labels := make([]string, 0, len(options)+1)
		for _, o := range options {
			labels = append(labels, o.Name)
		}
		customIdx := len(labels)
		labels = append(labels, "Personalizado (tamaño en px)")

		var picked []int
		prompt := &survey.MultiSelect{
			Message: "Selecciona los tamaños deseados:",
			Options: labels,
		}
		if err := survey.AskOne(prompt, &picked); err != nil {
			log.Fatalf("Selección cancelada: %v", err)
		}

		customSelected := false
		for _, idx := range picked {
			if idx == customIdx {
				customSelected = true
				continue
			}
			targets = append(targets, options[idx])
		}

		if customSelected {
			var pxStr string
			if err := survey.AskOne(&survey.Input{
				Message: "Introduce los tamaños en px (ej: 800, 1200):",
			}, &pxStr); err != nil {
				log.Fatalf("Selección cancelada: %v", err)
			}
			targets = append(targets, parseCustomSizes(pxStr)...)
		}
	} else {
		fmt.Println(bold(cyan("\n📐 Selecciona los tamaños deseados (separados por coma, ej: 1,2,3):")))
		for i, opt := range options {
			fmt.Printf("%s%d) %s\n", cyan(""), i+1, opt.Name)
		}
		customNum := len(options) + 1
		fmt.Println(cyan(fmt.Sprintf("%d) Personalizado (tamaño en px)", customNum)))

		fmt.Print(yellow("\nOpción(es): "))
		choicesStr, _ := reader.ReadString('\n')
		choicesStr = strings.TrimSpace(choicesStr)

		for _, c := range strings.Split(choicesStr, ",") {
			idx, err := strconv.Atoi(strings.TrimSpace(c))
			if err != nil || idx < 1 || idx > customNum {
				continue
			}
			if idx == customNum {
				fmt.Print(yellow("Introduce los tamaños en px (ej: 800, 1200): "))
				pxStr, _ := reader.ReadString('\n')
				targets = append(targets, parseCustomSizes(pxStr)...)
			} else {
				targets = append(targets, options[idx-1])
			}
		}
	}

	if len(targets) == 0 {
		log.Fatal("No seleccionaste ningún tamaño válido.")
	}
	return targets
}

func parseCustomSizes(pxStr string) []TargetSize {
	var out []TargetSize
	for _, p := range strings.Split(strings.TrimSpace(pxStr), ",") {
		width, _ := strconv.Atoi(strings.TrimSpace(p))
		if width > 0 {
			out = append(out, TargetSize{Name: fmt.Sprintf("%dpx", width), Width: width, IsPct: false})
		}
	}
	return out
}

// selectQuality devuelve (lossless, quality) según la opción elegida.
func selectQuality(reader *bufio.Reader, interactive bool) (bool, float32) {
	qualityLabels := []string{
		"Sin pérdida / Pérdida mínima (Lossless)",
		"Calidad Máxima (95%)",
		"Calidad Alta (80%) [Predeterminado]",
		"Calidad Media-Alta (65%)",
		"Calidad Media (45%)",
		"Calidad Baja / Menor peso (25%)",
	}

	qualityLevel := 3
	if interactive {
		idx := 2
		prompt := &survey.Select{
			Message: "Selecciona la calidad de salida:",
			Options: qualityLabels,
			Default: qualityLabels[2],
		}
		if err := survey.AskOne(prompt, &idx); err != nil {
			log.Fatalf("Selección cancelada: %v", err)
		}
		qualityLevel = idx + 1
	} else {
		fmt.Println(bold(cyan("\n💎 Selecciona la calidad de salida:")))
		for i, l := range qualityLabels {
			fmt.Println(cyan(fmt.Sprintf("%d) %s", i+1, l)))
		}
		fmt.Print(yellow("\nOpción (1-6): "))
		qStr, _ := reader.ReadString('\n')
		qStr = strings.TrimSpace(qStr)
		qualityLevel, _ = strconv.Atoi(qStr)
		if qualityLevel < 1 || qualityLevel > 6 {
			qualityLevel = 3
		}
	}

	switch qualityLevel {
	case 1:
		return true, 100
	case 2:
		return false, 95
	case 3:
		return false, 80
	case 4:
		return false, 65
	case 5:
		return false, 45
	case 6:
		return false, 25
	default:
		return false, 80
	}
}

// selectFormat devuelve "webp" o "jpeg".
func selectFormat(reader *bufio.Reader, interactive bool) string {
	if interactive {
		idx := 0
		prompt := &survey.Select{
			Message: "Selecciona el formato de salida:",
			Options: []string{"WebP", "JPEG"},
			Default: "WebP",
		}
		if err := survey.AskOne(prompt, &idx); err != nil {
			log.Fatalf("Selección cancelada: %v", err)
		}
		if idx == 1 {
			return "jpeg"
		}
		return "webp"
	}

	fmt.Println(bold(cyan("\n🗂️  Selecciona el formato de salida:")))
	fmt.Println(cyan("1) WebP [Predeterminado]"))
	fmt.Println(cyan("2) JPEG"))
	fmt.Print(yellow("\nOpción (1-2): "))
	fStr, _ := reader.ReadString('\n')
	fStr = strings.TrimSpace(fStr)
	if fStr == "2" {
		return "jpeg"
	}
	return "webp"
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
