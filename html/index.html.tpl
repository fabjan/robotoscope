<!doctype html>
<html lang=en>

<head>
    <meta charset=utf-8>
    <title>{{.Title}}</title>
    <style>
        /* CSS extracted from http://bettermotherfuckingwebsite.com */
        
        body {
            margin: 40px auto;
            max-width: 650px;
            line-height: 1.6;
            font-size: 18px;
            color: #444;
            padding: 0 10px
        }
        
        h1,
        h2,
        h3 {
            line-height: 1.2
        }
    </style>
</head>

<body>

    <h2>Robots</h2>
    {{if not .Robots}}
    <p><em>(no data)</em></p>
    {{else}}
    <table>
        <tr>
            <th>Seen</th>
            <th>User-Agent</th>
        </tr>
        {{ range .Robots}}
        <tr>
            <td>{{ .Seen }}</td>
            <td>{{ .UserAgent }}</td>
        </tr>
        {{ end}}
    </table>
    {{end}}

    <h2>Cheaters</h2>
    {{if not .Cheaters}}
    <p><em>(no data)</em></p>
    {{else}}
    <table>
        <tr>
            <th>Seen</th>
            <th>User-Agent</th>
        </tr>
        {{ range .Cheaters}}
        <tr>
            <td>{{ .Seen }}</td>
            <td>{{ .UserAgent }}</td>
        </tr>
        {{ end}}
    </table>
    {{end}}
</body>

</html>