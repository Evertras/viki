{{define "base"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <style>
        body {
            margin: 0;
            font-family: sans-serif;
        }
        .container {
            display: flex;
            min-height: 100vh;
        }
        .sidebar {
            width: 400px;
            background: #f4f4f4;
            padding: 1rem;
            box-sizing: border-box;
        }
        .body {
            flex: 1;
            padding: 2rem;
            background: #fff;
        }
    </style>
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