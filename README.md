# mmd-exporter

Un pequeño CLI en Go para convertir archivos `.mmd` (Mermaid) en imágenes `.svg` y `.png` usando [`mmdc`](https://github.com/mermaid-js/mermaid-cli). Útil para automatizar la exportación de diagramas desde línea de comandos o al guardar cambios en tu editor.



## 🚀 Instalación

### Requisitos previos

- Tener instalado [Go](https://go.dev/doc/install) (`>=1.18`)
- Tener instalado [Node.js](https://nodejs.org/es/download/) (`>=20`)

#### Instalar dependencias

```bash
npm install -g playwright
npm install -g @mermaid-js/mermaid-cli
```

### 📦 Uso

#### Instalar CLI

```bash
go install github.com/deivguerrero/mmd-exporter@latest
```



#### Procesar todos los archivos .mmd en un directorio (solo el nivel actual)

```bash
mmd-exporter .
```

#### Procesar un archivo .mmd específico

```bash
mmd-exporter diagramas/flujo.mmd
```

#### Modo automático con --watch

```bash
# Observa un directorio y vuelve a exportar automáticamente los .svg y .png al modificar un .mmd.
mmd-exporter --watch .
```

### 📄 Qué genera

Para cada archivo diagrama.mmd, se crean:
	•	diagrama.svg
	•	diagrama.png (con --scale 2 para mejor calidad)

