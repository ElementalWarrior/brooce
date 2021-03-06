package tpl

var mainPageTpl = `
{{ define "mainpage" }}
{{ template "header" "overview" }}
<div class="row">
  <div class="col-md-8">
    <h3>Queues</h3>
    <table class="table">
      <thead>
        <tr>
          <th>Queue</th>
          <th>Pending</th>
          <th>Running</th>
          <th>Done</th>
          <th>Failed</th>
          <th>Delayed</th>
        </tr>
      </thead>
      <tbody>
        {{ range $i, $Queue := .Queues }}
          <tr>
            <td>{{ $Queue.QueueName }}</td>
            <td><a href="/pending/{{ $Queue.QueueName }}">{{ $Queue.Pending }}</a></td>
            <td>{{ $Queue.Running }}</td>
            <td><a href="/done/{{ $Queue.QueueName }}">{{ $Queue.Done }}</a></td>
            <td><a href="/failed/{{ $Queue.QueueName }}">{{ $Queue.Failed }}</a></td>
            <td><a href="/delayed/{{ $Queue.QueueName }}">{{ $Queue.Delayed }}</a></td>
          </tr>
        {{ end }}
      </tbody>
    </table>

  </div>
</div>


<div class="row">
  <div class="col-md-12">
    <h3>{{ len .RunningWorkers }} Workers Alive</h3>
    <table class="table">
      <thead>
        <tr>
          <th>Worker Name</th>
          <th>Machine Name</th>
          <th>Machine IP</th>
          <th>Process ID</th>
          <th>Queues</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          {{ range $i, $Worker := .RunningWorkers }}
          <tr>
            <td>{{ $Worker.ProcName }}</td>
            <td>{{ $Worker.Hostname }}</td>
            <td>{{ $Worker.IP }}</td>
            <td>{{ $Worker.PID }}</td>
            <td>
              {{ range $j, $Queue := $Worker.Queues }}
                {{ $Queue.Workers }}x<tt>{{ $Queue.Name }}</tt>
              {{ end }}
            </td>
          </tr>
        {{ end }}
        </tr>
      </tbody>
    </table>
  </div>
</div>



<div class="row">
  <div class="col-md-12">

    <h3>{{ len .RunningJobs }} of {{ .TotalThreads }} Threads Working</h3>
    <table class="table">
      <thead>
        <tr>
          <th>Thread Name</th>
          <th>Queue</th>
          <th>Started</th>
          <th>Command</th>
          <th>Params</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {{ range .RunningJobs }}
          <tr>
            <td>{{ .WorkerThreadName }}</td>
            <td>{{ .QueueName }}</td>
            <td><span title="{{FormatTime .StartTime}}">{{ TimeSince .StartTime }}</span></td>
            <td><code>{{ .Command }}</code></td>
            <td class="params">
              <ul>
                {{ if .Timeout }} <li>Timeout: {{ TimeDuration .Timeout }} {{ end }}
                {{ if gt .MaxTries 1 }} <li>Max Tries: {{ .MaxTries }} {{ end }}
                {{ if .Cron }} <li>Cron: {{ .Cron }} {{ end }}
                {{ if .Locks }} <li>Locks: {{ Join .Locks ", " }} {{ end }}
              </ul>
            </td>

            <td class="buttons">
              {{ if .HasLog }}
                <a href="/showlog/{{ .Id }}" target="_new" class="btn btn-info btn-xs">
                  <span class="glyphicon glyphicon-align-justify"></span>
                  Show Log
                </a>
              {{ end }}
            </td>
          </tr>
        {{ end }}
      </tbody>
    </table>

  </div>
</div>
{{ template "footer" }}
{{ end }}
`
