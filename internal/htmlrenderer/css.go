package htmlrenderer

import "github.com/vchitepu/weave/internal/theme"

// BaseCSS contains theme-independent layout styles using CSS custom properties.
const BaseCSS = `
*, *::before, *::after { box-sizing: border-box; }

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, sans-serif;
  line-height: 1.6;
  color: var(--color-text);
  background-color: var(--color-bg);
  max-width: 52em;
  margin: 0 auto;
  padding: 2em 1.5em;
}

h1, h2, h3, h4, h5, h6 {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  line-height: 1.3;
}
h1 { color: var(--color-h1); font-size: 2em; border-bottom: 2px solid var(--color-rule); padding-bottom: 0.3em; }
h2 { color: var(--color-h2); font-size: 1.5em; }
h3 { color: var(--color-h3); font-size: 1.25em; }
h4 { color: var(--color-h4); font-size: 1.1em; }
h5 { color: var(--color-h5); font-size: 1em; }
h6 { color: var(--color-h6); font-size: 0.9em; }

p { margin: 0.8em 0; }

a { color: var(--color-link); text-decoration: none; }
a:hover { text-decoration: underline; }

em { font-style: italic; }
strong { font-weight: bold; }
del { text-decoration: line-through; color: var(--color-dim); }

hr {
  border: none;
  border-top: 2px solid var(--color-rule);
  margin: 1.5em 0;
}

.inline-code {
  background-color: var(--color-code-bg);
  border: 1px solid var(--color-code-border);
  border-radius: 4px;
  padding: 0.15em 0.4em;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
  font-size: 0.9em;
}

blockquote {
  border-left: 4px solid var(--color-blockquote-bar);
  margin: 1em 0;
  padding: 0.5em 1em;
  color: var(--color-blockquote-text);
}
blockquote p { margin: 0.4em 0; }

ul, ol { padding-left: 2em; margin: 0.8em 0; }
li { margin: 0.3em 0; }

.task-item {
  list-style: none;
  margin-left: -1.5em;
}
.task-item .checkbox {
  display: inline-block;
  width: 1.2em;
  text-align: center;
  margin-right: 0.3em;
}
.task-item .checkbox.checked { color: var(--color-task-checked); }
.task-item .checkbox.unchecked { color: var(--color-task-unchecked); }

table {
  border-collapse: collapse;
  margin: 1em 0;
  width: auto;
}
th, td {
  border: 1px solid var(--color-table-border);
  padding: 0.4em 0.8em;
  text-align: left;
}
th {
  color: var(--color-table-header);
  background-color: var(--color-table-header-bg);
  font-weight: bold;
}

img { max-width: 100%; height: auto; border-radius: 4px; margin: 0.5em 0; }

.code-container {
  border: 1px solid var(--color-code-border);
  border-radius: 6px;
  margin: 1em 0;
  overflow: hidden;
}
.code-container .lang-badge {
  display: block;
  padding: 0.2em 0.8em;
  font-size: 0.8em;
  color: var(--color-code-badge-text);
  background-color: var(--color-code-bg);
  border-bottom: 1px solid var(--color-code-border);
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
}
.code-container pre {
  margin: 0;
  padding: 1em;
  overflow-x: auto;
  background-color: var(--color-code-bg);
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
  font-size: 0.9em;
  line-height: 1.5;
}
.code-container.container-tree { border-color: var(--color-tree-border); }
.code-container.container-tree .lang-badge { border-color: var(--color-tree-border); }
.code-container.container-diagram { border-color: var(--color-diagram-border); }
.code-container.container-diagram .lang-badge { border-color: var(--color-diagram-border); }
.code-container.container-shell { border-color: var(--color-shell-border); }
.code-container.container-shell .lang-badge { border-color: var(--color-shell-border); }

.file-separator {
  border: none;
  border-top: 1px dashed var(--color-rule);
  margin: 2em 0;
}
`

// ThemeCSS returns CSS :root declarations for the given theme.
func ThemeCSS(th theme.Theme) string {
	if th.ChromaStyle == "github-dark" {
		return darkThemeCSS
	}
	return lightThemeCSS
}

const darkThemeCSS = `:root {
  --color-text: #D0D4DC;
  --color-bg: #1A1E28;
  --color-h1: #8BA4D4;
  --color-h2: #A8A0D6;
  --color-h3: #C9A86A;
  --color-h4: #A88A55;
  --color-h5: #876C42;
  --color-h6: #665030;
  --color-rule: #3A3F4B;
  --color-link: #7FA3C8;
  --color-dim: #9BA3B2;
  --color-code-bg: #262B36;
  --color-code-border: #3A3F4B;
  --color-code-badge-text: #B7C0D0;
  --color-blockquote-bar: #3A3F4B;
  --color-blockquote-text: #9BA3B2;
  --color-table-border: #3A3F4B;
  --color-table-header: #8BA4D4;
  --color-table-header-bg: #262B36;
  --color-tree-border: #5D7290;
  --color-diagram-border: #7B74A6;
  --color-shell-border: #9F8656;
  --color-task-checked: #6BBF8A;
  --color-task-unchecked: #5A5F6B;
}
`

const lightThemeCSS = `:root {
  --color-text: #2C3340;
  --color-bg: #FAFBFC;
  --color-h1: #3F5F8A;
  --color-h2: #665E95;
  --color-h3: #8D6B3F;
  --color-h4: #A07850;
  --color-h5: #B38A63;
  --color-h6: #C69C78;
  --color-rule: #C2C7D0;
  --color-link: #496A92;
  --color-dim: #6C7483;
  --color-code-bg: #EEF1F5;
  --color-code-border: #C2C7D0;
  --color-code-badge-text: #5F6B7A;
  --color-blockquote-bar: #C2C7D0;
  --color-blockquote-text: #6C7483;
  --color-table-border: #C2C7D0;
  --color-table-header: #3F5F8A;
  --color-table-header-bg: #EEF1F5;
  --color-tree-border: #5E7699;
  --color-diagram-border: #7D73A3;
  --color-shell-border: #9A7B52;
  --color-task-checked: #3A8C5C;
  --color-task-unchecked: #9BA3B2;
}
`
