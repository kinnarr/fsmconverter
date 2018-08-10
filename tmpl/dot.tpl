digraph graphname
{
  {{range $stateName, $state := .States -}}
  {{range $nextName, $next := $state.Successors -}}
    {{$stateName}} -> {{$nextName}};
  {{end -}}
  {{range $elseName, $elseValue := $state.DefaultSuccessor -}}
    {{$stateName}} -> {{$elseName}} [style=dashed];
  {{end -}}
  {{end}}
}
