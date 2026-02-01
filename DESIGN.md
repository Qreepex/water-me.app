# Design Theme: Emerald Oasis (V3)

**Core-Farbe:** #00EE57 (Vibrant Emerald)

---

## 1. Farb-Modi 🌗

### Light Mode (Frisch & Belebend)

- **Hintergrund:** `#F2FAF5` (Leichtes Mint-Weiß)
- **Oberfläche/Cards:** `#FFFFFF` (Reinweiß)
- **Text (Primary):** `#061F12` (Dunkles Moos)
- **Text (Secondary):** `#4B6658` (Gedämpftes Salbei)

### Dark Mode (Edel & Immersiv)

- **Hintergrund:** `#08140E` (Tiefes Wald-Schwarzgrün)
- **Oberfläche/Cards:** `#12261D` (Dschungel-Grün)
- **Text (Primary):** `#E0F2E9` (Eukalyptus-Weiß)
- **Text (Secondary):** `#7A9485` (Nebliges Grün)

---

## 2. Funktions- & Warnfarben (Full Spectrum) 🚦

| Status                  | Farbe          | Hex       | Logik                                                   |
| :---------------------- | :------------- | :-------- | :------------------------------------------------------ |
| **Kritisch / Fehler**   | **Berry Red**  | `#E63946` | Kontrast zu Grün, signalisiert "Gießen vergessen!"      |
| **Warnung**             | **Amber Glow** | `#FFB703` | Warmes Gelb-Orange für "Boden wird trocken".            |
| **Erinnerung**          | **Azure Sky**  | `#219EBC` | Kühles Blau für allgemeine Termine (z.B. Umtopfen).     |
| **Info (Grün-basiert)** | **Deep Grass** | `#008F41` | Wichtige Info, die positiv ist (z.B. "Optimale Sonne"). |
| **Sanfter Hinweis**     | **Soft Lime**  | `#D9FFEA` | Hintergrundfarbe für Erfolgs-Banner (Light Mode).       |

---

## 3. UI-Komponenten (Gradients & Styles) 💎

### Header Gradient (CSS)

Nutze diesen Verlauf für den oberen Bereich der App, um Tiefe zu erzeugen:
`background: linear-gradient(135deg, #00EE57 0%, #00A343 100%);`

### Wichtige Buttons

- **Background:** `#00EE57`
- **Text-Farbe:** `#08140E` (Immer dunkler Text auf hellem Primärgrün für Lesbarkeit!)
- **Shadow:** `0px 4px 14px rgba(0, 238, 87, 0.4)` (Ein leuchtender Schatten lässt den Button schweben).

---

## 4. CSS Variablen (Code Ready) 🧩

```css
:root {
 /* Brand Core */
 --p-emerald: #00ee57;
 --p-emerald-dark: #00a343;

 /* Status */
 --status-error: #e63946;
 --status-warn: #ffb703;
 --status-info: #219ebc;
 --status-success: #008f41;

 /* Theme Light */
 --bg-light: #f2faf5;
 --card-light: #ffffff;
 --text-light-main: #061f12;

 /* Theme Dark */
 --bg-dark: #08140e;
 --card-dark: #12261d;
 --text-dark-main: #e0f2e9;
}
```
