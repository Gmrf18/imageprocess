# Optimizador de Imágenes WebP

Herramienta de línea de comandos escrita en Go que procesa imágenes en lote: las redimensiona a múltiples anchos y las convierte a **WebP** o **JPEG**, manteniendo la proporción original.

## Qué hace

- Lee todas las imágenes de un directorio de entrada
- Detecta orientación (portrait / landscape)
- Redimensiona a uno o varios anchos predefinidos usando el filtro Lanczos
- Exporta en WebP o JPEG con el nivel de calidad elegido
- Soporta formatos de entrada: `.jpg`, `.jpeg`, `.png`, `.webp`, `.cr3`
- Para archivos `.cr3` (RAW Canon) extrae el JPEG embebido vía `exiftool`
- Nombra los archivos de salida como `nombre_orientacion_ancho.ext`

## Flujo de uso

Al ejecutar el programa, responde las siguientes preguntas interactivas:

### 1. Directorios

```
📁 Directorio de entrada [default: input]:
📁 Directorio de salida [default: output]:
```

Presiona Enter para usar los valores por defecto (`input/` y `output/`).

### 2. Tamaños de salida

```
1) 320px  (Móvil S)
2) 468px  (Móvil L)
3) 768px  (Tablet)
4) 1080px (Full HD)
5) 1440px (2K)
6) 2160px (4K)
7) Personalizado (tamaño en px)
```

Puedes elegir uno o varios separados por coma: `1,3,5`  
Si eliges `7`, introduce los anchos manualmente: `800, 1200`

> Si el ancho solicitado supera el de la imagen original, se limita al original para evitar pixelado.

### 3. Calidad

```
1) Sin pérdida / Pérdida mínima (Lossless)
2) Calidad Máxima  (95%)
3) Calidad Alta    (80%) — Predeterminado
4) Calidad Media-Alta (65%)
5) Calidad Media   (45%)
6) Calidad Baja    (25%)
```

### 4. Formato de salida

```
1) WebP  — Predeterminado
2) JPEG
```

> JPEG no admite modo lossless; si se combinan ambas opciones, la calidad se ajusta a 100.

### 5. Procesado

El programa recorre el directorio de entrada y por cada imagen genera los archivos solicitados en el directorio de salida.

Ejemplo de archivos generados:

```
foto_landscape_768.webp
foto_landscape_1080.webp
foto_landscape_1440.webp
```

---

## Compilación

### Linux / WSL

```bash
go build -o imageprocess main.go
```

### Windows (.exe) desde Linux / WSL

Requiere el compilador cruzado `mingw-w64`:

```bash
CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o imageprocess.exe main.go
```

---

## Dependencias

| Paquete | Uso |
|---|---|
| `github.com/chai2010/webp` | Codificación y decodificación WebP (CGO) |
| `github.com/disintegration/imaging` | Redimensionado y apertura de imágenes |
| `exiftool` *(externo)* | Extracción de JPEG embebido en archivos CR3 |
