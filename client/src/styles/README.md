# Styles Directory

This directory contains all global styles for the Dossier application, organized into modular CSS files for maintainability and reusability.

## File Structure

```
styles/
├── main.css          # Main entry point that imports all other styles
├── variables.css     # Design tokens (colors, spacing, etc.)
├── reset.css         # Browser reset and base styles
├── buttons.css       # Global button styles and variants
├── forms.css         # Input, textarea, and select styles
├── components.css    # Reusable component styles (cards, alerts, etc.)
└── utilities.css     # Single-purpose utility classes
```

## Usage

The global styles are automatically imported in `src/main.js`:

```javascript
import "./styles/main.css";
```

## File Purposes

### `variables.css`

Contains all CSS custom properties (design tokens) for:

- Colors (backgrounds, text, borders, semantic colors)
- Spacing values
- Border radius values
- Shadow definitions
- Effects (blur, transitions)
- Layout constraints
- Z-index layers

Use these variables throughout your component styles for consistency.

### `reset.css`

Provides baseline styles and browser normalization:

- Global box-sizing
- Body font and colors
- Margin/padding reset

### `buttons.css`

All button styles including variants:

- Base button styles
- `.primary` - Primary action buttons
- `.secondary` - Secondary action buttons
- `.success` - Success/confirmation buttons
- `.danger` - Destructive action buttons
- Disabled states

### `forms.css`

Form element styling:

- Input fields
- Textareas
- Select dropdowns
- Focus states
- Placeholder styles

### `components.css`

Reusable UI component styles:

- `.card` - Card container
- `.error` - Error message alert
- `.success` - Success message alert
- `.loading` - Loading spinner animation

### `utilities.css`

Single-purpose helper classes:

- Typography utilities (`.text-sm`, `.font-medium`, etc.)
- Layout utilities (`.flex`, `.items-center`, etc.)
- Spacing utilities (`.gap-2`, `.mb-4`, etc.)
- Opacity utilities

## Best Practices

1. **Use CSS Variables**: Reference design tokens from `variables.css` instead of hardcoding values
2. **Scoped Styles**: Use `<style scoped>` in Vue components for component-specific styles
3. **Utility-First**: Use utility classes for simple one-off styling needs
4. **Component Classes**: Create reusable component classes for repeated UI patterns
5. **Avoid !important**: Structure your CSS to avoid needing !important flags

## Adding New Styles

- **New design tokens**: Add to `variables.css`
- **New component patterns**: Add to `components.css`
- **New utilities**: Add to `utilities.css`
- **Major new category**: Create a new file and import it in `main.css`

## Naming Conventions

- Use kebab-case for class names
- Use BEM (Block Element Modifier) for complex components
- Prefix utility classes with their purpose (`.text-`, `.flex-`, etc.)
- Use semantic names over presentational names when possible
