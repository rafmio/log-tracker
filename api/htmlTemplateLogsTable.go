package main

const htmlTemplateLogsTable = `
<table>
        <thead>
            <tr>
                <th>SeqNum</th>
                <th>TmStmp</th>
                <th>SrcIP</th>
                <th>Len</th>
                <th>Ttl</th>
                <th>Id</th>
                <th>Spt</th>
                <th>Dpt</th>
                <th>Window</th>
            </tr>
        </thead>
        <tbody>
            {{range .}}
            <tr>
                <td>{{.SeqNum}}</td>
                <td>{{.TmStmp.Format "2006-01-02 15:04:05"}}</td>
                <td>{{.SrcIP}}</td>
                <td>{{.Len}}</td>
                <td>{{.Ttl}}</td>
                <td>{{.Id}}</td>
                <td>{{.Spt}}</td>
                <td>{{.Dpt}}</td>
                <td>{{.Window}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
`
