{{/*
  Copyright 2018 Franz Schmidt

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  		http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/}}
module fsm (
  clock,
  reset,
  {{ range $inputName, $inputLenght := .Inputs -}}
  {{$inputName}},
  {{end -}}
  state
  );
  {{$countStates := len $.States}}
  {{$binaryStateSize := getBinarySize $countStates }}
  parameter SIZE = {{ $binaryStateSize }};

  input wire clock, reset;
  {{ range $inputName, $inputLenght := .Inputs -}}
  input wire {{if gt $inputLenght 1}}[{{minus $inputLenght 1}}:0] {{end}}{{$inputName}};
  {{end -}}

  output reg [SIZE-1:0] state;

  {{ $counter := 0 }}
  parameter {{ range $stateName, $state := .States }}{{upper $stateName}} = {{convertBinary $state.Statenumber $binaryStateSize}}{{if ne $counter (minus $countStates 1)}}, {{end}}{{ $counter = inc $counter}}{{end}};

  reg [SIZE-1:0] next_state;

  always @ (state{{ range $inputName, $inputLenght := .Inputs }} or {{$inputName}}{{end}})
  begin
    next_state = {{convertBinary 0 $binaryStateSize}};
    case(state)
    {{- range $stateName, $state := .States}}
      {{$countCondition := 0}}
      {{ upper $stateName}} : {{range $nextName, $next := $state.Successors -}}
                {{if not (emptyCondition $next)}}{{$countCondition = inc $countCondition}}if {{if conditionIsAnd $next}}({{conditionToString $next "&&"}}){{end -}}
                {{- if conditionIsOr $next}}({{conditionToString $next "||"}}){{end}}{{end}} begin
                  next_state = {{upper $nextName}};
                end {{if not (emptyCondition $next)}}else {{end}}{{end -}}{{if ne $countCondition 0}}begin
                  {{- range $elseName, $elseValue := $state.DefaultSuccessor}}
                  next_state = {{upper $elseName}};
                  {{- else}}
                  next_state = {{upper $stateName}};
                  {{- end}}
                end
                {{- end}}
    {{end -}}
      default : next_state = {{upper .Defaults.State}};
    endcase
  end

  initial state <= {{upper .Defaults.State}};

  always @ (posedge clock)
  begin
      state <= next_state;
  end
endmodule
