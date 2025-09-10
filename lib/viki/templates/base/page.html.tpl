{{define "base"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Caveat:wght@400..700&family=Patrick+Hand&display=swap" rel="stylesheet">
    <style>
        body {
            margin: 0;
        }
        .container {
            display: flex;
            min-height: 100vh;
        }
        .sidebar {
            width: 400px;
            padding: 1rem;
            box-sizing: border-box;
        }
        .sidebar-list {
            list-style: none;
            margin: 0;
            padding-left: 0;
        }
        .sidebar-list li {
            margin: 0;
            padding: 0 0 0 0.5em;
        }
        .collapsible {
            cursor: pointer;
            user-select: none;
            font-weight: bold;
            display: inline-block;
        }
        .collapsible-content {
            margin-left: 1.2em;
        }
        .body {
            flex: 1;
            padding: 2rem;
        }
        .wikilink {
            font-weight: bold;
        }
    </style>
    <link rel="stylesheet" href="/theme.css">
</head>
<body>
    <div class="container">
        <aside class="sidebar">
            {{ .SidebarHtml }}
        </aside>
        <main class="body">
            {{ .BodyHtml }}
        </main>
    </div>
</body>
</html>
{{end}}