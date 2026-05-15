# 📸 Optimizador de Imágenes — Manual de Usuario

## 🎯 ¿Qué hace esta aplicación?

**Optimizador de Imágenes** es una herramienta que redimensiona y convierte tus fotos a formatos web optimizados. Procesa lotes completos de imágenes en un solo paso, generando versiones múltiples en diferentes tamaños (móvil, tablet, desktop, etc.) y eligiendo el formato de salida que prefieras.

### Características principales:
- ✅ Redimensionamiento inteligente manteniendo la proporción original
- ✅ Soporte para múltiples formatos de entrada: **JPG**, **PNG**, **WebP** y **CR3** (Canon RAW)
- ✅ Salida flexible: **WebP** o **JPEG**
- ✅ 6 niveles de calidad configurables
- ✅ Genera múltiples tamaños en un único paso
- ✅ Detección automática de orientación (portrait/landscape)

---

## 🛠️ Requisitos previos

### Para Windows:
1. **Go** (compilador) — opcional si solo vas a ejecutar el `.exe` precompilado
2. **ExifTool** — **OBLIGATORIO si usarás archivos CR3**
   - Descarga desde: https://exiftool.org/
   - Windows Executable: `exiftool(-k).exe`
   - Renómbralo a `exiftool.exe` y colócalo en una carpeta accesible (ej: `C:\Tools\exiftool\`)
   - **Añade la carpeta al PATH del sistema** (Variables de entorno)
   - Verifica en PowerShell: `exiftool -ver`

### Para Linux/WSL:
```bash
sudo apt-get install exiftool
```

---

## ▶️ Cómo usar la aplicación

### Paso 1: Prepara tus imágenes
1. Crea una carpeta llamada `input/` (o con el nombre que prefieras) en el mismo directorio que el programa
2. Coloca tus imágenes (JPG, PNG, WebP, CR3) en esa carpeta

### Paso 2: Ejecuta el programa

**Windows (PowerShell):**
```powershell
.\imageprocess.exe
```

**Linux/WSL:**
```bash
./imageprocess
```

---

## 📋 Flujo de preguntas interactivo

### 1️⃣ **Directorio de entrada** [default: `input`]
```
📁 Directorio de entrada [default: input]: 
```
- Presiona **Enter** para usar la carpeta por defecto `input/`
- O escribe la ruta de otra carpeta, ej: `C:\mis_fotos` o `/home/user/pictures`

### 2️⃣ **Directorio de salida** [default: `output`]
```
📁 Directorio de salida [default: output]: 
```
- Presiona **Enter** para usar la carpeta por defecto `output/`
- O especifica dónde guardar las imágenes procesadas

### 3️⃣ **Selecciona tamaños** (múltiple opción)
```
📐 Selecciona los tamaños deseados (separados por coma, ej: 1,2,3):
1) 420px (Móvil)
2) 720px (HD)
3) 1080px (Full HD)
4) 1440px (2K)
5) 2160px (4K)
6) Personalizado (tamaño en px)
```

**Ejemplos:**
- `1,3,5` → genera versiones de 420px, 1080px y 2160px
- `6` → te pide ingresar tamaños a medida, ej: `600, 1200, 1600`
- `1,2,3,4,5` → todas las opciones estándar

**Nota:** Si pides un tamaño mayor al original, se limitará al ancho de la imagen para evitar pixelado.

### 4️⃣ **Selecciona la calidad**
```
💎 Selecciona la calidad de salida:
1) Sin pérdida / Pérdida mínima (Lossless)
2) Calidad Máxima (95%)
3) Calidad Alta (80%) [Predeterminado]
4) Calidad Media-Alta (65%)
5) Calidad Media (45%)
6) Calidad Baja / Menor peso (25%)
```

**Guía de elección:**
| Opción | Uso recomendado | Tamaño aprox. |
|--------|-----------------|---------------|
| 1 (Lossless) | Logos, gráficos vectoriales | Grande |
| 2 (95%) | Fotografía profesional | Medio-Grande |
| 3 (80%) | **Sitios web estándar** | Medio |
| 4 (65%) | Web optimizado con buen aspecto | Pequeño |
| 5 (45%) | Miniaturas, previsualizaciones | Muy pequeño |
| 6 (25%) | Máxima compresión (calidad baja) | Mínimo |

### 5️⃣ **Selecciona el formato de salida**
```
🗂️  Selecciona el formato de salida:
1) WebP [Predeterminado]
2) JPEG
```

**Diferencias:**
| Formato | Ventajas | Desventajas |
|---------|----------|------------|
| **WebP** | Mejor compresión, soporte moderno | No soportado en IE, algunos navegadores antiguos |
| **JPEG** | Amplia compatibilidad, navegadores antiguos | Tamaño más grande, siempre con pérdida |

**Nota:** Si eliges **JPEG + Lossless**, se convertirá automáticamente a JPEG con calidad 100 (ya que JPEG siempre tiene pérdida).

---

## 📁 Estructura de archivos de salida

Los archivos generados tienen este formato:
```
<nombre_original>_<orientación>_<tamaño>.<extensión>
```

**Ejemplo:** Si procesan `foto.jpg` con tamaños 420px, 1080px en formato WebP:
```
output/
├── foto_landscape_420.webp
├── foto_landscape_1080.webp
└── ...
```

O si eliges JPEG:
```
output/
├── foto_landscape_420.jpg
├── foto_landscape_1080.jpg
└── ...
```

**Componentes:**
- `foto` = nombre original sin extensión
- `landscape` o `portrait` = orientación detectada automáticamente
- `420`, `1080` = ancho en píxeles
- `webp` o `jpg` = formato elegido

---

## 💡 Casos de uso

### Caso 1: Fotos para sitio web responsive
```
📁 Entrada: input
📁 Salida: output
Tamaños: 1,2,3,4,5 (todos)
Calidad: 3 (80%) — buena relación tamaño/calidad
Formato: 1 (WebP)
```
**Resultado:** 5 versiones por imagen (móvil, tablet, desktop) en formato moderno

### Caso 2: Galería de miniaturas
```
Tamaños: 6 → 300
Calidad: 5 (45%) — tamaño mínimo
Formato: 2 (JPEG)
```
**Resultado:** Miniaturas rápidas de cargar, máxima compatibilidad

### Caso 3: Procesamiento de CR3 (Canon RAW)
```
📁 Entrada: mi_viaje/
Tamaños: 1080
Calidad: 2 (95%) — máxima calidad
Formato: 1 (WebP)
```
**Requisito:** `exiftool` instalado y en PATH
**Nota:** Extrae el JPEG embebido en el CR3 (no revelado RAW completo)

---

## ⚠️ Solución de problemas

### ❌ "exiftool no encontrado en PATH"
**Problema:** Intentas procesar un CR3 pero la aplicación no encuentra `exiftool`.

**Solución:**
1. Instala ExifTool desde https://exiftool.org/
2. Asegúrate de que `exiftool.exe` esté en una carpeta agregada al PATH
3. **Reinicia PowerShell/Terminal** después de cambiar el PATH
4. Verifica: `exiftool -ver` (debe mostrar la versión)
5. Si sigue sin funcionar, mueve `exiftool.exe` a `C:\Windows\System32\` (requiere admin)

### ❌ "Error al leer entrada"
**Problema:** La carpeta de entrada no existe o está vacía.

**Solución:**
- Verifica que la ruta sea correcta
- Asegúrate de que las imágenes estén dentro de la carpeta indicada
- Usa rutas completas: `C:\Users\Usuario\Imágenes` en lugar de solo `Imágenes`

### ❌ "No seleccionaste ningún tamaño válido"
**Problema:** Las opciones de tamaño que escribiste no son válidas.

**Solución:**
- Separa los números con comas: `1,2,3` (no espacios entre comas)
- Solo ingresa números del 1 al 6
- Si usas opción 6 (personalizado), escribe números separados por coma: `600,1200,1800`

### ⚠️ Las imágenes CR3 se saltan
**Problema:** Tienes archivos CR3 pero se reportan en rojo como error.

**Causa probable:** `exiftool` no está instalado o no está en el PATH
**Solución:** Ver sección anterior "exiftool no encontrado en PATH"

### 💾 Los archivos generados son muy grandes
**Problema:** Las imágenes de salida pesan demasiado.

**Solución:**
- Reduce la calidad: elige opción 5 o 6
- Usa WebP en lugar de JPEG (mejor compresión)
- Solicita tamaños menores
- Ejemplo: para web, 80% de calidad (opción 3) es generalmente suficiente

---

## 📊 Consejos de optimización

| Caso | Tamaño | Calidad | Formato |
|------|--------|---------|---------|
| Sitio web moderno | 420, 720, 1080 | 80% | WebP |
| Compatibilidad máxima | 800, 1200 | 85% | JPEG |
| Redes sociales | 1080, 1440 | 75% | WebP |
| Impresión web | 2160 | 95% | WebP |
| Miniaturas | 300 | 50% | JPEG |

---

## ✨ Preguntas frecuentes

**P: ¿Puedo procesar varias carpetas?**
A: Actualmente solo procesa una carpeta a la vez. Repite el programa para cada carpeta.

**P: ¿Qué pasa si hay archivos no compatibles en la carpeta?**
A: Se ignoran automáticamente. Solo se procesan JPG, PNG, WebP y CR3.

**P: ¿Sobrescribe archivos existentes?**
A: Sí. Si ya existe un archivo con el mismo nombre en `output/`, se reemplaza.

**P: ¿WebP tiene buena compatibilidad?**
A: Sí. Soportado en Chrome, Edge, Firefox (moderno), Safari 16+. Para IE o navegadores muy antiguos, usa JPEG.

**P: ¿Dónde veo el progreso?**
A: Aparece en la consola con emojis: ✅ para éxito, ❌ para errores, ⚠️ para advertencias.

---

**Versión:** 1.0 | **Última actualización:** Mayo 2026
