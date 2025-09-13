
:root {
	--color-bg: {{ .BgColor }};
	--color-fg: {{ .FgColor }};
	--color-link: {{ .LinkColor }};
	--color-link-hover: {{ .LinkHoverColor }};
    --color-strong: {{ .StrongColor }};
	--color-header: {{ .HeaderColor }};
	--color-sidebar-bg: {{ .SidebarBgColor }};
	--color-sidebar-fg: {{ .SidebarFgColor }};
	--color-code-bg: {{ .CodeBgColor }};
	--color-code-fg: {{ .CodeFgColor }};
	--color-blockquote-bg: {{ .BlockquoteBgColor }};
	--color-blockquote-border: {{ .BlockquoteBorderColor }};
    --color-li-marker: {{ .ListBulletColor }};
	
	--external-link-icon-size: 0.75em;
}

body {
	background: var(--color-bg);
	color: var(--color-fg);
    font-family: "Comfortaa", sans-serif;
    font-size: 20px;
}

a {
	position: relative;
	color: var(--color-link);
	text-decoration: none;
	padding-right: var(--external-link-icon-size);
}
a:hover {
	color: var(--color-link-hover);
}
a::after {
	content: "";
	position: absolute;
	left: auto;
	right: 0;
	top: calc(-0.25 * var(--external-link-icon-size));
	height: var(--external-link-icon-size);
	width: var(--external-link-icon-size);
	background-image: url('_viki_static/external-link.png');
	background-repeat: no-repeat;
	background-size: contain;
	pointer-events: none;
}

header, h1, h2, h3, h4, h5, h6 {
    font-family: "Caveat";
	color: var(--color-header);
}

pre, code {
	background: var(--color-code-bg);
	color: var(--color-code-fg);
	font-family: "Fira Mono", "Consolas", "Monaco", monospace;
	border-radius: 4px;
	padding: 0.2em 0.4em;
}
pre {
	padding: 1em;
	overflow-x: auto;
}

strong {
    color: var(--color-strong);
}

blockquote {
	background: var(--color-blockquote-bg);
	border-left: 4px solid var(--color-blockquote-border);
	margin: 1em 0;
	padding: 0.5em 1em;
	color: var(--color-fg);
}

ul, ol {
	color: var(--color-fg);
}

li::marker {
    color: var(--color-li-marker);
}

hr {
	border: none;
	border-top: 1px solid var(--color-blockquote-border);
	margin: 2em 0;
}

table {
	border-collapse: collapse;
	width: 100%;
	background: var(--color-bg);
	color: var(--color-fg);
}
th, td {
	border: 1px solid var(--color-blockquote-border);
	padding: 0.5em 1em;
}
th {
	background: var(--color-sidebar-bg);
	color: var(--color-header);
}

.sidebar {
	background: var(--color-sidebar-bg);
	color: var(--color-sidebar-fg);
}
