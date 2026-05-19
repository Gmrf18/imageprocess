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
- ✅ **Selección interactiva con checkboxes y listas** (flechas + espacio + Enter), con menú numérico clásico como alternativa
- ✅ Validación del directorio de entrada: si no existe, vuelve a preguntar

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

> ℹ️ **Dos modos de selección**
>
> - **Modo interactivo** (al ejecutar en una terminal normal): los tamaños se eligen con **checkboxes** y la calidad/formato con **listas navegables**. Usa **↑/↓** para moverte, **espacio** para marcar/desmarcar y **Enter** para confirmar.
> - **Modo clásico/numérico** (cuando la entrada viene de una tubería, script o sin terminal interactiva): se muestran los menús numerados y escribes el/los número(s), igual que antes.
>
> La aplicación detecta el modo automáticamente; no hay que configurar nada. Los directorios de entrada/salida siempre se piden como texto.

### 1️⃣ **Directorio de entrada** [default: `input`]
```
📁 Directorio de entrada [default: input]: 
```
- Presiona **Enter** para usar la carpeta por defecto `input/`
- O escribe la ruta de otra carpeta, ej: `C:\mis_fotos` o `/home/user/pictures`
- ⚠️ Si la carpeta indicada **no existe**, la aplicación lo avisa y vuelve a preguntar hasta que indiques una carpeta válida

### 2️⃣ **Directorio de salida** [default: `output`]
```
📁 Directorio de salida [default: output]: 
```
- Presiona **Enter** para usar la carpeta por defecto `output/`
- O especifica dónde guardar las imágenes procesadas

### 3️⃣ **Selecciona tamaños** (múltiple opción)

Tamaños disponibles:

| # | Tamaño |
|---|--------|
| 1 | 240px (Watch/Legacy) |
| 2 | 280px (Galaxy Fold cerrado) |
| 3 | 320px (Móvil S) |
| 4 | 468px (Móvil L) |
| 5 | 768px (Tablet) |
| 6 | 1080px (Full HD) |
| 7 | 1440px (2K) |
| 8 | 2160px (4K) |
| 9 | Personalizado (tamaño en px) |

**🖱️ Modo interactivo (checkboxes):**
```
? Selecciona los tamaños deseados:  [usa flechas, espacio para marcar, enter para confirmar]
  [ ] 240px (Watch/Legacy)
  [x] 320px (Móvil S)
  [x] 768px (Tablet)
  [ ] Personalizado (tamaño en px)
```
- **↑/↓**: mover · **espacio**: marcar/desmarcar · **Enter**: confirmar
- Marca todos los que necesites en una sola pantalla
- Si marcas **Personalizado**, después se te pedirá: `Introduce los tamaños en px (ej: 800, 1200):`

**⌨️ Modo clásico (numérico):** escribe los números separados por coma. Ejemplos:
- `1,3,5` → genera 240px, 320px y 768px
- `7,8` → genera 1440px (2K) y 2160px (4K)
- `9` → te pide tamaños a medida, ej: `600, 1200, 1600`

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

En **modo interactivo** se muestra como una lista de una sola opción (↑/↓ + Enter), con **Calidad Alta (80%)** preseleccionada por defecto. En **modo clásico** escribes el número `1`–`6` (si lo dejas vacío o inválido, usa la opción 3).

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

En **modo interactivo** es una lista de una sola opción (↑/↓ + Enter) con **WebP** por defecto. En **modo clásico** escribes `1` (WebP) o `2` (JPEG).

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
Tamaños: marca 320px, 768px y 1080px   (clásico: 3,5,6)
Calidad: Calidad Alta (80%)            (clásico: 3)
Formato: WebP                          (clásico: 1)
```
**Resultado:** 3 versiones por imagen (móvil, tablet, desktop) en formato moderno

### Caso 2: Galería de miniaturas
```
Tamaños: marca "Personalizado" → 300   (clásico: 9 → 300)
Calidad: Calidad Media (45%)           (clásico: 5)
Formato: JPEG                          (clásico: 2)
```
**Resultado:** Miniaturas rápidas de cargar, máxima compatibilidad

### Caso 3: Procesamiento de CR3 (Canon RAW)
```
📁 Entrada: mi_viaje/
Tamaños: marca "Personalizado" → 1080  (clásico: 9 → 1080)
Calidad: Calidad Máxima (95%)          (clásico: 2)
Formato: WebP                          (clásico: 1)
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

### ❌ "La carpeta '...' no existe o no es un directorio"
**Problema:** La ruta de entrada que indicaste no existe.

**Solución:**
- La aplicación **vuelve a preguntar** automáticamente; escribe una ruta válida o pulsa Enter para usar `input/`
- Usa rutas completas: `C:\Users\Usuario\Imágenes` en lugar de solo `Imágenes`
- Asegúrate de que las imágenes estén dentro de la carpeta indicada

### ❌ "No seleccionaste ningún tamaño válido"
**Problema:** No marcaste/elegiste ningún tamaño válido.

**Solución:**
- **Modo interactivo:** usa **espacio** para marcar al menos un tamaño antes de pulsar Enter (las flechas solo mueven el cursor, no seleccionan)
- **Modo clásico:** separa los números con comas: `1,2,3` (sin espacios), usando números del **1 al 9**
- Si usas la opción **9 (personalizado)**, escribe los px separados por coma: `600,1200,1800`

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
| Sitio web moderno | 320, 768, 1080 | Alta (80%) | WebP |
| Compatibilidad máxima | 768, 1080 (custom) | Máxima (95%) | JPEG |
| Redes sociales | 1080, 1440 | Media-Alta (65%) | WebP |
| Impresión web | 2160 (4K) | Máxima (95%) | WebP |
| Miniaturas | 300 (custom) | Media (45%) | JPEG |

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

**P: Me salen menús numerados en vez de checkboxes, ¿por qué?**
A: La selección interactiva (checkboxes/listas) requiere una terminal real. Si ejecutas el programa con la entrada redirigida (tubería, script, automatización) o sin terminal interactiva, se usa automáticamente el menú numérico clásico. Ejecútalo directamente en PowerShell/Terminal para ver los checkboxes.

**P: ¿En modo checkbox por qué no se marca nada con las flechas?**
A: Las flechas **↑/↓** solo mueven el cursor. Para marcar/desmarcar una opción usa la tecla **espacio**, y **Enter** para confirmar la selección.

---

**Versión:** 1.1 | **Última actualización:** Mayo 2026
