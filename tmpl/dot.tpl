digraph graphname
{
  {{range $stateName, $state := .States -}}
  {{range $nextName, $next := $state.Successors -}}
    {{$stateName}} -> {{$nextName}};
  {{end -}}
  {{range $elseName, $elseValue := $state.DefaultSuccessor -}}
    {{$stateName}} -> {{$elseName}} [style=dashed];
  {{end -}}
  {{if and (eq 0 (len $state.Outputs)) (not $state.Preserve)}}{{$stateName}} [color=red];{{end}}
  {{end}}
}
