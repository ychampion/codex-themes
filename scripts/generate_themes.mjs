import { mkdirSync, writeFileSync, rmSync } from "node:fs";
import { join } from "node:path";

const root = process.cwd();
const themesDir = join(root, "themes");

const themes = [
  makeTheme({
    id: "amber-nocturne",
    name: "Amber Nocturne",
    family: "signature",
    appearance: "dark",
    description: "Warm premium amber theme designed specifically for long coding sessions.",
    tags: ["amber", "premium", "warm", "dark"],
    colors: base({
      bg: "#12110F", panel: "#191714", panel2: "#211D18", fg: "#F6F1E8", secondary: "#E7DDCF", muted: "#B8AFA2",
      accent: "#D6A24A", accent2: "#C97A2B", red: "#EA6D5D", yellow: "#E2A84B", green: "#78C26D",
      cyan: "#63D2CF", blue: "#80BFFF", purple: "#E2A7FF"
    })
  }),
  makeTheme({
    id: "dracula",
    name: "Dracula",
    family: "classic",
    appearance: "dark",
    description: "Punchy purple-pink neon contrast with readable violet surfaces.",
    tags: ["dracula", "purple", "dark"],
    colors: base({
      bg: "#282A36", panel: "#303341", panel2: "#3A3E4D", fg: "#F8F8F2", secondary: "#E9E8F2", muted: "#B8B6C9",
      accent: "#BD93F9", accent2: "#FF79C6", red: "#FF5555", yellow: "#F1FA8C", green: "#50FA7B",
      cyan: "#8BE9FD", blue: "#6272A4", purple: "#FF79C6"
    })
  }),
  makeTheme({
    id: "nord",
    name: "Nord",
    family: "classic",
    appearance: "dark",
    description: "Cool arctic blues with understated contrast and crisp structure.",
    tags: ["nord", "arctic", "dark"],
    colors: base({
      bg: "#2E3440", panel: "#3B4252", panel2: "#434C5E", fg: "#ECEFF4", secondary: "#E5E9F0", muted: "#A3B1C6",
      accent: "#88C0D0", accent2: "#81A1C1", red: "#BF616A", yellow: "#EBCB8B", green: "#A3BE8C",
      cyan: "#8FBCBB", blue: "#5E81AC", purple: "#B48EAD"
    })
  }),
  makeTheme({
    id: "catppuccin-mocha",
    name: "Catppuccin Mocha",
    family: "modern",
    appearance: "dark",
    description: "Soft dessert-like dark palette with polished pastel accents.",
    tags: ["catppuccin", "mocha", "pastel"],
    colors: base({
      bg: "#1E1E2E", panel: "#252538", panel2: "#313244", fg: "#CDD6F4", secondary: "#BAC2DE", muted: "#9399B2",
      accent: "#89B4FA", accent2: "#F5C2E7", red: "#F38BA8", yellow: "#F9E2AF", green: "#A6E3A1",
      cyan: "#94E2D5", blue: "#89B4FA", purple: "#CBA6F7"
    })
  }),
  makeTheme({
    id: "catppuccin-latte",
    name: "Catppuccin Latte",
    family: "modern",
    appearance: "light",
    description: "Creamy light palette with pastel accents and strong readability.",
    tags: ["catppuccin", "latte", "light"],
    colors: base({
      bg: "#EFF1F5", panel: "#E6E9EF", panel2: "#DCE0E8", fg: "#4C4F69", secondary: "#5C5F77", muted: "#7C7F93",
      accent: "#1E66F5", accent2: "#EA76CB", red: "#D20F39", yellow: "#DF8E1D", green: "#40A02B",
      cyan: "#179299", blue: "#1E66F5", purple: "#8839EF",
      inverse: "#FAFBFC"
    })
  }),
  makeTheme({
    id: "tokyo-night",
    name: "Tokyo Night",
    family: "modern",
    appearance: "dark",
    description: "Deep midnight blues with vibrant cyber accents and balanced contrast.",
    tags: ["tokyo-night", "blue", "dark"],
    colors: base({
      bg: "#1A1B26", panel: "#1F2335", panel2: "#24283B", fg: "#C0CAF5", secondary: "#A9B1D6", muted: "#7A88CF",
      accent: "#7AA2F7", accent2: "#BB9AF7", red: "#F7768E", yellow: "#E0AF68", green: "#9ECE6A",
      cyan: "#7DCFFF", blue: "#7AA2F7", purple: "#BB9AF7"
    })
  }),
  makeTheme({
    id: "gruvbox-dark",
    name: "Gruvbox Dark",
    family: "retro",
    appearance: "dark",
    description: "Warm retro earth tones with strong terminal heritage.",
    tags: ["gruvbox", "retro", "dark"],
    colors: base({
      bg: "#282828", panel: "#32302F", panel2: "#3C3836", fg: "#EBDBB2", secondary: "#D5C4A1", muted: "#A89984",
      accent: "#FABD2F", accent2: "#D79921", red: "#FB4934", yellow: "#FABD2F", green: "#B8BB26",
      cyan: "#8EC07C", blue: "#83A598", purple: "#D3869B"
    })
  }),
  makeTheme({
    id: "rose-pine",
    name: "Rosé Pine",
    family: "modern",
    appearance: "dark",
    description: "Muted mauve and pine tones for a calm premium coding surface.",
    tags: ["rose-pine", "mauve", "dark"],
    colors: base({
      bg: "#191724", panel: "#1F1D2E", panel2: "#26233A", fg: "#E0DEF4", secondary: "#C4A7E7", muted: "#908CAA",
      accent: "#9CCFD8", accent2: "#EBBCBA", red: "#EB6F92", yellow: "#F6C177", green: "#31748F",
      cyan: "#9CCFD8", blue: "#6E94B2", purple: "#C4A7E7"
    })
  }),
  makeTheme({
    id: "everforest",
    name: "Everforest",
    family: "classic",
    appearance: "dark",
    description: "Forest-inspired green neutrals with gentle contrast and low glare.",
    tags: ["everforest", "green", "dark"],
    colors: base({
      bg: "#2B3339", panel: "#323C41", panel2: "#3A464C", fg: "#D3C6AA", secondary: "#E1D8C3", muted: "#9DA9A0",
      accent: "#A7C080", accent2: "#DBBC7F", red: "#E67E80", yellow: "#DBBC7F", green: "#A7C080",
      cyan: "#83C092", blue: "#7FBBB3", purple: "#D699B6"
    })
  }),
  makeTheme({
    id: "monokai-pro",
    name: "Monokai Pro",
    family: "classic",
    appearance: "dark",
    description: "High-energy Monokai lineage with premium balanced saturation.",
    tags: ["monokai", "pro", "dark"],
    colors: base({
      bg: "#2D2A2E", panel: "#363338", panel2: "#403E41", fg: "#FCFCFA", secondary: "#F0EEE9", muted: "#A59FA1",
      accent: "#FFD866", accent2: "#AB9DF2", red: "#FF6188", yellow: "#FFD866", green: "#A9DC76",
      cyan: "#78DCE8", blue: "#78DCE8", purple: "#AB9DF2"
    })
  }),
  makeTheme({
    id: "one-dark-pro",
    name: "One Dark Pro",
    family: "classic",
    appearance: "dark",
    description: "Balanced editor classic with dependable neutral contrast.",
    tags: ["one-dark", "editor", "dark"],
    colors: base({
      bg: "#282C34", panel: "#2F343F", panel2: "#3A404B", fg: "#ABB2BF", secondary: "#C8CCD4", muted: "#7F8793",
      accent: "#61AFEF", accent2: "#C678DD", red: "#E06C75", yellow: "#E5C07B", green: "#98C379",
      cyan: "#56B6C2", blue: "#61AFEF", purple: "#C678DD"
    })
  }),
  makeTheme({
    id: "kanagawa",
    name: "Kanagawa",
    family: "modern",
    appearance: "dark",
    description: "Ink-wash inspired Japanese palette with elegant subdued contrast.",
    tags: ["kanagawa", "ink", "dark"],
    colors: base({
      bg: "#1F1F28", panel: "#26262F", panel2: "#2A2A37", fg: "#DCD7BA", secondary: "#C8C093", muted: "#727169",
      accent: "#7E9CD8", accent2: "#957FB8", red: "#C34043", yellow: "#C0A36E", green: "#76946A",
      cyan: "#6A9589", blue: "#7E9CD8", purple: "#957FB8"
    })
  }),
  makeTheme({
    id: "cyberpunk",
    name: "Cyberpunk",
    family: "experimental",
    appearance: "dark",
    description: "Neon magenta and cyan over pitch-black surfaces for dramatic sessions.",
    tags: ["cyberpunk", "neon", "dark"],
    colors: base({
      bg: "#0B0F14", panel: "#101722", panel2: "#152033", fg: "#E7F7FF", secondary: "#BFE8FF", muted: "#7CA6C1",
      accent: "#00F5FF", accent2: "#FF2BD6", red: "#FF5C8A", yellow: "#FFD166", green: "#72F1B8",
      cyan: "#00F5FF", blue: "#35A7FF", purple: "#FF2BD6"
    })
  }),
  makeTheme({
    id: "solarized-dark",
    name: "Solarized Dark",
    family: "classic",
    appearance: "dark",
    description: "Classic low-contrast precision palette tuned for long-form reading.",
    tags: ["solarized", "dark", "classic"],
    colors: base({
      bg: "#002B36", panel: "#073642", panel2: "#0A4755", fg: "#93A1A1", secondary: "#EEE8D5", muted: "#657B83",
      accent: "#268BD2", accent2: "#2AA198", red: "#DC322F", yellow: "#B58900", green: "#859900",
      cyan: "#2AA198", blue: "#268BD2", purple: "#6C71C4"
    })
  }),
  makeTheme({
    id: "solarized-light",
    name: "Solarized Light",
    family: "classic",
    appearance: "light",
    description: "Well-known paper-like light palette with precise muted accents.",
    tags: ["solarized", "light", "classic"],
    colors: base({
      bg: "#FDF6E3", panel: "#EEE8D5", panel2: "#E6DFCC", fg: "#586E75", secondary: "#657B83", muted: "#93A1A1",
      accent: "#268BD2", accent2: "#2AA198", red: "#DC322F", yellow: "#B58900", green: "#859900",
      cyan: "#2AA198", blue: "#268BD2", purple: "#6C71C4",
      inverse: "#FBFEF5"
    })
  })
];

rmSync(themesDir, { recursive: true, force: true });
mkdirSync(themesDir, { recursive: true });

const catalog = themes.map(({ id, name, family, appearance, description, tags }) => ({ id, name, family, appearance, description, tags }));

for (const theme of themes) {
  const dir = join(themesDir, theme.id);
  mkdirSync(dir, { recursive: true });
  writeFileSync(join(dir, "theme.toml"), renderTheme(theme));
  writeFileSync(join(dir, "preview.md"), `# ${theme.name}\n\n${theme.description}\n\nTags: ${theme.tags.join(", ")}\n`);
}

writeFileSync(join(themesDir, "catalog.json"), JSON.stringify(catalog, null, 2) + "\n");

function base({ bg, panel, panel2, fg, secondary, muted, accent, accent2, red, yellow, green, cyan, blue, purple, inverse }) {
  const bgOverlay = `${bg.replace('#', '#')}${bg.length === 7 ? 'CC' : ''}`;
  const inverseText = inverse ?? (luminance(bg) < 0.45 ? "#0F1115" : "#FAFBFC");
  return {
    bg_primary: bg,
    bg_secondary: panel,
    bg_tertiary: panel2,
    bg_overlay: bgOverlay,
    fg_primary: fg,
    fg_secondary: secondary,
    fg_muted: muted,
    fg_inverse: inverseText,
    border: panel2,
    border_subtle: panel,
    border_focus: accent,
    accent_primary: accent,
    accent_secondary: accent2,
    accent_soft: mix(bg, accent, 0.18),
    assistant_msg: blue,
    user_msg: purple,
    tool_output: green,
    shell_stdout: secondary,
    shell_stderr: red,
    status_bg: panel,
    status_fg: fg,
    status_accent: accent,
    input_bg: mix(bg, panel, 0.65),
    input_fg: fg,
    input_placeholder: muted,
    selection_bg: mix(panel2, accent, 0.22),
    selection_fg: fg,
    cursor: accent,
    link: blue,
    info: cyan,
    success: green,
    warning: yellow,
    error: red,
    diff_add: green,
    diff_remove: red,
    diff_change: yellow,
    code_inline: blue,
    syntax_plain: fg,
    syntax_comment: muted,
    syntax_keyword: purple,
    syntax_string: green,
    syntax_number: yellow,
    syntax_function: blue,
    syntax_type: cyan,
    syntax_variable: fg,
    syntax_operator: accent2,
    ansi: buildAnsi({ bg, fg, red, yellow, green, cyan, blue, purple })
  };
}

function buildAnsi({ bg, fg, red, yellow, green, cyan, blue, purple }) {
  return {
    black: bg,
    red,
    green,
    yellow,
    blue,
    magenta: purple,
    cyan,
    white: fg,
    bright_black: lighten(bg, 0.22),
    bright_red: lighten(red, 0.12),
    bright_green: lighten(green, 0.12),
    bright_yellow: lighten(yellow, 0.12),
    bright_blue: lighten(blue, 0.12),
    bright_magenta: lighten(purple, 0.12),
    bright_cyan: lighten(cyan, 0.12),
    bright_white: lighten(fg, 0.08)
  };
}

function makeTheme({ id, name, family, appearance, description, tags, colors }) {
  return {
    schema: "https://codexthemes.dev/schemas/theme.schema.json",
    id,
    name,
    family,
    version: "1.0.0",
    appearance,
    description,
    author: "Codex Themes",
    license: "MIT",
    homepage: "https://github.com/ychampion/codex-themes",
    tags,
    codex_min_version: "0.1.0",
    accessibility: {
      min_contrast_body: 4.5,
      min_contrast_ui: 3.0,
      color_blind_safe: true
    },
    tokens: Object.fromEntries(Object.entries(colors).filter(([key]) => key !== "ansi")),
    ansi: colors.ansi,
    targets: {
      codex: { config_patch: { mode: "merge", path: "config/codex.toml" } },
      windows_terminal: { scheme_name: name },
      iterm2: { profile_name: name },
      ghostty: { theme_name: name },
      alacritty: { theme_name: name }
    }
  };
}

function renderTheme(theme) {
  const lines = [
    `schema = "${theme.schema}"`,
    `id = "${theme.id}"`,
    `name = "${theme.name}"`,
    `family = "${theme.family}"`,
    `version = "${theme.version}"`,
    `appearance = "${theme.appearance}"`,
    `description = "${theme.description}"`,
    `author = "${theme.author}"`,
    `license = "${theme.license}"`,
    `homepage = "${theme.homepage}"`,
    `tags = [${theme.tags.map((tag) => `"${tag}"`).join(", ")}]`,
    `codex_min_version = "${theme.codex_min_version}"`,
    "",
    "[accessibility]",
    `min_contrast_body = ${theme.accessibility.min_contrast_body}`,
    `min_contrast_ui = ${theme.accessibility.min_contrast_ui}`,
    `color_blind_safe = ${theme.accessibility.color_blind_safe}`,
    "",
    "[tokens]"
  ];

  for (const [key, value] of Object.entries(theme.tokens)) {
    lines.push(`${key} = "${value}"`);
  }

  lines.push("", "[ansi]");
  for (const [key, value] of Object.entries(theme.ansi)) {
    lines.push(`${key} = "${value}"`);
  }

  lines.push("", "[targets.codex]", 'config_patch = { mode = "merge", path = "config/codex.toml" }');
  lines.push("", "[targets.windows_terminal]", `scheme_name = "${theme.name}"`);
  lines.push("", "[targets.iterm2]", `profile_name = "${theme.name}"`);
  lines.push("", "[targets.ghostty]", `theme_name = "${theme.name}"`);
  lines.push("", "[targets.alacritty]", `theme_name = "${theme.name}"`, "");
  return lines.join("\n");
}

function mix(hexA, hexB, ratio) {
  const a = rgb(hexA);
  const b = rgb(hexB);
  return toHex({
    r: a.r + (b.r - a.r) * ratio,
    g: a.g + (b.g - a.g) * ratio,
    b: a.b + (b.b - a.b) * ratio
  });
}

function lighten(hex, ratio) {
  return mix(hex, "#FFFFFF", ratio);
}

function luminance(hex) {
  const { r, g, b } = rgb(hex);
  return (0.2126 * r + 0.7152 * g + 0.0722 * b) / 255;
}

function rgb(hex) {
  const value = hex.replace("#", "").slice(0, 6);
  return {
    r: Number.parseInt(value.slice(0, 2), 16),
    g: Number.parseInt(value.slice(2, 4), 16),
    b: Number.parseInt(value.slice(4, 6), 16)
  };
}

function toHex({ r, g, b }) {
  return `#${[r, g, b].map((value) => Math.max(0, Math.min(255, Math.round(value))).toString(16).padStart(2, "0")).join("").toUpperCase()}`;
}


