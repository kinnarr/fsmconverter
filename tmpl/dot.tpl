digraph graphname
{
  {{range $stateName, $state := .States -}}
  {{range $nextName, $next := $state.Successors -}}
    {{$stateName}} -> {{$nextName}};
  {{end}}
  {{- end}}
}
